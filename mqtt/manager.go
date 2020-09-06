package mqtt

import (
	"log"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/genus-machina/ganglia"
)

type Manager struct {
	broker *Broker
	logger *log.Logger

	analogChannels  []chan ganglia.AnalogEvent
	digitalChannels []chan ganglia.DigitalEvent
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
	input := make(chan ganglia.AnalogEvent, 1)
	manager.broker.Subscribe(topic, manager.handleAnalogEvent(input))
	manager.analogChannels = append(manager.analogChannels, input)
	return input
}

func (manager *Manager) DigitalInput(topic string) ganglia.DigitalInput {
	input := make(chan ganglia.DigitalEvent, 1)
	manager.broker.Subscribe(topic, manager.handleDigitalEvent(input))
	manager.digitalChannels = append(manager.digitalChannels, input)
	return input
}

func (manager *Manager) DigitalOutput(topic string) ganglia.DigitalOutput {
	output := make(chan ganglia.DigitalValue)
	go manager.watchOutput(output, topic)
	return output
}

func (manager *Manager) Halt() {
	manager.broker.Close()

	for _, channel := range manager.analogChannels {
		close(channel)
	}

	for _, channel := range manager.digitalChannels {
		close(channel)
	}
}

func (manager *Manager) handleAnalogEvent(input chan ganglia.AnalogEvent) mqtt.MessageHandler {
	return func(client mqtt.Client, message mqtt.Message) {
		if event, err := parseAnalogEvent(message); err == nil {
			input <- event
		} else {
			manager.logger.Printf("Failed to parse analog event. %s.\n", err.Error())
		}

		message.Ack()
	}
}

func (manager *Manager) handleDigitalEvent(input chan ganglia.DigitalEvent) mqtt.MessageHandler {
	return func(client mqtt.Client, message mqtt.Message) {
		if event, err := parseDigitalEvent(message); err == nil {
			input <- event
		} else {
			manager.logger.Printf("Failed to parse digital event. %s.\n", err.Error())
		}

		message.Ack()
	}
}

func (manager *Manager) PublishAnalogInput(input ganglia.AnalogInput, topic string) {
	go manager.watchAnalogInput(input, topic)
}

func (manager *Manager) PublishDigitalInput(input ganglia.DigitalInput, topic string) {
	go manager.watchDigitalInput(input, topic)
}

func (manager *Manager) watchAnalogInput(input ganglia.AnalogInput, topic string) {
	for event := range input {
		manager.broker.Publish(encodeAnalogEvent(event), topic)
	}
}

func (manager *Manager) watchDigitalInput(input ganglia.DigitalInput, topic string) {
	for event := range input {
		manager.broker.Publish(encodeDigitalEvent(event), topic)
	}
}

func (manager *Manager) watchOutput(output chan ganglia.DigitalValue, topic string) {
	for value := range output {
		event := &DigitalEvent{bool(value), time.Now()}
		manager.broker.Publish(encodeDigitalEvent(event), topic)
	}
}
