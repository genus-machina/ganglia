package widgets

import (
	"image"
	"testing"

	"github.com/genus-machina/ganglia/widgets/test"
)

func TestPNG(t *testing.T) {
	display := test.NewTestDisplay(t)
	png, _ := NewPNG("dialog-warning-symbolic.png")
	display.Render(png)
}

func TestPNGScale(t *testing.T) {
	display := test.NewTestDisplay(t)
	png, _ := NewPNG("dialog-warning-symbolic.png")

	display.
		Bounds(image.Rect(0, 0, 12, 12)).
		Render(png)
}
