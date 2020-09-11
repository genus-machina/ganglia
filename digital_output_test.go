package ganglia

import (
	"testing"
)

func TestDigitalOutput(t *testing.T) {
	channel := make(chan DigitalValue, 1)
	var results []DigitalValue

	go func() {
		output := DigitalOutput(channel)
		output.Write(High)
		output.Write(Low)
		output.Write(High)
		output.Close()
	}()

	for value := range channel {
		results = append(results, value)
	}

	if count := len(results); count != 3 {
		t.Errorf("expected %d value but got %d", 3, count)
	}

	if results[0] != High {
		t.Error("expected high")
	}

	if results[1] != Low {
		t.Error("expected low")
	}

	if results[2] != High {
		t.Error("expected high")
	}
}

func TestDigitalOutputInvert(t *testing.T) {
	channel := make(chan DigitalValue, 1)
	var results []DigitalValue

	go func() {
		output := DigitalOutput(channel).Invert()
		output.Write(High)
		output.Write(Low)
		output.Write(High)
		output.Close()
	}()

	for value := range channel {
		results = append(results, value)
	}

	if count := len(results); count != 3 {
		t.Errorf("expected %d value but got %d", 3, count)
	}

	if results[0] != Low {
		t.Error("expected low")
	}

	if results[1] != High {
		t.Error("expected high")
	}

	if results[2] != Low {
		t.Error("expected low")
	}
}
