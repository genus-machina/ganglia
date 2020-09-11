package monitors

import (
	"github.com/genus-machina/ganglia"
)

type DigitalMonitorGroup []DigitalMonitor

func (group DigitalMonitorGroup) Subscribe(observer *DigitalEventObserver) {
	for _, monitor := range group {
		monitor.Subscribe(observer)
	}
}

func (group DigitalMonitorGroup) Unsubscribe(observer *DigitalEventObserver) {
	for _, monitor := range group {
		monitor.Unsubscribe(observer)
	}
}

func (group DigitalMonitorGroup) Value() uint {
	var value uint
	count := len(group)

	for index, monitor := range group {
		if current := monitor.CurrentValue(); current != nil && current.Value == ganglia.High {
			value += 1 << (count - index - 1)
		}
	}

	return value
}
