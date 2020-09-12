package monitors

import (
	"github.com/genus-machina/ganglia"
)

type currentValueFactory func() *ganglia.DigitalEvent

type digitalMonitorComposite struct {
	currentValue currentValueFactory
	digitalNotifier
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
	composite.digitalNotifier.handleEvent(compositeEvent)
}

func (composite *digitalMonitorComposite) Subscribe(observer *DigitalEventObserver) {
	if len(composite.digitalNotifier.observers) == 0 {
		composite.left.Subscribe(composite.observer)
		composite.right.Subscribe(composite.observer)
	}

	composite.digitalNotifier.Subscribe(observer)
}

func (composite *digitalMonitorComposite) Unsubscribe(observer *DigitalEventObserver) {
	composite.digitalNotifier.Unsubscribe(observer)

	if len(composite.digitalNotifier.observers) == 0 {
		composite.left.Unsubscribe(composite.observer)
		composite.right.Unsubscribe(composite.observer)
	}
}
