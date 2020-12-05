package monitors

import (
	"time"

	"github.com/genus-machina/ganglia"
)

type HoldMonitor struct {
	current *ganglia.DigitalEvent
	DigitalNotifier
	observer *DigitalEventObserver
	source   DigitalMonitor
	timer    *time.Timer
}

func NewHoldMonitor(source DigitalMonitor) *HoldMonitor {
	monitor := new(HoldMonitor)
	monitor.observer = NewDigitalEventObserver(monitor.handleEvent)
	monitor.source = source
	return monitor
}

func (monitor *HoldMonitor) CurrentValue() *ganglia.DigitalEvent {
	return monitor.current
}

func (monitor *HoldMonitor) handleEvent(event *ganglia.DigitalEvent) {
	if event.Value == ganglia.High && monitor.timer == nil {
		monitor.timer = time.AfterFunc(
			time.Second,
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

func (monitor *HoldMonitor) Subscribe(observer *DigitalEventObserver) ganglia.Trigger {
	if len(monitor.DigitalNotifier.observers) == 0 {
		monitor.source.Subscribe(monitor.observer)
	}

	monitor.DigitalNotifier.Subscribe(observer)
	return monitor.triggerUnsubscribe(observer)
}

func (monitor *HoldMonitor) triggerUnsubscribe(observer *DigitalEventObserver) ganglia.Trigger {
	return func() {
		monitor.Unsubscribe(observer)
	}
}

func (monitor *HoldMonitor) Unsubscribe(observer *DigitalEventObserver) {
	monitor.DigitalNotifier.Unsubscribe(observer)

	if len(monitor.DigitalNotifier.observers) == 0 {
		monitor.source.Unsubscribe(monitor.observer)
	}
}

func (monitor *HoldMonitor) update(event *ganglia.DigitalEvent) {
	monitor.current = event
	monitor.DigitalNotifier.Notify(monitor.current)
}
