package ganglia

import (
	"testing"
	"time"
)

func TestAnalogMonitorValue(t *testing.T) {
	source := make(chan AnalogEvent, 1)
	monitor := NewAnalogMonitor(source)
	defer close(source)

	if value := monitor.CurrentValue(); value != nil {
		t.Error("got initial value", value)
	}

	event := &AnalogTestEvent{time.Now(), 42}
	source <- event
	time.Sleep(100 * time.Millisecond)

	if eventTime := monitor.CurrentValue().Time(); eventTime != event.T {
		t.Errorf("wanted %s got %s", event.T.String(), eventTime.String())
	}

	if eventValue := monitor.CurrentValue().Value(); eventValue != event.V {
		t.Errorf("wanted %d got %d", event.V, eventValue)
	}

	event = &AnalogTestEvent{time.Now(), 43}
	source <- event
	time.Sleep(100 * time.Millisecond)

	if eventTime := monitor.CurrentValue().Time(); eventTime != event.T {
		t.Errorf("wanted %s got %s", event.T.String(), eventTime.String())
	}

	if eventValue := monitor.CurrentValue().Value(); eventValue != event.V {
		t.Errorf("wanted %d got %d", event.V, eventValue)
	}
}

func TestDigitalMonitorValue(t *testing.T) {
	source := make(chan DigitalEvent, 1)
	monitor := NewDigitalMonitor(source)
	defer close(source)

	if value := monitor.CurrentValue(); value != nil {
		t.Error("got initial value", value)
	}

	event := &DigitalTestEvent{time.Now(), High}
	source <- event
	time.Sleep(100 * time.Millisecond)

	if eventTime := monitor.CurrentValue().Time(); eventTime != event.T {
		t.Errorf("wanted %s got %s", event.T.String(), eventTime.String())
	}

	if eventValue := monitor.CurrentValue().Value(); eventValue != event.V {
		t.Errorf("wanted %t got %t", event.V, eventValue)
	}

	event = &DigitalTestEvent{time.Now(), Low}
	source <- event
	time.Sleep(100 * time.Millisecond)

	if eventTime := monitor.CurrentValue().Time(); eventTime != event.T {
		t.Errorf("wanted %s got %s", event.T.String(), eventTime.String())
	}

	if eventValue := monitor.CurrentValue().Value(); eventValue != event.V {
		t.Errorf("wanted %t got %t", event.V, eventValue)
	}
}
