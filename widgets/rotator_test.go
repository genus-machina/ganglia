package widgets

import (
	"testing"

	"github.com/genus-machina/ganglia"
	"github.com/genus-machina/ganglia/widgets/test"
)

func TestRotator0(t *testing.T) {
	display := test.NewTestDisplay(t)
	text := NewText(BuildFontFace(20), "Hello")
	rotated := NewRotator(text, ganglia.Rotate0)
	display.Render(rotated)
}

func TestRotator90(t *testing.T) {
	display := test.NewTestDisplay(t)
	text := NewText(BuildFontFace(20), "Hello")
	rotated := NewRotator(text, ganglia.Rotate90)
	display.Render(rotated)
}

func TestRotator180(t *testing.T) {
	display := test.NewTestDisplay(t)
	text := NewText(BuildFontFace(20), "Hello")
	rotated := NewRotator(text, ganglia.Rotate180)
	display.Render(rotated)
}

func TestRotator270(t *testing.T) {
	display := test.NewTestDisplay(t)
	text := NewText(BuildFontFace(20), "Hello")
	rotated := NewRotator(text, ganglia.Rotate270)
	display.Render(rotated)
}
