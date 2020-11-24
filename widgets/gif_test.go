package widgets

import (
	"image"
	"testing"

	"github.com/genus-machina/ganglia/widgets/test"
)

func TestGIFRenderAnimated(t *testing.T) {
	display := test.NewTestDisplay(t)
	gif, _ := OpenGIF("ballerine.gif")

	display.
		MaxFrames(10).
		Render(gif)
}

func TestGIFScaleAnimated(t *testing.T) {
	display := test.NewTestDisplay(t)
	gif, _ := OpenGIF("ballerine.gif")

	display.
		Bounds(image.Rect(0, 0, 12, 12)).
		MaxFrames(3).
		Render(gif)
}

func TestGIFRenderStatic(t *testing.T) {
	display := test.NewTestDisplay(t)
	gif, _ := OpenGIF("bunny.gif")

	display.
		MaxFrames(2).
		Render(gif)
}
