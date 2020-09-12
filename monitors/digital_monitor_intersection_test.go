package monitors

import (
	"testing"
	"time"

	"github.com/genus-machina/ganglia"
)

func TestDigitarMonitorIntersection(t *testing.T) {
	input1 := make(chan *ganglia.DigitalEvent, 1)
	input2 := make(chan *ganglia.DigitalEvent, 1)
	defer close(input1)
	defer close(input2)

	monitor1 := NewDigitalInputMonitor(input1)
	monitor2 := NewDigitalInputMonitor(input2)
	intersection := NewDigitalMonitorIntersection(monitor1, monitor2)
	time.Sleep(100 * time.Millisecond)

	if event := intersection.CurrentValue(); event != nil {
		t.Error("expected nil")
	}

	input1 <- &ganglia.DigitalEvent{Value: ganglia.Low}
	time.Sleep(100 * time.Millisecond)

	if event := intersection.CurrentValue(); event.Value != ganglia.Low {
		t.Error("expected low")
	}

	input2 <- &ganglia.DigitalEvent{Value: ganglia.High}
	time.Sleep(100 * time.Millisecond)

	if event := intersection.CurrentValue(); event.Value != ganglia.Low {
		t.Error("expected low")
	}

	input1 <- &ganglia.DigitalEvent{Value: ganglia.High}
	time.Sleep(100 * time.Millisecond)

	if event := intersection.CurrentValue(); event.Value != ganglia.High {
		t.Error("expected high")
	}
}

func TestDigitarMonitorIntersectionSubscribe(t *testing.T) {
	input1 := make(chan *ganglia.DigitalEvent, 1)
	input2 := make(chan *ganglia.DigitalEvent, 1)
	defer close(input1)
	defer close(input2)

	monitor1 := NewDigitalInputMonitor(input1)
	monitor2 := NewDigitalInputMonitor(input2)
	intersection := NewDigitalMonitorIntersection(monitor1, monitor2)

	var currentValue ganglia.DigitalValue
	observer := NewDigitalEventObserver(func(event *ganglia.DigitalEvent) {
		currentValue = event.Value
	})
	intersection.Subscribe(observer)
	time.Sleep(100 * time.Millisecond)

	if currentValue != ganglia.Low {
		t.Error("expected low")
	}

	input1 <- &ganglia.DigitalEvent{Value: ganglia.Low}
	time.Sleep(100 * time.Millisecond)

	if currentValue != ganglia.Low {
		t.Error("expected low")
	}

	input2 <- &ganglia.DigitalEvent{Value: ganglia.High}
	time.Sleep(100 * time.Millisecond)

	if currentValue != ganglia.Low {
		t.Error("expected low")
	}

	input1 <- &ganglia.DigitalEvent{Value: ganglia.High}
	time.Sleep(100 * time.Millisecond)

	if currentValue != ganglia.High {
		t.Error("expected high")
	}

	intersection.Unsubscribe(observer)

	input1 <- &ganglia.DigitalEvent{Value: ganglia.Low}
	input2 <- &ganglia.DigitalEvent{Value: ganglia.Low}
	time.Sleep(100 * time.Millisecond)

	if currentValue != ganglia.High {
		t.Error("expected high")
	}
}
