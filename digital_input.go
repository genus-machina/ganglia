package ganglia

import (
	"time"
)

const (
	High DigitalValue = true
	Low  DigitalValue = false
)

type DigitalEvent struct {
	Time  time.Time    `json:"time"`
	Value DigitalValue `json:"value"`
}

type DigitalInput <-chan *DigitalEvent

func (input DigitalInput) Debounce(duration time.Duration) DigitalInput {
	debounced := make(chan *DigitalEvent)
	go input.debounce(debounced, duration)
	return debounced
}

func (input DigitalInput) debounce(debounced chan<- *DigitalEvent, duration time.Duration) {
	var next time.Time
	defer close(debounced)

	for event := range input {
		if now := time.Now(); now.After(next) {
			next = now.Add(duration)
			debounced <- event
		}
	}
}

func (input DigitalInput) Invert() DigitalInput {
	inverted := make(chan *DigitalEvent)
	go input.invert(inverted)
	return inverted
}

func (input DigitalInput) invert(inverted chan<- *DigitalEvent) {
	defer close(inverted)

	for event := range input {
		inverted <- &DigitalEvent{event.Time, !event.Value}
	}
}

func (input DigitalInput) Read() *DigitalEvent {
	value := <-input
	return value
}

type DigitalValue bool
