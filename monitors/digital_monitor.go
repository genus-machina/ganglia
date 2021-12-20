package monitors

import (
	"github.com/genus-machina/ganglia"
)

type DigitalMonitor interface {
	CurrentValue() *ganglia.DigitalEvent
	Once(*DigitalEventObserver) ganglia.Trigger
	Subscribe(*DigitalEventObserver) ganglia.Trigger
	Unsubscribe(*DigitalEventObserver)
}

type DigitalInputMonitor struct {
	DigitalNotifier
	currentValue *ganglia.DigitalEvent
	source       ganglia.DigitalInput
}

func NewDigitalInputMonitor(source ganglia.DigitalInput) *DigitalInputMonitor {
	monitor := new(DigitalInputMonitor)
	monitor.source = source
	go monitor.watchSource()
	return monitor
}

func (monitor *DigitalInputMonitor) CurrentValue() *ganglia.DigitalEvent {
	return monitor.currentValue
}

func (monitor *DigitalInputMonitor) watchSource() {
	for event := range monitor.source {
		monitor.currentValue = event
		monitor.Notify(event)
	}
}
