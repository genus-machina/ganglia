package ganglia

import (
	"testing"
)

func TestAnalogInput(t *testing.T) {
	channel := make(chan *AnalogEvent, 1)
	channel <- &AnalogEvent{Value: 42}
	close(channel)

	result := AnalogInput(channel).Read()

	if result.Value != 42 {
		t.Errorf("expected value %d but got %d", 42, result.Value)
	}
}
