package widgets

import (
	"testing"
)

func TestTextArea(t *testing.T) {
	display := NewTestDisplay(t)
	area := NewTextArea(BuildFontFace(20), "Hello World")
	display.Render(area)
}

func TestTextAreaSentence(t *testing.T) {
	display := NewTestDisplay(t)
	area := NewTextArea(BuildFontFace(9), "Hello super duper fantastically sassy world!")
	display.Render(area)
}

func TestTextAreaOverflow(t *testing.T) {
	display := NewTestDisplay(t)
	area := NewTextArea(BuildFontFace(20), "Hello Happy World")
	display.Render(area)
}
