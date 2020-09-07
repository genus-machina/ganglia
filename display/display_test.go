package display

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"sync"
	"testing"
)

type TestDisplay struct {
	count     int
	bounds    image.Rectangle
	frames    sync.WaitGroup
	mutex     sync.Mutex
	rendering bool
	t         *testing.T
}

func NewTestDisplay(t *testing.T) *TestDisplay {
	display := new(TestDisplay)
	display.bounds = image.Rect(0, 0, 128, 64)
	display.t = t
	return display
}

func (display *TestDisplay) Bounds(bounds image.Rectangle) *TestDisplay {
	display.bounds = bounds
	return display
}

func (display *TestDisplay) fileName() string {
	display.mutex.Lock()
	defer display.mutex.Unlock()

	name := fmt.Sprintf("%s_%d.png", display.t.Name(), display.count)
	display.count++
	return name
}

func (display *TestDisplay) MaxFrames(count int) *TestDisplay {
	display.frames.Add(count - 1)
	return display
}

func (display *TestDisplay) Render(content Widget) {
	display.start()
	display.render(content, display.rerender(content))
	display.frames.Wait()
	display.stop()
}

func (display *TestDisplay) render(content Widget, rerender Trigger) {
	file, err := os.Create(display.fileName())
	if err != nil {
		display.t.Errorf("Failed to create image. %s.", err.Error())
	}
	defer file.Close()

	image := content.Render(display.bounds, rerender)
	if err := png.Encode(file, image); err != nil {
		display.t.Errorf("Failed to save image. %s.", err.Error())
	}
}

func (display *TestDisplay) rerender(content Widget) Trigger {
	return func() {
		if display.running() {
			display.render(content, display.rerender(content))
			display.frames.Done()
		}
	}
}

func (display *TestDisplay) running() bool {
	return display.rendering
}

func (display *TestDisplay) start() {
	display.rendering = true
}

func (display *TestDisplay) stop() {
	display.rendering = false
}
