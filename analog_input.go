package ganglia

import (
	"time"
)

type AnalogValue uint

type AnalogEvent struct {
	Time  time.Time   `json:"time"`
	Value AnalogValue `json:"value"`
}

type AnalogInput <-chan *AnalogEvent

func (input AnalogInput) Read() *AnalogEvent {
	value := <-input
	return value
}
