package i2c

import (
	"log"
	"time"

	"github.com/genus-machina/ganglia"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/devices/ssd1306"
	"periph.io/x/periph/host"
)

type Device interface {
	Halt() error
}

type Manager struct {
	bme280  *BME280
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
	if _, err := host.Init(); err != nil {
		return nil, err
	}

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

func (manager *Manager) EnvironmentalInput() (ganglia.EnvironmentalInput, error) {
	if manager.bme280 == nil {
		manager.bme280 = newBME280(manager.logger, manager.bus)
		manager.devices = append(manager.devices, manager.bme280)
	}

	input, err := manager.bme280.SenseContinuous(time.Minute)
	if err != nil {
		return nil, err
	}

	return input, nil
}
