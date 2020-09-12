package monitors

import (
	"testing"
	"time"

	"github.com/genus-machina/ganglia"
)

func TestSustainedDigitalMonitor(t *testing.T) {
	input := make(chan *ganglia.DigitalEvent, 1)
	monitor := NewDigitalInputMonitor(input)
	sustained := NewSustainedDigitalMonitor(monitor, 500*time.Millisecond)
	defer close(input)

	var current ganglia.DigitalValue
	observer := NewDigitalEventObserver(func(event *ganglia.DigitalEvent) {
		current = event.Value
	})
	sustained.Subscribe(observer)
	time.Sleep(500 * time.Millisecond)

	if value := sustained.CurrentValue(); value != nil {
		t.Error("expected nil")
	}

	if current != ganglia.Low {
		t.Error("expected low")
	}

	input <- &ganglia.DigitalEvent{Value: ganglia.Low}
	time.Sleep(100 * time.Millisecond)

	if value := sustained.CurrentValue(); value.Value != ganglia.Low {
		t.Error("expected low")
	}

	if current != ganglia.Low {
		t.Error("expected low")
	}

	input <- &ganglia.DigitalEvent{Value: ganglia.High}
	time.Sleep(100 * time.Millisecond)

	if value := sustained.CurrentValue(); value.Value != ganglia.High {
		t.Error("expected high")
	}

	if current != ganglia.High {
		t.Error("expected high")
	}

	input <- &ganglia.DigitalEvent{Value: ganglia.Low}
	time.Sleep(100 * time.Millisecond)

	if value := sustained.CurrentValue(); value.Value != ganglia.High {
		t.Error("expected high")
	}

	if current != ganglia.High {
		t.Error("expected high")
	}

	time.Sleep(500 * time.Millisecond)

	if value := sustained.CurrentValue(); value.Value != ganglia.Low {
		t.Error("expected low")
	}

	if current != ganglia.Low {
		t.Error("expected low")
	}

	sustained.Unsubscribe(observer)
	input <- &ganglia.DigitalEvent{Value: ganglia.High}
	time.Sleep(100 * time.Millisecond)

	if value := sustained.CurrentValue(); value.Value != ganglia.Low {
		t.Error("expected low")
	}

	if current != ganglia.Low {
		t.Error("expected low")
	}
}
