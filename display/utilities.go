package display

import (
	"image"
	"math"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
)

var (
	ttf *truetype.Font
)

func init() {
	ttf, _ = truetype.Parse(gomono.TTF)
}

func BuildFontFace(size float64) font.Face {
	options := new(truetype.Options)
	options.DPI = 133
	options.Hinting = font.HintingFull
	options.Size = size
	return truetype.NewFace(ttf, options)
}

func minInt(left, right int) int {
	return int(math.Min(float64(left), float64(right)))
}

func pointToFixed(point image.Point) fixed.Point26_6 {
	return fixed.P(point.X, point.Y)
}

func rectangleFromFixed(rectangle fixed.Rectangle26_6) image.Rectangle {
	return image.Rect(
		rectangle.Min.X.Floor(),
		rectangle.Min.Y.Floor(),
		rectangle.Max.X.Ceil(),
		rectangle.Max.Y.Ceil(),
	)
}
