package monitors

import (
	"github.com/genus-machina/ganglia"
)

type DigitalMonitorUnion struct {
	*digitalMonitorComposite
}

func NewDigitalMonitorUnion(left, right DigitalMonitor) *DigitalMonitorUnion {
	union := new(DigitalMonitorUnion)
	union.digitalMonitorComposite = newDigitalMonitorComposite(left, right, union.CurrentValue)
	return union
}

func (union *DigitalMonitorUnion) CurrentValue() *ganglia.DigitalEvent {
	left := union.left.CurrentValue()
	right := union.right.CurrentValue()

	if left == nil {
		return right
	}

	if right == nil {
		return left
	}

	if left.Value == ganglia.Low {
		return right
	}

	return left
}
