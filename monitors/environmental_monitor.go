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

type EnivormentalMonitor interface {
	CurrentValue() *ganglia.EnvironmentalEvent
	Subscribe(*EnvironmentalEventObserver)
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

func (monitor *EnvironmentalInputMonitor) handleEvent(event *ganglia.EnvironmentalEvent) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	for _, observer := range monitor.observers {
		observer.handler(event)
	}
}

func (monitor *EnvironmentalInputMonitor) Subscribe(observer *EnvironmentalEventObserver) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()
	monitor.observers = append(monitor.observers, observer)
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
