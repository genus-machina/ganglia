package monitors

import (
	"time"

	"github.com/genus-machina/ganglia"
)

type ThresholdDigitalMonitor struct {
	current *ganglia.DigitalEvent
	digitalNotifier
	duration time.Duration
	observer *DigitalEventObserver
	source   DigitalMonitor
	timer    *time.Timer
}

func NewThresholdDigitalMonitor(source DigitalMonitor, duration time.Duration) *ThresholdDigitalMonitor {
	monitor := new(ThresholdDigitalMonitor)
	monitor.duration = duration
	monitor.observer = NewDigitalEventObserver(monitor.handleEvent)
	monitor.source = source
	return monitor
}

func (monitor *ThresholdDigitalMonitor) CurrentValue() *ganglia.DigitalEvent {
	return monitor.current
}

func (monitor *ThresholdDigitalMonitor) handleEvent(event *ganglia.DigitalEvent) {
	if event.Value == ganglia.High && monitor.timer == nil {
		monitor.timer = time.AfterFunc(
			monitor.duration,
			func() { monitor.update(event) },
		)
	}

	if event.Value == ganglia.Low {
		if monitor.timer != nil {
			monitor.timer.Stop()
			monitor.timer = nil
		}

		monitor.update(event)
	}
}

func (monitor *ThresholdDigitalMonitor) Subscribe(observer *DigitalEventObserver) ganglia.Trigger {
	if len(monitor.digitalNotifier.observers) == 0 {
		monitor.source.Subscribe(monitor.observer)
	}

	monitor.digitalNotifier.Subscribe(observer)
	return monitor.triggerUnsubscribe(observer)
}

func (monitor *ThresholdDigitalMonitor) triggerUnsubscribe(observer *DigitalEventObserver) ganglia.Trigger {
	return func() {
		monitor.Unsubscribe(observer)
	}
}

func (monitor *ThresholdDigitalMonitor) Unsubscribe(observer *DigitalEventObserver) {
	monitor.digitalNotifier.Unsubscribe(observer)

	if len(monitor.digitalNotifier.observers) == 0 {
		monitor.source.Unsubscribe(monitor.observer)
	}
}

func (monitor *ThresholdDigitalMonitor) update(event *ganglia.DigitalEvent) {
	monitor.current = event
	monitor.digitalNotifier.handleEvent(monitor.current)
}
