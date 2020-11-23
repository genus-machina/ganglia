package widgets

import (
	"image"

	"github.com/genus-machina/ganglia"
)

type Image struct {
	content image.Image
}

func NewImage(image image.Image) *Image {
	widget := new(Image)
	widget.content = image
	return widget
}

func (widget *Image) Render(bounds image.Rectangle, rerender ganglia.Trigger) image.Image {
	scaledBounds := computeImageBounds(widget.content, bounds)
	buffer := image.NewNRGBA(scaledBounds)
	scaleImage(buffer, scaledBounds, widget.content, widget.content.Bounds())
	return buffer
}
