package ganglia

import (
	"sync"
)

type AnalogEventObserver interface {
	HandleAnalogEvent(AnalogEvent)
}

type AnalogMonitor struct {
	currentValue AnalogEvent
	mutex        sync.Mutex
	observers    []AnalogEventObserver
	source       AnalogInput
}

func NewAnalogMonitor(source AnalogInput) *AnalogMonitor {
	monitor := new(AnalogMonitor)
	monitor.source = source
	go monitor.watchSource()
	return monitor
}

func (monitor *AnalogMonitor) CurrentValue() AnalogEvent {
	return monitor.currentValue
}

func (monitor *AnalogMonitor) handleEvent(event AnalogEvent) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	for _, observer := range monitor.observers {
		observer.HandleAnalogEvent(event)
	}
}

func (monitor *AnalogMonitor) Subscribe(observer AnalogEventObserver) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()
	monitor.observers = append(monitor.observers, observer)
}

func (monitor *AnalogMonitor) Unsubscribe(observer AnalogEventObserver) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	var observers []AnalogEventObserver

	for _, existing := range monitor.observers {
		if existing != observer {
			observers = append(observers, existing)
		}
	}

	monitor.observers = observers
}

func (monitor *AnalogMonitor) watchSource() {
	for event := range monitor.source {
		monitor.currentValue = event
		monitor.handleEvent(event)
	}
}

type DigitalEventObserver interface {
	HandleDigitalEvent(DigitalEvent)
}

type DigitalMonitor struct {
	currentValue DigitalEvent
	mutex        sync.Mutex
	observers    []DigitalEventObserver
	source       DigitalInput
}

func NewDigitalMonitor(source DigitalInput) *DigitalMonitor {
	monitor := new(DigitalMonitor)
	monitor.source = source
	go monitor.watchSource()
	return monitor
}

func (monitor *DigitalMonitor) CurrentValue() DigitalEvent {
	return monitor.currentValue
}

func (monitor *DigitalMonitor) handleEvent(event DigitalEvent) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	for _, observer := range monitor.observers {
		observer.HandleDigitalEvent(event)
	}
}

func (monitor *DigitalMonitor) Subscribe(observer DigitalEventObserver) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()
	monitor.observers = append(monitor.observers, observer)
}

func (monitor *DigitalMonitor) Unsubscribe(observer DigitalEventObserver) {
	monitor.mutex.Lock()
	defer monitor.mutex.Unlock()

	var observers []DigitalEventObserver

	for _, existing := range monitor.observers {
		if existing != observer {
			observers = append(observers, existing)
		}
	}

	monitor.observers = observers
}

func (monitor *DigitalMonitor) watchSource() {
	for event := range monitor.source {
		monitor.currentValue = event
		monitor.handleEvent(event)
	}
}
