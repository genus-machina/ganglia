package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"sync"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

type Broker struct {
	client        mqtt.Client
	logger        *log.Logger
	mutex         sync.Mutex
	subscriptions SubscriptionMap
}

func NewBroker(logger *log.Logger, options *MqttOptions) (*Broker, error) {
	broker := new(Broker)
	broker.logger = log.New(logger.Writer(), "[mqtt] ", logger.Flags())
	broker.subscriptions = make(SubscriptionMap)

	clientOptions := broker.buildClientOptions(options)
	if tlsConfig, err := broker.buildTlsConfig(options); err == nil {
		clientOptions.SetTLSConfig(tlsConfig)
	} else {
		return nil, err
	}

	broker.client = mqtt.NewClient(clientOptions)
	broker.client.Connect()
	return broker, nil
}

func (broker *Broker) buildClientOptions(options *MqttOptions) *mqtt.ClientOptions {
	return mqtt.NewClientOptions().
		AddBroker(options.Broker).
		SetAutoReconnect(true).
		SetCleanSession(true).
		SetClientID(options.ClientId).
		SetConnectionLostHandler(broker.handleConnectionLost).
		SetKeepAlive(time.Minute).
		SetOnConnectHandler(broker.handleConnect).
		SetWill(DeviceStatusTopic(options.ClientId), StatusMessage(Online), AtLeastOnce, true)
}

func (broker *Broker) buildTlsConfig(options *MqttOptions) (*tls.Config, error) {
	caContents, err := ioutil.ReadFile(options.CaFile)
	if err != nil {
		return nil, err
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caContents)

	keyPair, err := tls.LoadX509KeyPair(options.CertFile, options.KeyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := new(tls.Config)
	tlsConfig.Certificates = []tls.Certificate{keyPair}
	tlsConfig.RootCAs = caPool

	return tlsConfig, nil
}

func (broker *Broker) Close() {
	broker.logger.Println("Closing MQTT broker.")
	broker.client.Disconnect(250)
}

func (broker *Broker) handleConnect(client mqtt.Client) {
	broker.logger.Println("Connected to MQTT broker.")
	broker.publishBirthMessage()
	broker.resubscribe()
}

func (broker *Broker) handleConnectionLost(client mqtt.Client, err error) {
	broker.logger.Printf("Lost connection to MQTT broker. %s.\n", err.Error())
}

func (broker *Broker) publishBirthMessage() {
	reader := broker.client.OptionsReader()
	options := &reader
	statusTopic := DeviceStatusTopic(options.ClientID())
	broker.client.Publish(statusTopic, AtLeastOnce, true, StatusMessage(Online))
}

func (broker *Broker) resubscribe() {
	broker.mutex.Lock()
	defer broker.mutex.Unlock()

	broker.logger.Println("Resubscribing to all previously subscribed topics.")
	for topic, handlers := range broker.subscriptions {
		for _, handler := range handlers {
			broker.client.Subscribe(topic, AtMostOnce, handler)
		}
	}
}

func (broker *Broker) Publish(message []byte, topic string) {
	broker.client.Publish(topic, AtMostOnce, true, message)
}

func (broker *Broker) Subscribe(topic string, handler mqtt.MessageHandler) {
	broker.mutex.Lock()
	defer broker.mutex.Unlock()

	broker.logger.Printf("Subscribing to topic '%s'.\n", topic)
	broker.client.Subscribe(topic, AtMostOnce, handler)
	handlers := broker.subscriptions[topic]
	handlers = append(handlers, handler)
	broker.subscriptions[topic] = handlers
}

type MqttOptions struct {
	Broker   string
	CaFile   string
	CertFile string
	ClientId string
	KeyFile  string
}

type SubscriptionMap map[string][]mqtt.MessageHandler
