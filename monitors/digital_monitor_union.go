package monitors

import (
	"github.com/genus-machina/ganglia"
)

type DigitalMonitorUnion struct {
	*digitalMonitorComposite
}

func NewDigitalMonitorUnion(left, right DigitalMonitor) *DigitalMonitorUnion {
	intersection := new(DigitalMonitorUnion)
	intersection.digitalMonitorComposite = newDigitalMonitorComposite(left, right, intersection.CurrentValue)
	return intersection
}

func (intersection *DigitalMonitorUnion) CurrentValue() *ganglia.DigitalEvent {
	left := intersection.left.CurrentValue()
	right := intersection.right.CurrentValue()

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
