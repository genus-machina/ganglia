package indicators

import (
	"time"

	"github.com/genus-machina/ganglia"
	"github.com/genus-machina/ganglia/monitors"
)

type Blinker struct {
	active   bool
	done     chan bool
	monitors monitors.DigitalMonitorGroup
	observer *monitors.DigitalEventObserver
	on, off  time.Duration
	outputs  []ganglia.DigitalOutputGroup
	value    uint
}

func NewBlinker(inputs monitors.DigitalMonitorGroup, on, off time.Duration, outputs ...ganglia.DigitalOutputGroup) *Blinker {
	indicator := new(Blinker)
	indicator.monitors = inputs
	indicator.observer = monitors.NewDigitalEventObserver(indicator.handleEvent)
	indicator.on = on
	indicator.off = off
	indicator.outputs = outputs
	return indicator
}

func (indicator *Blinker) blink() {
	for running := true; running; {
		var timer <-chan time.Time

		if indicator.active {
			timer = time.After(indicator.on)
		} else {
			timer = time.After(indicator.off)
		}

		select {
		case <-indicator.done:
			running = false
		case <-timer:
			indicator.active = !indicator.active
		}
	}
}

func (indicator *Blinker) handleEvent(event *ganglia.DigitalEvent) {
	indicator.update()
}

func (indicator *Blinker) Start() {
	indicator.active = true
	indicator.done = make(chan bool)
	indicator.update()
	indicator.write()
	go indicator.blink()
}

func (indicator *Blinker) Stop() {
	close(indicator.done)
}

func (indicator *Blinker) update() {
	indicator.value = indicator.monitors.Value()
}

func (indicator *Blinker) write() {
	var value uint

	if indicator.active {
		value = indicator.value
	}

	for _, output := range indicator.outputs {
		output.Write(value)
	}
}
