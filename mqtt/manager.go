package mqtt

import (
	"encoding/json"
	"log"
	"reflect"
	"time"

	"github.com/genus-machina/ganglia"
	"github.com/genus-machina/ganglia/monitors"
)

type Manager struct {
	broker *Broker
	logger *log.Logger

	analogChannels        []chan *ganglia.AnalogEvent
	channels              []interface{}
	digitalChannels       []chan *ganglia.DigitalEvent
	environmentalChannels []chan *ganglia.EnvironmentalEvent
	rawChannels           []chan Message
}

func NewManager(logger *log.Logger, broker *Broker) *Manager {
	manager := new(Manager)
	manager.broker = broker
	manager.logger = log.New(logger.Writer(), "[mqtt] ", logger.Flags())
	return manager
}

func Connect(logger *log.Logger, options *MqttOptions) (*Manager, error) {
	broker, err := NewBroker(logger, options)
	if err != nil {
		return nil, err
	}

	return NewManager(logger, broker), nil
}

func (manager *Manager) AnalogInput(topic string) ganglia.AnalogInput {
	input := make(chan *ganglia.AnalogEvent, 1)
	manager.subscribe(topic, input)
	return input
}

func (manager *Manager) DigitalInput(topic string) ganglia.DigitalInput {
	input := make(chan *ganglia.DigitalEvent, 1)
	manager.subscribe(topic, input)
	return input
}

func (manager *Manager) DigitalOutput(topic string) ganglia.DigitalOutput {
	output := make(chan ganglia.DigitalValue)
	go manager.watchDigitalOutput(output, topic)
	return output
}

func (manager *Manager) EnvironmentalInput(topic string) ganglia.EnvironmentalInput {
	input := make(chan *ganglia.EnvironmentalEvent, 1)
	manager.subscribe(topic, input)
	return input
}

func (manager *Manager) EnvironmentalOutput(topic string) ganglia.EnvironmentalOutput {
	output := make(chan *ganglia.EnvironmentalEvent)
	go manager.watchEnvironmentalOutput(output, topic)
	return output
}

func (manager *Manager) Halt() {
	manager.broker.Close()

	for _, channel := range manager.channels {
		reflect.ValueOf(channel).Close()
	}
}

func (manager *Manager) handleEvent(input interface{}) MessageHandler {
	return func(message Message) {
		var event interface{}

		switch input.(type) {
		case chan *ganglia.AnalogEvent:
			event = new(ganglia.AnalogEvent)
		case chan *ganglia.DigitalEvent:
			event = new(ganglia.DigitalEvent)
		case chan *ganglia.EnvironmentalEvent:
			event = new(ganglia.EnvironmentalEvent)
		default:
			panic("Unknown event type.")
		}

		if err := json.Unmarshal(message, event); err == nil {
			channel := reflect.ValueOf(input)
			payload := reflect.ValueOf(event)
			channel.Send(payload)
		} else {
			manager.logger.Printf("Failed to parse event. %s.\n", err.Error())
		}
	}
}

func (manager *Manager) handleMessage(input chan Message) MessageHandler {
	return func(message Message) {
		input <- message
	}
}

func (manager *Manager) Monitor() monitors.DigitalMonitor {
	return NewMonitor(manager.broker)
}

func (manager *Manager) Publish(message interface{}, topic string) error {
	var err error
	var payload []byte

	if payload, err = json.Marshal(message); err == nil {
		manager.broker.Publish(payload, topic)
	}

	return err
}

func (manager *Manager) PublishAnalogInput(input ganglia.AnalogInput, topic string) {
	go manager.watchInput(input, topic)
}

func (manager *Manager) PublishDigitalInput(input ganglia.DigitalInput, topic string) {
	go manager.watchInput(input, topic)
}

func (manager *Manager) PublishEnvironmentalInput(input ganglia.EnvironmentalInput, topic string) {
	go manager.watchInput(input, topic)
}

func (manager *Manager) Subscribe(topic string) <-chan Message {
	input := make(chan Message, 1)
	manager.broker.Subscribe(topic, manager.handleMessage(input))
	manager.rawChannels = append(manager.rawChannels, input)
	return input
}

func (manager *Manager) subscribe(topic string, input interface{}) {
	manager.broker.Subscribe(topic, manager.handleEvent(input))
	manager.channels = append(manager.channels, input)
}

func (manager *Manager) watchDigitalOutput(output chan ganglia.DigitalValue, topic string) {
	for value := range output {
		event := &ganglia.DigitalEvent{time.Now(), value}
		payload, _ := json.Marshal(event)
		manager.broker.Publish(payload, topic)
	}
}

func (manager *Manager) watchEnvironmentalOutput(output chan *ganglia.EnvironmentalEvent, topic string) {
	for event := range output {
		payload, _ := json.Marshal(event)
		manager.broker.Publish(payload, topic)
	}
}

func (manager *Manager) watchInput(input interface{}, topic string) {
	channel := reflect.ValueOf(input)

	for value, ok := channel.Recv(); ok; value, ok = channel.Recv() {
		event := value.Interface()
		manager.Publish(event, topic)
	}
}
