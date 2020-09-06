package ganglia

import (
	"time"
)

type AnalogTestEvent struct {
	T time.Time
	V AnalogValue
}

func (event *AnalogTestEvent) Time() time.Time {
	return event.T
}

func (event *AnalogTestEvent) Value() AnalogValue {
	return event.V
}

type DigitalTestEvent struct {
	T time.Time
	V DigitalValue
}

func (event *DigitalTestEvent) Time() time.Time {
	return event.T
}

func (event *DigitalTestEvent) Value() DigitalValue {
	return event.V
}
