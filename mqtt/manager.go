package mqtt

import (
	"encoding/json"
	"log"
	"time"

	"github.com/genus-machina/ganglia"
)

type Manager struct {
	broker *Broker
	logger *log.Logger

	analogChannels        []chan *ganglia.AnalogEvent
	environmentalChannels []chan *ganglia.EnvironmentalEvent
	digitalChannels       []chan *ganglia.DigitalEvent
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
	manager.broker.Subscribe(topic, manager.handleAnalogEvent(input))
	manager.analogChannels = append(manager.analogChannels, input)
	return input
}

func (manager *Manager) DigitalInput(topic string) ganglia.DigitalInput {
	input := make(chan *ganglia.DigitalEvent, 1)
	manager.broker.Subscribe(topic, manager.handleDigitalEvent(input))
	manager.digitalChannels = append(manager.digitalChannels, input)
	return input
}

func (manager *Manager) DigitalOutput(topic string) ganglia.DigitalOutput {
	output := make(chan ganglia.DigitalValue)
	go manager.watchDigitalOutput(output, topic)
	return output
}

func (manager *Manager) EnvironmentalInput(topic string) ganglia.EnvironmentalInput {
	input := make(chan *ganglia.EnvironmentalEvent, 1)
	manager.broker.Subscribe(topic, manager.handleEnvironmentalEvent(input))
	manager.environmentalChannels = append(manager.environmentalChannels, input)
	return input
}

func (manager *Manager) Halt() {
	manager.broker.Close()

	for _, channel := range manager.analogChannels {
		close(channel)
	}

	for _, channel := range manager.digitalChannels {
		close(channel)
	}

	for _, channel := range manager.environmentalChannels {
		close(channel)
	}
}

func (manager *Manager) handleAnalogEvent(input chan *ganglia.AnalogEvent) MessageHandler {
	return func(message Message) {
		event := new(ganglia.AnalogEvent)

		if err := json.Unmarshal(message, event); err == nil {
			input <- event
		} else {
			manager.logger.Printf("Failed to parse analog event. %s.\n", err.Error())
		}
	}
}

func (manager *Manager) handleDigitalEvent(input chan *ganglia.DigitalEvent) MessageHandler {
	return func(message Message) {
		event := new(ganglia.DigitalEvent)

		if err := json.Unmarshal(message, event); err == nil {
			input <- event
		} else {
			manager.logger.Printf("Failed to parse digital event. %s.\n", err.Error())
		}
	}
}

func (manager *Manager) handleEnvironmentalEvent(input chan *ganglia.EnvironmentalEvent) MessageHandler {
	return func(message Message) {
		event := new(ganglia.EnvironmentalEvent)

		if err := json.Unmarshal(message, event); err == nil {
			input <- event
		} else {
			manager.logger.Printf("Failed to parse environmental event. %s.\n", err.Error())
		}
	}
}

func (manager *Manager) PublishAnalogInput(input ganglia.AnalogInput, topic string) {
	go manager.watchAnalogInput(input, topic)
}

func (manager *Manager) PublishDigitalInput(input ganglia.DigitalInput, topic string) {
	go manager.watchDigitalInput(input, topic)
}

func (manager *Manager) PublishEnvironmentalInput(input ganglia.EnvironmentalInput, topic string) {
	go manager.watchEnvironmentalInput(input, topic)
}

func (manager *Manager) watchAnalogInput(input ganglia.AnalogInput, topic string) {
	for event := range input {
		payload, _ := json.Marshal(event)
		manager.broker.Publish(payload, topic)
	}
}

func (manager *Manager) watchDigitalInput(input ganglia.DigitalInput, topic string) {
	for event := range input {
		payload, _ := json.Marshal(event)
		manager.broker.Publish(payload, topic)
	}
}

func (manager *Manager) watchDigitalOutput(output chan ganglia.DigitalValue, topic string) {
	for value := range output {
		event := &ganglia.DigitalEvent{time.Now(), value}
		payload, _ := json.Marshal(event)
		manager.broker.Publish(payload, topic)
	}
}

func (manager *Manager) watchEnvironmentalInput(input ganglia.EnvironmentalInput, topic string) {
	for event := range input {
		payload, _ := json.Marshal(event)
		manager.broker.Publish(payload, topic)
	}
}
