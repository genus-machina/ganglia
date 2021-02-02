package widgets

import (
	"image"
	"image/color"
	"time"

	"github.com/genus-machina/ganglia"
	"golang.org/x/image/font"
)

const (
	delay = 100 * time.Millisecond
)

type Marquee struct {
	face       font.Face
	nextUpdate time.Time
	offset     int
	text       string
	textBounds image.Rectangle
	timer      *time.Timer
}

func NewMarquee(face font.Face, text string) *Marquee {
	TextMutex.Lock()
	defer TextMutex.Unlock()

	bounds, _ := font.BoundString(face, text)
	textBounds := rectangleFromFixed(bounds)

	widget := new(Marquee)
	widget.face = face
	widget.text = text
	widget.textBounds = textBounds
	return widget
}

func (widget *Marquee) advance(bounds image.Rectangle) {
	if time.Now().After(widget.nextUpdate) {
		widget.nextUpdate = time.Now().Add(delay)
		widget.offset = widget.offset + 5
	}

	if widget.offset > bounds.Dx()+widget.textBounds.Dx() {
		widget.offset = 0
	}
}

func (widget *Marquee) computeBounds(bounds image.Rectangle) image.Rectangle {
	delta := bounds.Min.Sub(widget.textBounds.Min)
	translated := widget.textBounds.Add(delta)
	return bounds.Intersect(translated)
}

func (widget *Marquee) Render(bounds image.Rectangle, rerender ganglia.Trigger) image.Image {
	buffer := image.NewNRGBA(widget.computeBounds(bounds))
	widget.advance(bounds)
	translation := widget.textBounds.Min.Sub(image.Pt(bounds.Dx()-widget.offset, 0))

	drawer := new(font.Drawer)
	drawer.Dot = pointToFixed(bounds.Min.Sub(translation))
	drawer.Dst = buffer
	drawer.Face = widget.face
	drawer.Src = image.NewUniform(color.White)
	drawer.DrawString(widget.text)

	if widget.timer == nil {
		widget.timer = time.AfterFunc(
			delay,
			func() {
				widget.timer = nil
				rerender()
			},
		)
	}

	return buffer
}
