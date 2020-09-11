package monitors

import (
	"testing"
	"time"

	"github.com/genus-machina/ganglia"
)

func TestDigitalMonitorMap(t *testing.T) {
	channel1 := make(chan *ganglia.DigitalEvent, 1)
	channel2 := make(chan *ganglia.DigitalEvent, 1)
	defer close(channel1)
	defer close(channel2)

	monitors := make(DigitalMonitorMap)
	monitors["one"] = NewDigitalInputMonitor(channel1)
	monitors["two"] = NewDigitalInputMonitor(channel2)

	count := 0
	observer := NewDigitalEventObserver(func(event *ganglia.DigitalEvent) {
		if event.Value == ganglia.High {
			count++
		}
	})

	monitors.Subscribe(observer)

	channel1 <- &ganglia.DigitalEvent{Value: ganglia.High}
	channel1 <- &ganglia.DigitalEvent{Value: ganglia.Low}
	channel2 <- &ganglia.DigitalEvent{Value: ganglia.High}
	time.Sleep(100 * time.Millisecond)

	if count != 2 {
		t.Errorf("expected %d but got %d", 2, count)
	}

	monitors.Unsubscribe(observer)

	channel1 <- &ganglia.DigitalEvent{Value: ganglia.High}
	channel1 <- &ganglia.DigitalEvent{Value: ganglia.Low}
	channel2 <- &ganglia.DigitalEvent{Value: ganglia.High}
	time.Sleep(100 * time.Millisecond)

	if count != 2 {
		t.Errorf("expected %d but got %d", 2, count)
	}
}
