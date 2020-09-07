package display

import (
	"image"
	"image/png"
	"os"
)

type PNG struct {
	content image.Image
}

func NewPNG(path string) (*PNG, error) {
	widget := new(PNG)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if widget.content, err = png.Decode(file); err != nil {
		return nil, err
	}
	return widget, nil
}

func (widget *PNG) Render(bounds image.Rectangle, rerender Trigger) image.Image {
	scaledBounds := computeImageBounds(widget.content, bounds)
	buffer := image.NewNRGBA(scaledBounds)
	scaleImage(buffer, scaledBounds, widget.content, widget.content.Bounds())
	return buffer
}
