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

func (monitor *AnalogInputMonitor) handleEvent(event *ganglia.AnalogEvent) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	for _, observer := range monitor.observers {
		observer.handler(event)
	}
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
