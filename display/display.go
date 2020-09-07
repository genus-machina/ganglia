package display

import (
	"image"
)

type Display interface {
	Render(content Widget)
}

type Trigger func()

type Widget interface {
	Render(bounds image.Rectangle, rerender Trigger) image.Image
}
