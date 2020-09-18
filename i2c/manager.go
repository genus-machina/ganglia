package i2c

import (
	"log"

	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
)

type Device interface {
	Halt() error
}

type Manager struct {
	bus     i2c.BusCloser
	devices []Device
	logger  *log.Logger
}

func newManager(logger *log.Logger, bus i2c.BusCloser) *Manager {
	manager := new(Manager)
	manager.bus = bus
	manager.logger = log.New(logger.Writer(), "[i2c] ", logger.Flags())
	return manager
}

func Connect(logger *log.Logger) (*Manager, error) {
	bus, err := i2creg.Open("")
	if err != nil {
		return nil, err
	}

	return newManager(logger, bus), nil
}

func (manager *Manager) Halt() {
	for _, device := range manager.devices {
		device.Halt()
	}
	manager.bus.Close()
}

func (manager *Manager) Display() (*SSD1306, error) {
	device, err := ssd1306.NewI2C(manager.bus, &ssd1306.DefaultOpts)
	if err != nil {
		return nil, err
	}

	manager.devices = append(manager.devices, device)
	return newSSD1306(device), nil
}
