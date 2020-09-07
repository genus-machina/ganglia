package display

import (
	"image"
	"image/gif"
	"os"
	"time"
)

type GIF struct {
	content  *gif.GIF
	index    int
	maxIndex int
}

func NewGIF(path string) (*GIF, error) {
	widget := new(GIF)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if widget.content, err = gif.DecodeAll(file); err != nil {
		return nil, err
	}

	switch widget.content.LoopCount {
	case 0:
		widget.maxIndex = -1
	case -1:
		widget.maxIndex = len(widget.content.Image)
	default:
		widget.maxIndex = len(widget.content.Image) * (widget.content.LoopCount + 1)
	}

	return widget, nil
}

func (widget *GIF) advanceFrame(rerender Trigger) {
	delay := widget.getDelay()
	time.Sleep(delay)
	widget.index++
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

func (widget *GIF) Render(bounds image.Rectangle, rerender Trigger) image.Image {
	frame := widget.getFrame()
	scaledBounds := computeImageBounds(frame, bounds)
	buffer := image.NewNRGBA(scaledBounds)
	scaleImage(buffer, scaledBounds, frame, frame.Bounds())

	if widget.maxIndex < 0 || widget.index < widget.maxIndex {
		go widget.advanceFrame(rerender)
	}

	return buffer
}
