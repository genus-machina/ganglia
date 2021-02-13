package mqtt

import (
	"github.com/genus-machina/ganglia"
	"github.com/genus-machina/ganglia/monitors"
)

type Monitor struct {
	broker *Broker
}

func NewMonitor(broker *Broker) *Monitor {
	monitor := new(Monitor)
	monitor.broker = broker
	return monitor
}

func (monitor *Monitor) CurrentValue() *ganglia.DigitalEvent {
	return monitor.broker.event
}

func (monitor *Monitor) Once(observer *monitors.DigitalEventObserver) ganglia.Trigger {
	return monitor.broker.notifier.Once(observer)
}

func (monitor *Monitor) Subscribe(observer *monitors.DigitalEventObserver) ganglia.Trigger {
	return monitor.broker.notifier.Subscribe(observer)
}

func (monitor *Monitor) Unsubscribe(observer *monitors.DigitalEventObserver) {
	monitor.broker.notifier.Unsubscribe(observer)
}
