package i2c

import (
	"log"
	"sync"
	"time"

	"github.com/genus-machina/ganglia"
	"periph.io/x/periph/conn/i2c"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/devices/bmxx80"
)

type BME280 struct {
	bus    i2c.Bus
	device *bmxx80.Dev
	done   chan bool
	logger *log.Logger
	mutex  sync.Mutex
}

func newBME280(logger *log.Logger, bus i2c.Bus) *BME280 {
	device := new(BME280)
	device.bus = bus
	device.logger = log.New(logger.Writer(), "[bme280] ", logger.Flags())
	return device
}

func (device *BME280) buildDevice() (*bmxx80.Dev, error) {
	options := &bmxx80.Opts{
		Humidity:    bmxx80.O1x,
		Pressure:    bmxx80.O1x,
		Temperature: bmxx80.O1x,
	}

	return bmxx80.NewI2C(device.bus, 0x76, options)
}

func (device *BME280) buildEvent(value *physic.Env) *ganglia.Environmental {
	event := new(ganglia.Environmental)
	event.Humidity = float64(value.Humidity) / float64(physic.PercentRH)
	event.Pressure = float64(value.Pressure) / float64(physic.Pascal) / 3386.38816
	event.Temperature = value.Temperature.Fahrenheit()
	event.Time = time.Now()
	return event
}

func (device *BME280) Halt() error {
	close(device.done)
	return device.halt()
}

func (device *BME280) halt() error {
	return device.device.Halt()
}

func (device *BME280) readValue() (*physic.Env, error) {
	device.mutex.Lock()
	defer device.mutex.Unlock()

	var err error
	value := new(physic.Env)
	err = device.device.Sense(value)
	return value, err
}

func (device *BME280) readValues(output ganglia.EnvironmentalOutput, interval time.Duration) {
	defer close(output)

	for reading := true; reading; {
		if value, err := device.readValue(); err == nil {
			event := device.buildEvent(value)
			output <- event
		} else {
			device.logger.Printf("BME280 failed. Attempting to rebuild. %s\n", err.Error())
			device.rebuildDevice()
		}

		select {
		case <-time.After(interval):
		case <-device.done:
			reading = false
		}
	}
}

func (device *BME280) rebuildDevice() {
	device.mutex.Lock()
	defer device.mutex.Unlock()

	device.halt()
	device.resetDevice()

	var err error
	if device.device, err = device.buildDevice(); err != nil {
		device.logger.Fatalf("Failed to rebuild BME280. %s\n", err.Error())
	}
}

func (device *BME280) resetDevice() error {
	address := uint16(0xE0)
	command := []byte{0xB6}
	device.device.Halt()
	return device.bus.Tx(address, command, nil)
}

func (device *BME280) SenseContinuous(interval time.Duration) (ganglia.EnvironmentalInput, error) {
	if err := device.start(); err != nil {
		return nil, err
	}

	values := make(chan *ganglia.Environmental)
	go device.readValues(values, interval)
	return values, nil
}

func (device *BME280) start() error {
	device.mutex.Lock()
	defer device.mutex.Unlock()

	var err error
	if device.device == nil {
		device.done = make(chan bool)
		device.device, err = device.buildDevice()
	}

	return err
}
