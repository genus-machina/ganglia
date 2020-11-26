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
	Once(*AnalogEventObserver)
	Subscribe(*AnalogEventObserver)
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

func (notifier *AnalogInputMonitor) Once(observer *AnalogEventObserver) {
	var wrapped *AnalogEventObserver

	handler := func(event *ganglia.AnalogEvent) {
		observer.handler(event)
		notifier.Unsubscribe(wrapped)
	}

	wrapped = NewAnalogEventObserver(handler)
	notifier.Subscribe(wrapped)
}

func (monitor *AnalogInputMonitor) Subscribe(observer *AnalogEventObserver) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()
	monitor.observers = append(monitor.observers, observer)
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
