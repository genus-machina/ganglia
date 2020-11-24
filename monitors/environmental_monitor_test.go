package monitors

import (
	"testing"
	"time"

	"github.com/genus-machina/ganglia"
)

func TestEnvironmentalInputMonitorValue(t *testing.T) {
	source := make(chan *ganglia.EnvironmentalEvent, 1)
	monitor := NewEnvironmentalInputMonitor(source)
	defer close(source)

	if value := monitor.CurrentValue(); value != nil {
		t.Error("got initial value", value)
	}

	event := &ganglia.EnvironmentalEvent{Time: time.Now()}
	source <- event
	time.Sleep(100 * time.Millisecond)

	if eventTime := monitor.CurrentValue().Time; eventTime != event.Time {
		t.Errorf("wanted %s got %s", event.Time.String(), eventTime.String())
	}

	event = &ganglia.EnvironmentalEvent{Time: time.Now()}
	source <- event
	time.Sleep(100 * time.Millisecond)

	if eventTime := monitor.CurrentValue().Time; eventTime != event.Time {
		t.Errorf("wanted %s got %s", event.Time.String(), eventTime.String())
	}
}

func TestEnvironmentalInputMonitorSubscription(t *testing.T) {
	source := make(chan *ganglia.EnvironmentalEvent, 1)
	monitor := NewEnvironmentalInputMonitor(source)
	defer close(source)

	var results []time.Time
	observer := NewEnvironmentalEventObserver(func(event *ganglia.EnvironmentalEvent) {
		results = append(results, event.Time)
	})
	monitor.Subscribe(observer)

	events := []*ganglia.EnvironmentalEvent{
		&ganglia.EnvironmentalEvent{Time: time.Now()},
		&ganglia.EnvironmentalEvent{Time: time.Now()},
	}

	for _, event := range events {
		source <- event
	}
	time.Sleep(100 * time.Millisecond)

	if count := len(results); count != 2 {
		t.Errorf("expected 2 results but got %d", count)
	}

	for index, result := range results {
		if expected := events[index].Time; result != expected {
			t.Errorf("index %d: expected %s but got %s", index, expected.String(), result.String())
		}
	}

	monitor.Unsubscribe(observer)
	source <- &ganglia.EnvironmentalEvent{Time: time.Now()}
	time.Sleep(100 * time.Millisecond)

	if count := len(results); count != 2 {
		t.Errorf("expected 2 results but got %d", count)
	}
}
