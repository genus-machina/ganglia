package display

import (
	"testing"
)

func TestText(t *testing.T) {
	display := NewTestDisplay(t)
	text := NewText(BuildFontFace(20), "Hello")
	display.Render(text)
}

func TestTextOverflow(t *testing.T) {
	display := NewTestDisplay(t)
	text := NewText(BuildFontFace(20), "Greetings")
	display.Render(text)
}
