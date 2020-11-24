package widgets

import (
	"image"
	"image/gif"
	"os"
	"sync"
	"time"

	"github.com/genus-machina/ganglia"
)

type GIF struct {
	content  *gif.GIF
	index    int
	maxIndex int
	once     *sync.Once
}

func NewGIF(content *gif.GIF) *GIF {
	widget := new(GIF)
	widget.content = content

	switch widget.content.LoopCount {
	case 0:
		widget.maxIndex = -1
	case -1:
		widget.maxIndex = len(widget.content.Image)
	default:
		widget.maxIndex = len(widget.content.Image) * (widget.content.LoopCount + 1)
	}

	widget.once = new(sync.Once)
	return widget
}

func OpenGIF(path string) (*GIF, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if content, err := gif.DecodeAll(file); err == nil {
		return NewGIF(content), nil
	} else {
		return nil, err
	}
}

func (widget *GIF) advanceFrame(rerender ganglia.Trigger) {
	delay := widget.getDelay()
	time.Sleep(delay)
	widget.index++
	widget.once = new(sync.Once)
	rerender()
}

func (widget *GIF) getDelay() time.Duration {
	length := len(widget.content.Delay)
	delayIndex := widget.index % length
	return time.Duration(widget.content.Delay[delayIndex]) * 10 * time.Millisecond
}

func (widget *GIF) getFrame() image.Image {
	length := len(widget.content.Image)
	frameIndex := widget.index % length
	return widget.content.Image[frameIndex]
}

func (widget *GIF) Render(bounds image.Rectangle, rerender ganglia.Trigger) image.Image {
	frame := widget.getFrame()
	scaledBounds := computeImageBounds(frame, bounds)
	buffer := image.NewNRGBA(scaledBounds)
	scaleImage(buffer, scaledBounds, frame, frame.Bounds())

	if widget.maxIndex < 0 || widget.index < widget.maxIndex {
		widget.once.Do(func() { go widget.advanceFrame(rerender) })
	}

	return buffer
}
