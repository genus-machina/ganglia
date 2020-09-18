package widgets

import (
	"image"
	"image/color"

	"github.com/genus-machina/ganglia"
	"periph.io/x/periph/devices/ssd1306"
)

type SSD1306 struct {
	device  *ssd1306.Dev
	last    *context
	updates chan *context
}

func NewSSD1306(device *ssd1306.Dev) *SSD1306 {
	display := new(SSD1306)
	display.device = device
	display.updates = make(chan *context, 10)
	go display.watchUpdates()
	return display
}

func (display *SSD1306) Halt() {
	close(display.updates)
	display.device.Halt()
}

func (display *SSD1306) Render(content ganglia.Widget) {
	context := createContext(content, display.updates)
	context.Render()
}

func (display *SSD1306) render(context *context) {
	bounds := display.device.Bounds()
	content := context.Content()
	rerender := context.Render

	var rendered image.Image
	if content == nil {
		rendered = image.NewUniform(color.Black)
	} else {
		rendered = content.Render(bounds, rerender)
	}

	display.device.Draw(rendered.Bounds(), rendered, rendered.Bounds().Min)
}

func (display *SSD1306) watchUpdates() {
	for context := range display.updates {
		if display.last != nil && display.last != context {
			display.last.Halt()
		}

		display.last = context
		display.render(context)
	}
}
