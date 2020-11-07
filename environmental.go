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

func (input EnvironmentalInput) CalibrateTemperature(coefficient, intercept float64) EnvironmentalInput {
	calibrated := make(chan *EnvironmentalEvent, 1)
	go input.calibrateTemperature(calibrated, coefficient, intercept)
	return calibrated
}

func (input EnvironmentalInput) calibrateTemperature(calibrated chan<- *EnvironmentalEvent, coefficient, intercept float64) {
	defer close(calibrated)

	for event := range input {
		calibrated <- &EnvironmentalEvent{
			Humidity:    event.Humidity,
			Pressure:    event.Pressure,
			Temperature: event.Temperature*coefficient + intercept,
			Time:        event.Time,
		}
	}
}

func (input EnvironmentalInput) Read() *EnvironmentalEvent {
	value := <-input
	return value
}

type EnvironmentalOutput chan<- *EnvironmentalEvent

func (output EnvironmentalOutput) Write(value *EnvironmentalEvent) {
	output <- value
}
