package widgets

import (
	"testing"

	"github.com/genus-machina/ganglia/widgets/test"
)

func TestText(t *testing.T) {
	display := test.NewTestDisplay(t)
	text := NewText(BuildFontFace(20), "Hello")
	display.Render(text)
}

func TestTextOverflow(t *testing.T) {
	display := test.NewTestDisplay(t)
	text := NewText(BuildFontFace(20), "Greetings")
	display.Render(text)
}
