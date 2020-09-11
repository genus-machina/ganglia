package ganglia

import (
	"sync"
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
	var pending sync.WaitGroup
	var timer *time.Timer
	defer close(debounced)

	for event := range input {
		if timer != nil {
			timer.Stop()
		} else {
			pending.Add(1)
		}

		timer = time.AfterFunc(duration, func() {
			debounced <- event
			pending.Done()
		})
	}

	pending.Wait()
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
