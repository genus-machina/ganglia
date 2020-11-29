package monitors

import (
	"time"

	"github.com/genus-machina/ganglia"
)

type DepressMonitor struct {
	digitalNotifier
	current, last *ganglia.DigitalEvent
	observer      *DigitalEventObserver
	source        DigitalMonitor
}

func NewDepressMonitor(source DigitalMonitor) *DepressMonitor {
	monitor := new(DepressMonitor)
	monitor.observer = NewDigitalEventObserver(monitor.handleEvent)
	monitor.source = source
	return monitor
}

func (monitor *DepressMonitor) CurrentValue() *ganglia.DigitalEvent {
	return monitor.current
}

func (monitor *DepressMonitor) handleEvent(event *ganglia.DigitalEvent) {
	isLow := event.Value == ganglia.Low
	wasHigh := monitor.last != nil && monitor.last.Value == ganglia.High

	if wasHigh && isLow && event.Time.Before(monitor.last.Time.Add(500*time.Millisecond)) {
		monitor.update(monitor.last)
		monitor.update(event)
	}

	monitor.last = event
}

func (monitor *DepressMonitor) Subscribe(observer *DigitalEventObserver) ganglia.Trigger {
	if len(monitor.digitalNotifier.observers) == 0 {
		monitor.source.Subscribe(monitor.observer)
	}

	monitor.digitalNotifier.Subscribe(observer)
	return monitor.triggerUnsubscribe(observer)
}

func (monitor *DepressMonitor) triggerUnsubscribe(observer *DigitalEventObserver) ganglia.Trigger {
	return func() {
		monitor.Unsubscribe(observer)
	}
}

func (monitor *DepressMonitor) Unsubscribe(observer *DigitalEventObserver) {
	monitor.digitalNotifier.Unsubscribe(observer)

	if len(monitor.digitalNotifier.observers) == 0 {
		monitor.source.Unsubscribe(monitor.observer)
	}
}

func (monitor *DepressMonitor) update(event *ganglia.DigitalEvent) {
	monitor.current = event
	monitor.digitalNotifier.handleEvent(monitor.current)
}
