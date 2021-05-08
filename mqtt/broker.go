package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/genus-machina/ganglia"
	"github.com/genus-machina/ganglia/monitors"
)

type Broker struct {
	client        mqtt.Client
	event         *ganglia.DigitalEvent
	logger        *log.Logger
	mutex         sync.Mutex
	notifier      *monitors.DigitalNotifier
	subscriptions SubscriptionMap
}

func NewBroker(logger *log.Logger, options *MqttOptions) (*Broker, error) {
	broker := new(Broker)
	broker.logger = log.New(logger.Writer(), "[mqtt] ", logger.Flags())
	broker.notifier = new(monitors.DigitalNotifier)
	broker.subscriptions = make(SubscriptionMap)

	clientOptions := broker.buildClientOptions(options)
	if tlsConfig, err := broker.buildTlsConfig(options); err == nil {
		clientOptions.SetTLSConfig(tlsConfig)
	} else {
		return nil, err
	}

	broker.client = mqtt.NewClient(clientOptions)
	broker.client.Connect()
	broker.buildEvent(ganglia.Low)
	return broker, nil
}

func (broker *Broker) buildClientOptions(options *MqttOptions) *mqtt.ClientOptions {
	return mqtt.NewClientOptions().
		AddBroker(options.Broker).
		SetAutoReconnect(true).
		SetCleanSession(true).
		SetClientID(options.ClientId).
		SetConnectionAttemptHandler(broker.handleConnectionAttempt).
		SetConnectionLostHandler(broker.handleConnectionLost).
		SetConnectRetry(true).
		SetConnectRetryInterval(time.Second).
		SetKeepAlive(time.Minute).
		SetMaxReconnectInterval(time.Minute).
		SetOnConnectHandler(broker.handleConnect).
		SetReconnectingHandler(broker.handleReconnect).
		SetWill(DeviceStatusTopic(options.ClientId), StatusMessage(Offline), AtLeastOnce, true)
}

func (broker *Broker) buildEvent(value ganglia.DigitalValue) {
	broker.event = &ganglia.DigitalEvent{
		Time:  time.Now(),
		Value: value,
	}
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

func (broker *Broker) getHandlers(topic string) []MessageHandler {
	broker.mutex.Lock()
	defer broker.mutex.Unlock()

	var handlers []MessageHandler
	for _, handler := range broker.subscriptions[topic] {
		handlers = append(handlers, handler)
	}
	return handlers
}

func (broker *Broker) handleConnect(client mqtt.Client) {
	broker.logger.Println("Connected to MQTT broker.")
	broker.publishBirthMessage()
	broker.resubscribe()
	broker.buildEvent(ganglia.High)
	broker.notify()
}

func (broker *Broker) handleConnectionAttempt(brokerUrl *url.URL, config *tls.Config) *tls.Config {
	broker.logger.Printf("Attempting to connect to %s...\n", brokerUrl.String())
	return config
}

func (broker *Broker) handleConnectionLost(client mqtt.Client, err error) {
	broker.logger.Printf("Lost connection to MQTT broker. %s.\n", err.Error())
	broker.buildEvent(ganglia.Low)
	broker.notify()
}

func (broker *Broker) handleMessage(client mqtt.Client, message mqtt.Message) {
	wrappedMessage := Message(message.Payload())
	for _, handler := range broker.getHandlers(message.Topic()) {
		handler(wrappedMessage)
	}
	message.Ack()
}

func (broker *Broker) handleReconnect(client mqtt.Client, options *mqtt.ClientOptions) {
	broker.logger.Println("Reconnecting to broker...")
}

func (broker *Broker) notify() {
	broker.notifier.Notify(broker.event)
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
	for topic, _ := range broker.subscriptions {
		broker.client.Subscribe(topic, AtMostOnce, broker.handleMessage)
	}
}

func (broker *Broker) Publish(message []byte, topic string) {
	broker.client.Publish(topic, AtMostOnce, true, message)
}

func (broker *Broker) Subscribe(topic string, handler MessageHandler) {
	broker.mutex.Lock()
	defer broker.mutex.Unlock()

	broker.logger.Printf("Subscribing to topic '%s'.\n", topic)
	broker.client.Subscribe(topic, AtMostOnce, broker.handleMessage)
	handlers := broker.subscriptions[topic]
	handlers = append(handlers, handler)
	broker.subscriptions[topic] = handlers
}

type Message []byte

type MessageHandler func(Message)

type MqttOptions struct {
	Broker   string
	CaFile   string
	CertFile string
	ClientId string
	KeyFile  string
}

type SubscriptionMap map[string][]MessageHandler
