package monitors

import (
	"sync"

	"github.com/genus-machina/ganglia"
)

type DigitalEventHandler func(*ganglia.DigitalEvent)

type DigitalEventObserver struct {
	Handler DigitalEventHandler
}

func NewDigitalEventObserver(handler DigitalEventHandler) *DigitalEventObserver {
	observer := new(DigitalEventObserver)
	observer.Handler = handler
	return observer
}

type DigitalTarget interface {
	Write(ganglia.DigitalValue)
}

func NewDigitalForwarder(output DigitalTarget) *DigitalEventObserver {
	return NewDigitalEventObserver(
		func(event *ganglia.DigitalEvent) {
			output.Write(event.Value)
		},
	)
}

func NewDigitalHandler(handler ganglia.Trigger) *DigitalEventObserver {
	return NewDigitalEventObserver(
		func(event *ganglia.DigitalEvent) {
			handler()
		},
	)
}

func NewDigitalTrigger(trigger ganglia.Trigger) *DigitalEventObserver {
	return NewDigitalEventObserver(
		func(event *ganglia.DigitalEvent) {
			if event.Value == ganglia.High {
				trigger()
			}
		},
	)
}

type DigitalNotifier struct {
	mutex     sync.Mutex
	observers []*DigitalEventObserver
}

func (notifier *DigitalNotifier) getObservers() []*DigitalEventObserver {
	notifier.mutex.Lock()
	defer notifier.mutex.Unlock()

	var observers []*DigitalEventObserver

	for _, observer := range notifier.observers {
		observers = append(observers, observer)
	}

	return observers
}

func (notifier *DigitalNotifier) Notify(event *ganglia.DigitalEvent) {
	for _, observer := range notifier.getObservers() {
		observer.Handler(event)
	}
}

func (notifier *DigitalNotifier) Once(observer *DigitalEventObserver) ganglia.Trigger {
	var wrapped *DigitalEventObserver

	handler := func(event *ganglia.DigitalEvent) {
		observer.Handler(event)
		notifier.Unsubscribe(wrapped)
	}

	wrapped = NewDigitalEventObserver(handler)
	notifier.Subscribe(wrapped)
	return notifier.triggerUnsubscribe(wrapped)
}

func (notifier *DigitalNotifier) Subscribe(observer *DigitalEventObserver) ganglia.Trigger {
	notifier.mutex.Lock()
	defer notifier.mutex.Unlock()
	notifier.observers = append(notifier.observers, observer)
	return notifier.triggerUnsubscribe(observer)
}

func (notifier *DigitalNotifier) triggerUnsubscribe(observer *DigitalEventObserver) ganglia.Trigger {
	return func() {
		notifier.Unsubscribe(observer)
	}
}

func (notifier *DigitalNotifier) Unsubscribe(observer *DigitalEventObserver) {
	notifier.mutex.Lock()
	defer notifier.mutex.Unlock()

	var observers []*DigitalEventObserver

	for _, existing := range notifier.observers {
		if existing != observer {
			observers = append(observers, existing)
		}
	}

	notifier.observers = observers
}
