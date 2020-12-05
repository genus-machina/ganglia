package widgets

import (
	"testing"

	"github.com/genus-machina/ganglia/widgets/test"
)

func TestTextArea(t *testing.T) {
	display := test.NewTestDisplay(t)
	area := NewTextArea(BuildFontFace(20), "Hello World")
	display.Render(area)
}

func TestTextAreaSmall(t *testing.T) {
	display := test.NewTestDisplay(t)
	area := NewTextArea(BuildFontFace(10), "Hi")
	display.Render(area)
}

func TestTextAreaSentence(t *testing.T) {
	display := test.NewTestDisplay(t)
	area := NewTextArea(BuildFontFace(9), "Hello super duper fantastically sassy world!")
	display.Render(area)
}

func TestTextAreaOverflow(t *testing.T) {
	display := test.NewTestDisplay(t)
	area := NewTextArea(BuildFontFace(20), "Hello Happy World")
	display.Render(area)
}
