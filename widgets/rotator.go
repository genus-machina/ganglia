package widgets

import (
	"image"
	"image/color"

	"github.com/genus-machina/ganglia"
)

type Rotator struct {
	child    ganglia.Widget
	rotation ganglia.Rotation
}

func computeBounds(bounds image.Rectangle, rotation ganglia.Rotation) image.Rectangle {
	if rotation == ganglia.Rotate90 || rotation == ganglia.Rotate270 {
		return image.Rect(
			bounds.Min.X,
			bounds.Min.Y,
			bounds.Min.X+bounds.Dy(),
			bounds.Min.Y+bounds.Dx(),
		)
	}

	return bounds
}

func NewRotator(child ganglia.Widget, rotation ganglia.Rotation) *Rotator {
	rotator := new(Rotator)
	rotator.child = child
	rotator.rotation = rotation
	return rotator
}

func (rotator *Rotator) Render(bounds image.Rectangle, rerender ganglia.Trigger) image.Image {
	bounds = computeBounds(bounds, rotator.rotation)
	child := rotator.child.Render(bounds, rerender)
	return newRotatedImage(child, rotator.rotation)
}

type rotatedImage struct {
	bounds   image.Rectangle
	image    image.Image
	rotation ganglia.Rotation
}

func newRotatedImage(image image.Image, rotation ganglia.Rotation) *rotatedImage {
	rotated := new(rotatedImage)
	rotated.bounds = computeBounds(image.Bounds(), rotation)
	rotated.image = image
	rotated.rotation = rotation
	return rotated
}

func (rotated *rotatedImage) Bounds() image.Rectangle {
	return rotated.bounds
}

func (rotated *rotatedImage) ColorModel() color.Model {
	return rotated.image.ColorModel()
}

func (rotated *rotatedImage) At(x, y int) color.Color {
	tx, ty := rotated.translate(x, y)
	return rotated.image.At(tx, ty)
}

func (rotated *rotatedImage) translate(x, y int) (int, int) {
	var tx, ty int
	bounds := rotated.image.Bounds()

	switch rotated.rotation {
	default:
		fallthrough
	case ganglia.Rotate0:
		tx, ty = x, y
	case ganglia.Rotate90:
		tx, ty = y, bounds.Max.Y-x-1
	case ganglia.Rotate180:
		tx, ty = bounds.Max.X-x-1, bounds.Max.Y-y-1
	case ganglia.Rotate270:
		tx, ty = bounds.Max.X-y-1, x
	}

	return tx, ty
}
