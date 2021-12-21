package ganglia

import (
	"testing"
	"time"
)

func TestDigitalInputRead(t *testing.T) {
	input := make(chan *DigitalEvent, 1)
	event := &DigitalEvent{time.Now(), High}
	input <- event
	close(input)

	result := DigitalInput(input).Read()

	if value := result.Value; value != event.Value {
		t.Error("expected high")
	}

	if result.Time != event.Time {
		t.Errorf("expected %s got %s", event.Time.String(), result.Time.String())
	}
}

func TestDigitalInputDebounce(t *testing.T) {
	var results []*DigitalEvent
	input := make(chan *DigitalEvent, 1)

	go func() {
		input <- &DigitalEvent{time.Unix(0, 0), High}
		input <- &DigitalEvent{time.Unix(1, 0), Low}
		input <- &DigitalEvent{time.Unix(2, 0), High}
		close(input)
	}()

	for event := range DigitalInput(input).Debounce(100 * time.Millisecond) {
		results = append(results, event)
	}

	if count := len(results); count != 1 {
		t.Errorf("expected %d got %d", 1, count)
	}

	if results[0].Time != time.Unix(0, 0) {
		t.Errorf("expected %s got %s", time.Unix(2, 0), results[0].Time)
	}

	if results[0].Value != High {
		t.Errorf("expected high")
	}
}

func TestDigitalInputDebounceEmpty(t *testing.T) {
	var results []*DigitalEvent
	input := make(chan *DigitalEvent, 1)
	close(input)

	for event := range DigitalInput(input).Debounce(100 * time.Millisecond) {
		results = append(results, event)
	}

	if count := len(results); count != 0 {
		t.Errorf("expected %d got %d", 0, count)
	}
}

func TestDigitalInputDedup(t *testing.T) {
	var results []*DigitalEvent
	input := make(chan *DigitalEvent, 1)

	go func() {
		input <- &DigitalEvent{time.Unix(0, 0), High}
		input <- &DigitalEvent{time.Unix(1, 0), High}
		input <- &DigitalEvent{time.Unix(2, 0), Low}
		close(input)
	}()

	for event := range DigitalInput(input).Dedup() {
		results = append(results, event)
	}

	if count := len(results); count != 2 {
		t.Errorf("expected %d got %d", 2, count)
	}

	if results[0].Time != time.Unix(0, 0) {
		t.Errorf("expected %s got %s", time.Unix(2, 0), results[0].Time)
	}

	if results[0].Value != High {
		t.Errorf("expected high")
	}

	if results[1].Time != time.Unix(2, 0) {
		t.Errorf("expected %s got %s", time.Unix(2, 0), results[1].Time)
	}

	if results[1].Value != Low {
		t.Errorf("expected low")
	}
}

func TestDigitalInputInvert(t *testing.T) {
	input := make(chan *DigitalEvent, 1)
	event := &DigitalEvent{time.Now(), High}
	input <- event
	close(input)

	result := DigitalInput(input).Invert().Read()

	if value := result.Value; value != Low {
		t.Error("expected low")
	}

	if result.Time != event.Time {
		t.Errorf("expected %s got %s", event.Time.String(), result.Time.String())
	}
}
