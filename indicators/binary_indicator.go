package indicators

import (
	"github.com/genus-machina/ganglia"
	"github.com/genus-machina/ganglia/monitors"
)

type BinaryIndicator struct {
	monitors monitors.DigitalMonitorGroup
	observer *monitors.DigitalEventObserver
	outputs  ganglia.DigitalOutputGroup
}

func NewBinaryIndicator(inputs monitors.DigitalMonitorGroup, outputs ganglia.DigitalOutputGroup) *BinaryIndicator {
	indicator := new(BinaryIndicator)
	indicator.monitors = inputs
	indicator.observer = monitors.NewDigitalEventObserver(indicator.handleEvent)
	indicator.outputs = outputs
	return indicator
}

func (indicator *BinaryIndicator) handleEvent(event *ganglia.DigitalEvent) {
	indicator.update()
}

func (indicator *BinaryIndicator) Start() {
	indicator.monitors.Subscribe(indicator.observer)
	indicator.update()
}

func (indicator *BinaryIndicator) Stop() {
	indicator.monitors.Unsubscribe(indicator.observer)
}

func (indicator *BinaryIndicator) update() {
	value := indicator.monitors.Value()
	indicator.outputs.Write(value)
}
