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
	count            int
	bounds           image.Rectangle
	mutex            sync.Mutex
	remainingRenders int
	t                *testing.T
	widget           Widget
}

func NewTestDisplay(t *testing.T) *TestDisplay {
	display := new(TestDisplay)
	display.bounds = image.Rect(0, 0, 128, 64)
	display.t = t
	return display
}

func (display *TestDisplay) fileName() string {
	display.mutex.Lock()
	defer display.mutex.Unlock()

	name := fmt.Sprintf("%s_%d.png", display.t.Name(), display.count)
	display.count++
	return name
}

func (display *TestDisplay) Render(content Widget) {
	file, err := os.Create(display.fileName())
	if err != nil {
		display.t.Errorf("Failed to create image. %s.", err.Error())
	}
	defer file.Close()

	display.widget = content

	image := content.Render(display.bounds, display.rerender)
	if err := png.Encode(file, image); err != nil {
		display.t.Errorf("Failed to save image. %s.", err.Error())
	}
}

func (display *TestDisplay) rerender() {
	if display.remainingRenders > 0 {
		display.Render(display.widget)
	}

	display.remainingRenders--
}
