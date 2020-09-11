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

type DigitalMonitor interface {
	CurrentValue() *ganglia.DigitalEvent
	Subscribe(*DigitalEventObserver)
	Unsubscribe(*DigitalEventObserver)
}

type DigitalInputMonitor struct {
	currentValue *ganglia.DigitalEvent
	mutex        sync.Mutex
	observers    []*DigitalEventObserver
	source       ganglia.DigitalInput
}

func NewDigitalInputMonitor(source ganglia.DigitalInput) *DigitalInputMonitor {
	monitor := new(DigitalInputMonitor)
	monitor.source = source
	go monitor.watchSource()
	return monitor
}

func (monitor *DigitalInputMonitor) CurrentValue() *ganglia.DigitalEvent {
	return monitor.currentValue
}

func (monitor *DigitalInputMonitor) handleEvent(event *ganglia.DigitalEvent) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	for _, observer := range monitor.observers {
		observer.handler(event)
	}
}

func (monitor *DigitalInputMonitor) Subscribe(observer *DigitalEventObserver) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()
	monitor.observers = append(monitor.observers, observer)
}

func (monitor *DigitalInputMonitor) Unsubscribe(observer *DigitalEventObserver) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	var observers []*DigitalEventObserver

	for _, existing := range monitor.observers {
		if existing != observer {
			observers = append(observers, existing)
		}
	}

	monitor.observers = observers
}

func (monitor *DigitalInputMonitor) watchSource() {
	for event := range monitor.source {
		monitor.currentValue = event
		monitor.handleEvent(event)
	}
}
