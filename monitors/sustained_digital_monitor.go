package monitors

import (
	"time"

	"github.com/genus-machina/ganglia"
)

type SustainedDigitalMonitor struct {
	DigitalNotifier
	duration      time.Duration
	current, last *ganglia.DigitalEvent
	observer      *DigitalEventObserver
	source        DigitalMonitor
	timer         *time.Timer
}

func NewSustainedDigitalMonitor(source DigitalMonitor, duration time.Duration) *SustainedDigitalMonitor {
	monitor := new(SustainedDigitalMonitor)
	monitor.duration = duration
	monitor.observer = NewDigitalEventObserver(monitor.handleEvent)
	monitor.source = source
	return monitor
}

func (monitor *SustainedDigitalMonitor) CurrentValue() *ganglia.DigitalEvent {
	if monitor.timer == nil {
		return monitor.source.CurrentValue()
	}
	return monitor.last
}

func (monitor *SustainedDigitalMonitor) handleEvent(event *ganglia.DigitalEvent) {
	monitor.current = event

	if monitor.last == nil || event.Value == ganglia.High {
		if monitor.timer != nil {
			monitor.timer.Stop()
		}

		monitor.timer = time.AfterFunc(monitor.duration, monitor.update)

		monitor.last = event
		monitor.DigitalNotifier.Notify(event)
	}
}

func (monitor *SustainedDigitalMonitor) hasObservers() bool {
	return len(monitor.DigitalNotifier.observers) != 0
}

func (monitor *SustainedDigitalMonitor) Subscribe(observer *DigitalEventObserver) ganglia.Trigger {
	if !monitor.hasObservers() {
		monitor.source.Subscribe(monitor.observer)
	}

	monitor.DigitalNotifier.Subscribe(observer)
	return monitor.triggerUnsubscribe(observer)
}

func (monitor *SustainedDigitalMonitor) triggerUnsubscribe(observer *DigitalEventObserver) ganglia.Trigger {
	return func() {
		monitor.Unsubscribe(observer)
	}
}

func (monitor *SustainedDigitalMonitor) Unsubscribe(observer *DigitalEventObserver) {
	monitor.DigitalNotifier.Unsubscribe(observer)

	if !monitor.hasObservers() {
		monitor.source.Unsubscribe(monitor.observer)
	}
}

func (monitor *SustainedDigitalMonitor) update() {
	if monitor.hasObservers() {
		monitor.last = monitor.current
	} else {
		monitor.last = monitor.source.CurrentValue()
	}

	monitor.timer = nil
	monitor.DigitalNotifier.Notify(monitor.last)
}
