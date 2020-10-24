package widgets

import (
	"image"
	"strings"

	"github.com/genus-machina/ganglia"
	"golang.org/x/image/font"
)

type TextArea struct {
	face font.Face
	text string
}

func NewTextArea(face font.Face, text string) *TextArea {
	widget := new(TextArea)
	widget.face = face
	widget.text = text
	return widget
}

func (widget *TextArea) buildRows(width int) []string {
	var rows []string
	separator := " "
	tokens := strings.Split(widget.text, separator)
	current := tokens[0]
	tokens = tokens[1:]

	for _, token := range tokens {
		candidate := current + separator + token
		rowWidth := font.MeasureString(widget.face, candidate).Ceil()

		if rowWidth > width {
			rows = append(rows, current)
			current = token
		} else {
			current = candidate
		}
	}

	rows = append(rows, current)
	return rows
}

func (widget *TextArea) Render(bounds image.Rectangle, rerender ganglia.Trigger) image.Image {
	box := NewBox().Justify(Center)
	rows := widget.buildRows(bounds.Dx())

	for _, row := range rows {
		text := NewText(widget.face, row)
		box.Append(text)
	}

	return box.Render(bounds, rerender)
}
