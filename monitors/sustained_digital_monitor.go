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
	return monitor.last
}

func (monitor *SustainedDigitalMonitor) handleEvent(event *ganglia.DigitalEvent) {
	monitor.current = event

	if monitor.last == nil || event.Value == ganglia.High {
		monitor.last = event
		monitor.DigitalNotifier.Notify(event)

		if monitor.timer != nil {
			monitor.timer.Stop()
		}

		monitor.timer = time.AfterFunc(monitor.duration, monitor.update)
	}
}

func (monitor *SustainedDigitalMonitor) Subscribe(observer *DigitalEventObserver) ganglia.Trigger {
	if len(monitor.DigitalNotifier.observers) == 0 {
		monitor.source.Subscribe(monitor.observer)
	}

	return monitor.DigitalNotifier.Subscribe(observer)
}

func (monitor *SustainedDigitalMonitor) Unsubscribe(observer *DigitalEventObserver) {
	monitor.DigitalNotifier.Unsubscribe(observer)

	if len(monitor.DigitalNotifier.observers) == 0 {
		monitor.source.Unsubscribe(monitor.observer)
	}
}

func (monitor *SustainedDigitalMonitor) update() {
	monitor.last = monitor.current
	monitor.DigitalNotifier.Notify(monitor.last)
}
