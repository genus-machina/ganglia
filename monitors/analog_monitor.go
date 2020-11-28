package monitors

import (
	"sync"

	"github.com/genus-machina/ganglia"
)

type AnalogEventHandler func(*ganglia.AnalogEvent)

type AnalogEventObserver struct {
	handler AnalogEventHandler
}

func NewAnalogEventObserver(handler AnalogEventHandler) *AnalogEventObserver {
	observer := new(AnalogEventObserver)
	observer.handler = handler
	return observer
}

type AnalogMonitor interface {
	CurrentValue() *ganglia.AnalogEvent
	Once(*AnalogEventObserver) ganglia.Trigger
	Subscribe(*AnalogEventObserver) ganglia.Trigger
	Unsubscribe(*AnalogEventObserver)
}

type AnalogInputMonitor struct {
	currentValue *ganglia.AnalogEvent
	mutex        sync.Mutex
	observers    []*AnalogEventObserver
	source       ganglia.AnalogInput
}

func NewAnalogInputMonitor(source ganglia.AnalogInput) *AnalogInputMonitor {
	monitor := new(AnalogInputMonitor)
	monitor.source = source
	go monitor.watchSource()
	return monitor
}

func (monitor *AnalogInputMonitor) CurrentValue() *ganglia.AnalogEvent {
	return monitor.currentValue
}

func (monitor *AnalogInputMonitor) getObservers() []*AnalogEventObserver {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	var observers []*AnalogEventObserver

	for _, observer := range monitor.observers {
		observers = append(observers, observer)
	}

	return observers
}

func (monitor *AnalogInputMonitor) handleEvent(event *ganglia.AnalogEvent) {
	for _, observer := range monitor.getObservers() {
		observer.handler(event)
	}
}

func (monitor *AnalogInputMonitor) Once(observer *AnalogEventObserver) ganglia.Trigger {
	var wrapped *AnalogEventObserver

	handler := func(event *ganglia.AnalogEvent) {
		observer.handler(event)
		monitor.Unsubscribe(wrapped)
	}

	wrapped = NewAnalogEventObserver(handler)
	monitor.Subscribe(wrapped)

	return func() {
		monitor.Unsubscribe(observer)
	}
}

func (monitor *AnalogInputMonitor) Subscribe(observer *AnalogEventObserver) ganglia.Trigger {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()
	monitor.observers = append(monitor.observers, observer)
	return monitor.triggerUnsubscribe(observer)
}

func (monitor *AnalogInputMonitor) triggerUnsubscribe(observer *AnalogEventObserver) ganglia.Trigger {
	return func() {
		monitor.Unsubscribe(observer)
	}
}

func (monitor *AnalogInputMonitor) Unsubscribe(observer *AnalogEventObserver) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	var observers []*AnalogEventObserver

	for _, existing := range monitor.observers {
		if existing != observer {
			observers = append(observers, existing)
		}
	}

	monitor.observers = observers
}

func (monitor *AnalogInputMonitor) watchSource() {
	for event := range monitor.source {
		monitor.currentValue = event
		monitor.handleEvent(event)
	}
}
