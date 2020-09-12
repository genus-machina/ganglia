package monitors

import (
	"github.com/genus-machina/ganglia"
)

type DigitalMonitorIntersection struct {
	*digitalMonitorComposite
}

func NewDigitalMonitorIntersection(left, right DigitalMonitor) *DigitalMonitorIntersection {
	intersection := new(DigitalMonitorIntersection)
	intersection.digitalMonitorComposite = newDigitalMonitorComposite(left, right, intersection.CurrentValue)
	return intersection
}

func (intersection *DigitalMonitorIntersection) CurrentValue() *ganglia.DigitalEvent {
	left := intersection.left.CurrentValue()
	right := intersection.right.CurrentValue()

	if left == nil {
		return right
	}

	if right == nil {
		return left
	}

	if left.Value == ganglia.Low {
		return left
	}

	return right
}
