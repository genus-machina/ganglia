package ganglia

import (
	"time"
)

type Environmental struct {
	Humidity    float64   `json:"humidity,omitempty"`    // percent
	Pressure    float64   `json:"pressure,omitempty"`    // inHg
	Temperature float64   `json:"temperature,omitempty"` // fahrenheit
	Time        time.Time `json:"time"`
}

type EnvironmentalInput <-chan *Environmental

func (input EnvironmentalInput) Read() *Environmental {
	value := <-input
	return value
}

type EnvironmentalOutput chan<- *Environmental

func (output EnvironmentalOutput) Write(value *Environmental) {
	output <- value
}
