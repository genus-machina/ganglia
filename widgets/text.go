package widgets

import (
	"image"
	"image/color"
	"sync"

	"github.com/genus-machina/ganglia"
	"golang.org/x/image/font"
)

var (
	// The font drawer is not thread safe. Use this to limit
	// drawing to a single thread.
	// See https://github.com/golang/freetype/issues/65
	TextMutex = new(sync.Mutex)
)

type Text struct {
	face       font.Face
	text       string
	textBounds image.Rectangle
}

func NewText(face font.Face, text string) *Text {
	bounds, _ := font.BoundString(face, text)
	textBounds := rectangleFromFixed(bounds)

	widget := new(Text)
	widget.face = face
	widget.text = text
	widget.textBounds = textBounds
	return widget
}

func (widget *Text) computeBounds(bounds image.Rectangle) image.Rectangle {
	delta := bounds.Min.Sub(widget.textBounds.Min)
	translated := widget.textBounds.Add(delta)
	return bounds.Intersect(translated)
}

func (widget *Text) Render(bounds image.Rectangle, rerender ganglia.Trigger) image.Image {
	TextMutex.Lock()
	defer TextMutex.Unlock()

	buffer := image.NewNRGBA(widget.computeBounds(bounds))

	drawer := new(font.Drawer)
	drawer.Dot = pointToFixed(bounds.Min.Sub(widget.textBounds.Min))
	drawer.Dst = buffer
	drawer.Face = widget.face
	drawer.Src = image.NewUniform(color.White)
	drawer.DrawString(widget.text)

	return buffer
}
