package widgets

import (
	"image"
	"testing"
)

func TestGIFRenderAnimated(t *testing.T) {
	display := NewTestDisplay(t)
	gif, _ := NewGIF("ballerine.gif")

	display.
		MaxFrames(10).
		Render(gif)
}

func TestGIFScaleAnimated(t *testing.T) {
	display := NewTestDisplay(t)
	gif, _ := NewGIF("ballerine.gif")

	display.
		Bounds(image.Rect(0, 0, 12, 12)).
		MaxFrames(3).
		Render(gif)
}

func TestGIFRenderStatic(t *testing.T) {
	display := NewTestDisplay(t)
	gif, _ := NewGIF("bunny.gif")

	display.
		MaxFrames(2).
		Render(gif)
}
