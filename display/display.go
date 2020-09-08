package display

import (
	"image"
)

type Display interface {
	Halt()
	Render(content Widget)
}

type Trigger func()

type Widget interface {
	Render(bounds image.Rectangle, rerender Trigger) image.Image
}
