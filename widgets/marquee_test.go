package widgets

import (
	"testing"

	"github.com/genus-machina/ganglia/widgets/test"
)

func TestMarquee(t *testing.T) {
	display := test.NewTestDisplay(t)
	marquee := NewMarquee(BuildFontFace(10), "A long time ago in a galaxy far far away...")

	display.
		MaxFrames(100).
		Render(marquee)
}
