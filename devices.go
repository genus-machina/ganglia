package ganglia

import (
	"time"
)

const (
	High DigitalValue = true
	Low  DigitalValue = false
)

type AnalogValue uint
type DigitalValue bool

type AnalogEvent interface {
	Time() time.Time
	Value() AnalogValue
}

type DigitalEvent interface {
	Time() time.Time
	Value() DigitalValue
}

type AnalogInput <-chan AnalogEvent

func (input AnalogInput) Read() AnalogEvent {
	value := <-input
	return value
}

type DigitalInput <-chan DigitalEvent

func (input DigitalInput) Read() DigitalEvent {
	value := <-input
	return value
}

type DigitalOutput chan<- DigitalValue

func (output DigitalOutput) Write(value DigitalValue) {
	output <- value
}
