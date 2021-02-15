package widgets

import (
	"testing"

	"github.com/genus-machina/ganglia/widgets/test"
)

func TestMarqueeLong(t *testing.T) {
	display := test.NewTestDisplay(t)
	marquee := NewMarquee(BuildFontFace(10), "A long time ago in a galaxy far far away...")

	display.
		MaxFrames(100).
		Render(marquee)
}

func TestMarqueeShort(t *testing.T) {
	display := test.NewTestDisplay(t)
	marquee := NewMarquee(BuildFontFace(10), "Hi")

	display.
		MaxFrames(100).
		Render(marquee)
}
