package monitors

import (
	"time"

	"github.com/genus-machina/ganglia"
)

type SustainedDigitalMonitor struct {
	digitalNotifier
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
		monitor.digitalNotifier.handleEvent(event)

		if monitor.timer != nil {
			monitor.timer.Stop()
		}

		monitor.timer = time.AfterFunc(monitor.duration, monitor.update)
	}
}

func (monitor *SustainedDigitalMonitor) Subscribe(observer *DigitalEventObserver) ganglia.Trigger {
	if len(monitor.digitalNotifier.observers) == 0 {
		monitor.source.Subscribe(monitor.observer)
	}

	return monitor.digitalNotifier.Subscribe(observer)
}

func (monitor *SustainedDigitalMonitor) Unsubscribe(observer *DigitalEventObserver) {
	monitor.digitalNotifier.Unsubscribe(observer)

	if len(monitor.digitalNotifier.observers) == 0 {
		monitor.source.Unsubscribe(monitor.observer)
	}
}

func (monitor *SustainedDigitalMonitor) update() {
	monitor.last = monitor.current
	monitor.digitalNotifier.handleEvent(monitor.last)
}
