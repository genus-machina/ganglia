package widgets

import (
	"image/png"
	"os"
)

func OpenPNG(path string) (*Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if content, err := png.Decode(file); err == nil {
		return NewImage(content), nil
	} else {
		return nil, err
	}
}
