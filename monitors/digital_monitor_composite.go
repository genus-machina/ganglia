package monitors

import (
	"github.com/genus-machina/ganglia"
)

type currentValueFactory func() *ganglia.DigitalEvent

type digitalMonitorComposite struct {
	currentValue currentValueFactory
	DigitalNotifier
	left, right DigitalMonitor
	observer    *DigitalEventObserver
}

func newDigitalMonitorComposite(left, right DigitalMonitor, currentValue currentValueFactory) *digitalMonitorComposite {
	composite := new(digitalMonitorComposite)
	composite.currentValue = currentValue
	composite.left = left
	composite.right = right
	composite.observer = NewDigitalEventObserver(composite.handleEvent)
	return composite
}

func (composite *digitalMonitorComposite) handleEvent(event *ganglia.DigitalEvent) {
	compositeEvent := composite.currentValue()
	composite.DigitalNotifier.Notify(compositeEvent)
}

func (composite *digitalMonitorComposite) Subscribe(observer *DigitalEventObserver) ganglia.Trigger {
	if len(composite.DigitalNotifier.observers) == 0 {
		composite.left.Subscribe(composite.observer)
		composite.right.Subscribe(composite.observer)
	}

	return composite.DigitalNotifier.Subscribe(observer)
}

func (composite *digitalMonitorComposite) Unsubscribe(observer *DigitalEventObserver) {
	composite.DigitalNotifier.Unsubscribe(observer)

	if len(composite.DigitalNotifier.observers) == 0 {
		composite.left.Unsubscribe(composite.observer)
		composite.right.Unsubscribe(composite.observer)
	}
}
