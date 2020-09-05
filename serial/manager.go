package serial

import (
	"log"
	"time"

	"github.com/genus-machina/ganglia"
)

type AnalogEvent struct {
	time  time.Time
	value uint
}

func (event *AnalogEvent) Time() time.Time {
	return event.time
}

func (event *AnalogEvent) Value() ganglia.AnalogValue {
	return ganglia.AnalogValue(event.value)
}

type DigitalEvent struct {
	time  time.Time
	value bool
}

func (event *DigitalEvent) Time() time.Time {
	return event.time
}

func (event *DigitalEvent) Value() ganglia.DigitalValue {
	return ganglia.DigitalValue(event.value)
}

type Manager struct {
	broker *Broker
	done   chan bool
	logger *log.Logger
}

func NewManager(logger *log.Logger, broker *Broker) *Manager {
	manager := new(Manager)
	manager.broker = broker
	manager.done = make(chan bool)
	manager.logger = log.New(logger.Writer(), "[serial] ", logger.Flags())
	return manager
}

func Connect(logger *log.Logger, port string) (*Manager, error) {
	broker, err := Open(port)

	if err != nil {
		return nil, err
	}

	return NewManager(logger, broker), nil
}

func (manager *Manager) AnalogInput(pin int, interval time.Duration) ganglia.AnalogInput {
	input := make(chan ganglia.AnalogEvent)
	go manager.watchAnalogPin(pin, interval, input)
	return input
}

func (manager *Manager) DigitalInput(pin int, interval time.Duration) ganglia.DigitalInput {
	input := make(chan ganglia.DigitalEvent)
	go manager.watchDigitalPin(pin, interval, input)
	return input
}

func (manager *Manager) DigitalOutput(pin int) ganglia.DigitalOutput {
	output := make(chan ganglia.DigitalValue)
	go manager.watchDigitalOutput(output, pin)
	return output
}

func (manager *Manager) Halt() {
	if err := manager.broker.Close(); err != nil {
		manager.logger.Printf("Failed to close broker. %s\n", err.Error())
	}
}

func (manager *Manager) watchAnalogPin(pin int, interval time.Duration, input chan ganglia.AnalogEvent) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(input)

	for running := true; running; {
		select {
		case <-manager.done:
			running = false
		case <-ticker.C:
			if value, err := manager.broker.ReadAnalogValue(pin); err == nil {
				event := &AnalogEvent{time.Now(), value}
				input <- event
			} else {
				manager.logger.Printf(
					"Failed to read analog value from pin %d. %s\n",
					pin,
					err.Error(),
				)
			}
		}
	}
}

func (manager *Manager) watchDigitalOutput(output chan ganglia.DigitalValue, pin int) {
	defer close(output)

	for running := true; running; {
		select {
		case <-manager.done:
			running = false
		case value := <-output:
			if err := manager.broker.WriteDigitalValue(pin, bool(value)); err != nil {
				manager.logger.Printf(
					"Failed to write digital value to pin %d. %s\n",
					pin,
					err.Error(),
				)
			}
		}
	}
}

func (manager *Manager) watchDigitalPin(pin int, interval time.Duration, input chan ganglia.DigitalEvent) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	defer close(input)

	for running := true; running; {
		select {
		case <-manager.done:
			running = false
		case <-ticker.C:
			if value, err := manager.broker.ReadDigitalValue(pin); err == nil {
				event := &DigitalEvent{time.Now(), value}
				input <- event
			} else {
				manager.logger.Printf(
					"Failed to read digital value from pin %d. %s\n",
					pin,
					err.Error(),
				)
			}
		}
	}
}
