package i2c

import (
	"image"
	"image/color"

	"github.com/genus-machina/ganglia"
	"github.com/genus-machina/ganglia/widgets"
	"periph.io/x/periph/devices/ssd1306"
)

type SSD1306 struct {
	device   *ssd1306.Dev
	last     *displayContext
	rotation ganglia.Rotation
	updates  chan *displayContext
}

func newSSD1306(device *ssd1306.Dev) *SSD1306 {
	display := new(SSD1306)
	display.device = device
	display.updates = make(chan *displayContext, 10)
	go display.watchUpdates()
	return display
}

func (display *SSD1306) Halt() {
	close(display.updates)
	display.device.Halt()
}

func (display *SSD1306) Render(content ganglia.Widget) {
	displayContext := createDisplayContext(
		widgets.NewRotator(content, display.rotation),
		display.updates,
	)
	displayContext.Render()
}

func (display *SSD1306) render(displayContext *displayContext) {
	bounds := display.device.Bounds()
	content := displayContext.Content()
	rerender := displayContext.Render

	var rendered image.Image
	if content == nil {
		rendered = image.NewUniform(color.Black)
	} else {
		rendered = content.Render(bounds, rerender)
	}

	display.device.Draw(rendered.Bounds(), rendered, rendered.Bounds().Min)
}

func (display *SSD1306) Rotate(rotation ganglia.Rotation) *SSD1306 {
	display.rotation = rotation
	return display
}

func (display *SSD1306) watchUpdates() {
	for displayContext := range display.updates {
		if display.last != nil && display.last != displayContext {
			display.last.Halt()
		}

		display.last = displayContext
		display.render(displayContext)
	}
}
