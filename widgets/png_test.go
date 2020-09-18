package widgets

import (
	"image"
	"testing"
)

func TestPNG(t *testing.T) {
	display := NewTestDisplay(t)
	png, _ := NewPNG("dialog-warning-symbolic.png")
	display.Render(png)
}

func TestPNGScale(t *testing.T) {
	display := NewTestDisplay(t)
	png, _ := NewPNG("dialog-warning-symbolic.png")

	display.
		Bounds(image.Rect(0, 0, 12, 12)).
		Render(png)
}
