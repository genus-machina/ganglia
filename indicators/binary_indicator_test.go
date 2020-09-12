package indicators

import (
	"testing"
	"time"

	"github.com/genus-machina/ganglia"
	"github.com/genus-machina/ganglia/monitors"
)

func TestBinaryIndicator(t *testing.T) {
	input1 := make(chan *ganglia.DigitalEvent, 1)
	input2 := make(chan *ganglia.DigitalEvent, 1)
	output1 := make(chan ganglia.DigitalValue, 1)
	output2 := make(chan ganglia.DigitalValue, 1)
	defer close(input1)
	defer close(input2)
	defer close(output1)
	defer close(output2)

	inputGroup := []monitors.DigitalMonitor{
		monitors.NewDigitalInputMonitor(input1),
		monitors.NewDigitalInputMonitor(input2),
	}

	outputGroup := []ganglia.DigitalOutput{
		output1,
		output2,
	}

	var value1, value2 ganglia.DigitalValue

	go func() {
		for value := range output1 {
			value1 = value
		}
	}()

	go func() {
		for value := range output2 {
			value2 = value
		}
	}()

	indicator := NewBinaryIndicator(inputGroup, outputGroup)
	indicator.Start()
	time.Sleep(100 * time.Millisecond)

	if value1 != ganglia.Low {
		t.Errorf("expected low but got %t", value1)
	}

	if value2 != ganglia.Low {
		t.Errorf("expected low but got %t", value2)
	}

	input1 <- &ganglia.DigitalEvent{Value: ganglia.High}
	time.Sleep(100 * time.Millisecond)

	if value1 != ganglia.High {
		t.Errorf("expected high but got %t", value1)
	}

	if value2 != ganglia.Low {
		t.Errorf("expected low but got %t", value2)
	}

	input2 <- &ganglia.DigitalEvent{Value: ganglia.High}
	time.Sleep(100 * time.Millisecond)

	if value1 != ganglia.High {
		t.Errorf("expected high but got %t", value1)
	}

	if value2 != ganglia.High {
		t.Errorf("expected high but got %t", value2)
	}

	input1 <- &ganglia.DigitalEvent{Value: ganglia.Low}
	time.Sleep(100 * time.Millisecond)

	if value1 != ganglia.Low {
		t.Errorf("expected low but got %t", value1)
	}

	if value2 != ganglia.High {
		t.Errorf("expected high but got %t", value2)
	}

	indicator.Stop()

	input2 <- &ganglia.DigitalEvent{Value: ganglia.Low}
	time.Sleep(100 * time.Millisecond)

	if value1 != ganglia.Low {
		t.Errorf("expected low but got %t", value1)
	}

	if value2 != ganglia.High {
		t.Errorf("expected high but got %t", value2)
	}
}
