package monitors

import (
	"testing"
	"time"

	"github.com/genus-machina/ganglia"
)

func TestDigitalInputMonitorValue(t *testing.T) {
	source := make(chan *ganglia.DigitalEvent, 1)
	monitor := NewDigitalInputMonitor(source)
	defer close(source)

	if value := monitor.CurrentValue(); value != nil {
		t.Error("got initial value", value)
	}

	event := &ganglia.DigitalEvent{time.Now(), ganglia.High}
	source <- event
	time.Sleep(100 * time.Millisecond)

	if eventTime := monitor.CurrentValue().Time; eventTime != event.Time {
		t.Errorf("wanted %s got %s", event.Time.String(), eventTime.String())
	}

	if eventValue := monitor.CurrentValue().Value; eventValue != event.Value {
		t.Errorf("wanted %t got %t", event.Value, eventValue)
	}

	event = &ganglia.DigitalEvent{time.Now(), ganglia.Low}
	source <- event
	time.Sleep(100 * time.Millisecond)

	if eventTime := monitor.CurrentValue().Time; eventTime != event.Time {
		t.Errorf("wanted %s got %s", event.Time.String(), eventTime.String())
	}

	if eventValue := monitor.CurrentValue().Value; eventValue != event.Value {
		t.Errorf("wanted %t got %t", event.Value, eventValue)
	}
}

func TestDigitalInputMonitorSubscription(t *testing.T) {
	source := make(chan *ganglia.DigitalEvent, 1)
	monitor := NewDigitalInputMonitor(source)
	defer close(source)

	var results []ganglia.DigitalValue
	observer := NewDigitalEventObserver(func(event *ganglia.DigitalEvent) {
		results = append(results, event.Value)
	})
	monitor.Subscribe(observer)

	events := []*ganglia.DigitalEvent{
		&ganglia.DigitalEvent{Value: ganglia.High},
		&ganglia.DigitalEvent{Value: ganglia.Low},
	}

	for _, event := range events {
		source <- event
	}
	time.Sleep(100 * time.Millisecond)

	if count := len(results); count != 2 {
		t.Errorf("expected 2 results but got %d", count)
	}

	for index, result := range results {
		if expected := events[index].Value; result != expected {
			t.Errorf("index %d: expected %t but got %t", index, expected, result)
		}
	}

	monitor.Unsubscribe(observer)
	source <- &ganglia.DigitalEvent{Value: ganglia.High}
	time.Sleep(100 * time.Millisecond)

	if count := len(results); count != 2 {
		t.Errorf("expeted 2 results but got %d", count)
	}
}
