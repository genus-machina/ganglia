package monitors

import (
	"testing"
	"time"

	"github.com/genus-machina/ganglia"
)

func TestAnalogInputMonitorValue(t *testing.T) {
	source := make(chan *ganglia.AnalogEvent, 1)
	monitor := NewAnalogInputMonitor(source)
	defer close(source)

	if value := monitor.CurrentValue(); value != nil {
		t.Error("got initial value", value)
	}

	event := &ganglia.AnalogEvent{time.Now(), 42}
	source <- event
	time.Sleep(100 * time.Millisecond)

	if eventTime := monitor.CurrentValue().Time; eventTime != event.Time {
		t.Errorf("wanted %s got %s", event.Time.String(), eventTime.String())
	}

	if eventValue := monitor.CurrentValue().Value; eventValue != event.Value {
		t.Errorf("wanted %d got %d", event.Value, eventValue)
	}

	event = &ganglia.AnalogEvent{time.Now(), 43}
	source <- event
	time.Sleep(100 * time.Millisecond)

	if eventTime := monitor.CurrentValue().Time; eventTime != event.Time {
		t.Errorf("wanted %s got %s", event.Time.String(), eventTime.String())
	}

	if eventValue := monitor.CurrentValue().Value; eventValue != event.Value {
		t.Errorf("wanted %d got %d", event.Value, eventValue)
	}
}

func TestAnalogInputMonitorSubscription(t *testing.T) {
	source := make(chan *ganglia.AnalogEvent, 1)
	monitor := NewAnalogInputMonitor(source)
	defer close(source)

	var results []ganglia.AnalogValue
	observer := NewAnalogEventObserver(func(event *ganglia.AnalogEvent) {
		results = append(results, event.Value)
	})
	monitor.Subscribe(observer)

	events := []*ganglia.AnalogEvent{
		&ganglia.AnalogEvent{Value: 1},
		&ganglia.AnalogEvent{Value: 2},
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
			t.Errorf("index %d: expected %d but got %d", index, expected, result)
		}
	}

	monitor.Unsubscribe(observer)
	source <- &ganglia.AnalogEvent{Value: 3}
	time.Sleep(100 * time.Millisecond)

	if count := len(results); count != 2 {
		t.Errorf("expeted 2 results but got %d", count)
	}
}
