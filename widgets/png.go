package widgets

import (
	"image/png"
	"os"
)

type PNG struct {
	*Image
}

func NewPNG(path string) (*PNG, error) {
	widget := new(PNG)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if content, err := png.Decode(file); err == nil {
		widget.Image = NewImage(content)
	} else {
		return nil, err
	}

	return widget, nil
}
