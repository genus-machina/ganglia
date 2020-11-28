package monitors

import (
	"sync"

	"github.com/genus-machina/ganglia"
)

type EnvironmentalEventHandler func(*ganglia.EnvironmentalEvent)

type EnvironmentalEventObserver struct {
	handler EnvironmentalEventHandler
}

func NewEnvironmentalEventObserver(handler EnvironmentalEventHandler) *EnvironmentalEventObserver {
	observer := new(EnvironmentalEventObserver)
	observer.handler = handler
	return observer
}

func NewEnvironmentalForwarder(output ganglia.EnvironmentalOutput) *EnvironmentalEventObserver {
	return NewEnvironmentalEventObserver(
		func(event *ganglia.EnvironmentalEvent) {
			output <- event
		},
	)
}

type EnvironmentalMonitor interface {
	CurrentValue() *ganglia.EnvironmentalEvent
	Once(*EnvironmentalEventObserver) ganglia.Trigger
	Subscribe(*EnvironmentalEventObserver) ganglia.Trigger
	Unsubscribe(*EnvironmentalEventObserver)
}

type EnvironmentalInputMonitor struct {
	currentValue *ganglia.EnvironmentalEvent
	mutex        sync.Mutex
	observers    []*EnvironmentalEventObserver
	source       ganglia.EnvironmentalInput
}

func NewEnvironmentalInputMonitor(source ganglia.EnvironmentalInput) *EnvironmentalInputMonitor {
	monitor := new(EnvironmentalInputMonitor)
	monitor.source = source
	go monitor.watchSource()
	return monitor
}

func (monitor *EnvironmentalInputMonitor) CurrentValue() *ganglia.EnvironmentalEvent {
	return monitor.currentValue
}

func (monitor *EnvironmentalInputMonitor) getObservers() []*EnvironmentalEventObserver {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	var observers []*EnvironmentalEventObserver

	for _, observer := range monitor.observers {
		observers = append(observers, observer)
	}

	return observers
}

func (monitor *EnvironmentalInputMonitor) handleEvent(event *ganglia.EnvironmentalEvent) {
	for _, observer := range monitor.getObservers() {
		observer.handler(event)
	}
}

func (monitor *EnvironmentalInputMonitor) Once(observer *EnvironmentalEventObserver) ganglia.Trigger {
	var wrapped *EnvironmentalEventObserver

	handler := func(event *ganglia.EnvironmentalEvent) {
		observer.handler(event)
		monitor.Unsubscribe(wrapped)
	}

	wrapped = NewEnvironmentalEventObserver(handler)
	monitor.Subscribe(wrapped)
	return monitor.triggerUnsubscribe(observer)
}

func (monitor *EnvironmentalInputMonitor) Subscribe(observer *EnvironmentalEventObserver) ganglia.Trigger {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()
	monitor.observers = append(monitor.observers, observer)
	return monitor.triggerUnsubscribe(observer)
}

func (monitor *EnvironmentalInputMonitor) triggerUnsubscribe(observer *EnvironmentalEventObserver) ganglia.Trigger {
	return func() {
		monitor.Unsubscribe(observer)
	}
}

func (monitor *EnvironmentalInputMonitor) Unsubscribe(observer *EnvironmentalEventObserver) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	var observers []*EnvironmentalEventObserver

	for _, existing := range monitor.observers {
		if existing != observer {
			observers = append(observers, existing)
		}
	}

	monitor.observers = observers
}

func (monitor *EnvironmentalInputMonitor) watchSource() {
	for event := range monitor.source {
		monitor.currentValue = event
		monitor.handleEvent(event)
	}
}
