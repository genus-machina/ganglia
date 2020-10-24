package gpio

import (
	"log"
	"time"

	"github.com/genus-machina/ganglia"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
)

const (
	PullDown Pull = Pull(gpio.PullDown)
	PullUp   Pull = Pull(gpio.PullUp)
)

type Manager struct {
	done   chan bool
	logger *log.Logger
	pins   []gpio.PinIO
}

type Pull gpio.Pull

func NewManager(logger *log.Logger) *Manager {
	manager := new(Manager)
	manager.done = make(chan bool, 0)
	manager.logger = log.New(logger.Writer(), "[gpio] ", logger.Flags())

	if _, err := host.Init(); err != nil {
		manager.logger.Fatalf("Failed to initialize host. %s.\n", err.Error())
	}

	return nil
}

func (manager *Manager) Halt() {
	close(manager.done)

	for _, pin := range manager.pins {
		if err := pin.Halt(); err != nil {
			manager.logger.Printf("Failed to halt pin '%s'. %s\n", pin.Name(), err.Error())
		}
	}
}

func (manager *Manager) Input(name string, pull Pull) (ganglia.DigitalInput, error) {
	input := make(chan *ganglia.DigitalEvent, 1)
	pin := gpioreg.ByName(name)

	if err := pin.In(gpio.Pull(pull), gpio.BothEdges); err != nil {
		return nil, err
	}

	manager.pins = append(manager.pins, pin)
	go manager.watchInput(pin, input)
	return input, nil
}

func (manager *Manager) Output(name string, initialValue ganglia.DigitalValue) (ganglia.DigitalOutput, error) {
	output := make(chan ganglia.DigitalValue, 0)
	pin := gpioreg.ByName(name)

	if err := pin.Out(gpio.Level(initialValue)); err != nil {
		return nil, err
	}

	manager.pins = append(manager.pins, pin)
	go manager.watchOutput(output, pin)
	return output, nil
}

func (manager *Manager) watchInput(pin gpio.PinIn, input chan *ganglia.DigitalEvent) {
	defer close(input)

	for running := true; running; {
		select {
		case <-manager.done:
			running = false
		default:
			if pin.WaitForEdge(-1) {
				event := &ganglia.DigitalEvent{
					Time:  time.Now(),
					Value: ganglia.DigitalValue(pin.Read()),
				}

				input <- event
			}
		}
	}
}

func (manager *Manager) watchOutput(output chan ganglia.DigitalValue, pin gpio.PinOut) {
	for value := range output {
		if err := pin.Out(gpio.Level(value)); err != nil {
			manager.logger.Printf(
				"Failed to write value to pin '%s'. %s\n",
				pin.Name(),
				err.Error(),
			)
		}
	}
}
