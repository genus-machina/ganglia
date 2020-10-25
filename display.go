package ganglia

import (
	"image"
)

type Display interface {
	Halt()
	Render(content Widget)
	Rotate(rotation Rotation)
}

type Rotation int

const (
	Rotate0 Rotation = iota
	Rotate90
	Rotate180
	Rotate270
)

type Trigger func()

type Widget interface {
	Render(bounds image.Rectangle, rerender Trigger) image.Image
}
