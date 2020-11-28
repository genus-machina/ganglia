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
	var timer *time.Timer
	defer close(debounced)

	for event := range input {
		if timer != nil {
			timer.Stop()
		}

		if now := time.Now(); now.After(next) {
			next = now.Add(duration)
			debounced <- event
		} else {
			timer = time.AfterFunc(
				duration,
				func() {
					debounced <- event
				},
			)
		}
	}

	if timer != nil {
		timer.Stop()
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

func (input DigitalInput) Stabilize(duration time.Duration) DigitalInput {
	stabilized := make(chan *DigitalEvent)
	go input.stabilize(stabilized, duration)
	return stabilized
}

func (input DigitalInput) stabilize(stabilized chan<- *DigitalEvent, duration time.Duration) {
	var last *DigitalEvent
	var timer *time.Timer
	defer close(stabilized)

	for event := range input {
		if last == nil || last.Value != event.Value {
			if timer != nil {
				timer.Stop()
			}

			last = event

			timer = time.AfterFunc(
				duration,
				func() {
					stabilized <- event
				},
			)
		}
	}

	if timer != nil {
		timer.Stop()
	}
}

type DigitalValue bool
