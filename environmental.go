package ganglia

import (
	"time"
)

type EnvironmentalEvent struct {
	Humidity    float64   `json:"humidity,omitempty"`    // percent
	Pressure    float64   `json:"pressure,omitempty"`    // inHg
	Temperature float64   `json:"temperature,omitempty"` // fahrenheit
	Time        time.Time `json:"time"`
}

type EnvironmentalInput <-chan *EnvironmentalEvent

func (input EnvironmentalInput) Read() *EnvironmentalEvent {
	value := <-input
	return value
}

type EnvironmentalOutput chan<- *EnvironmentalEvent

func (output EnvironmentalOutput) Write(value *EnvironmentalEvent) {
	output <- value
}
