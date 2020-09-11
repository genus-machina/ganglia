package monitors

import (
	"time"

	"github.com/genus-machina/ganglia"
)

type AnalogTestEvent struct {
	T time.Time
	V ganglia.AnalogValue
}

func (event *AnalogTestEvent) Time() time.Time {
	return event.T
}

func (event *AnalogTestEvent) Value() ganglia.AnalogValue {
	return event.V
}

type DigitalTestEvent struct {
	T time.Time
	V ganglia.DigitalValue
}

func (event *DigitalTestEvent) Time() time.Time {
	return event.T
}

func (event *DigitalTestEvent) Value() ganglia.DigitalValue {
	return event.V
}
