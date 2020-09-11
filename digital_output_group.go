package ganglia

import (
	"time"
)

type DigitalOutputGroup []DigitalOutput

func (group DigitalOutputGroup) Blink(value uint, on, off time.Duration) Stopper {
	done := make(chan bool)
	go group.blink(value, on, off, done)
	return func() { close(done) }
}

func (group DigitalOutputGroup) blink(value uint, on, off time.Duration, done <-chan bool) {
	active := true
	timer := time.After(on)

	for running := true; running; {
		var out uint

		if active {
			out = value
		}

		group.Write(out)

		select {
		case <-done:
			running = false
		case <-timer:
			active = !active

			if active {
				timer = time.After(on)
			} else {
				timer = time.After(off)
			}
		}
	}
}

func (group DigitalOutputGroup) Write(value uint) {
	count := len(group)

	for index, output := range group {
		if value&(1<<(count-index-1)) > 0 {
			output <- High
		} else {
			output <- Low
		}
	}
}

type Stopper func()
