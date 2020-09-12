package ganglia

import (
	"testing"
	"time"
)

func TestDigitalOutputGroup(t *testing.T) {
	var level1, level2, level3 DigitalValue
	channel1 := make(chan DigitalValue, 1)
	channel2 := make(chan DigitalValue, 1)
	channel3 := make(chan DigitalValue, 1)
	defer close(channel1)
	defer close(channel2)
	defer close(channel3)

	go func() {
		for value := range channel1 {
			level1 = value
		}
	}()

	go func() {
		for value := range channel2 {
			level2 = value
		}
	}()

	go func() {
		for value := range channel3 {
			level3 = value
		}
	}()

	group := DigitalOutputGroup(
		[]DigitalOutput{channel1, channel2, channel3},
	)
	time.Sleep(100 * time.Millisecond)

	if level1 == High {
		t.Error("expected low")
	}

	if level2 == High {
		t.Error("expected low")
	}

	if level3 == High {
		t.Error("expected low")
	}

	group.Write(5)
	time.Sleep(100 * time.Millisecond)

	if level1 == Low {
		t.Error("expected high")
	}

	if level2 == High {
		t.Error("expected low")
	}

	if level3 == Low {
		t.Error("expected high")
	}
}
