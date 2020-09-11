package monitors

import (
	"testing"
	"time"

	"github.com/genus-machina/ganglia"
)

func TestDigitalMonitorGroup(t *testing.T) {
	channel1 := make(chan *ganglia.DigitalEvent, 1)
	channel2 := make(chan *ganglia.DigitalEvent, 1)
	channel3 := make(chan *ganglia.DigitalEvent, 1)
	defer close(channel1)
	defer close(channel2)
	defer close(channel3)

	group := DigitalMonitorGroup(
		[]DigitalMonitor{
			NewDigitalInputMonitor(channel1),
			NewDigitalInputMonitor(channel2),
			NewDigitalInputMonitor(channel3),
		},
	)
	time.Sleep(100 * time.Millisecond)

	if value := group.Value(); value != 0 {
		t.Errorf("expected %d but got %d", 0, value)
	}

	channel1 <- &ganglia.DigitalEvent{Value: ganglia.High}
	channel3 <- &ganglia.DigitalEvent{Value: ganglia.High}
	time.Sleep(100 * time.Millisecond)

	if value := group.Value(); value != 5 {
		t.Errorf("expected %d but got %d", 5, value)
	}

	channel2 <- &ganglia.DigitalEvent{Value: ganglia.Low}
	channel3 <- &ganglia.DigitalEvent{Value: ganglia.Low}
	time.Sleep(100 * time.Millisecond)

	if value := group.Value(); value != 4 {
		t.Errorf("expected %d but got %d", 4, value)
	}
}
