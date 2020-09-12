package monitors

import (
	"sync"

	"github.com/genus-machina/ganglia"
)

type DigitalEventHandler func(*ganglia.DigitalEvent)

type DigitalEventObserver struct {
	handler DigitalEventHandler
}

func NewDigitalEventObserver(handler DigitalEventHandler) *DigitalEventObserver {
	observer := new(DigitalEventObserver)
	observer.handler = handler
	return observer
}

type digitalNotifier struct {
	mutex     sync.Mutex
	observers []*DigitalEventObserver
}

func (notifier *digitalNotifier) handleEvent(event *ganglia.DigitalEvent) {
	notifier.mutex.Lock()
	defer notifier.mutex.Unlock()

	for _, observer := range notifier.observers {
		observer.handler(event)
	}
}

func (notifier *digitalNotifier) Subscribe(observer *DigitalEventObserver) {
	notifier.mutex.Lock()
	defer notifier.mutex.Unlock()
	notifier.observers = append(notifier.observers, observer)
}

func (notifier *digitalNotifier) Unsubscribe(observer *DigitalEventObserver) {
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
