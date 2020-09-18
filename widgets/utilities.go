package widgets

import (
	"image"
	"image/draw"
	"math"

	"github.com/golang/freetype/truetype"
	xdraw "golang.org/x/image/draw"
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

func computeImageBounds(content image.Image, bounds image.Rectangle) image.Rectangle {
	contentAspectRatio := float64(content.Bounds().Dx()) / float64(content.Bounds().Dy())
	widgetAspectRatio := float64(bounds.Dx()) / float64(bounds.Dy())

	var scale float64
	if contentAspectRatio > widgetAspectRatio {
		scale = float64(bounds.Dx()) / float64(content.Bounds().Dx())
	} else {
		scale = float64(bounds.Dy()) / float64(content.Bounds().Dy())
	}

	return image.Rect(
		bounds.Min.X,
		bounds.Min.Y,
		bounds.Min.X+int(float64(content.Bounds().Dx())*scale),
		bounds.Min.Y+int(float64(content.Bounds().Dy())*scale),
	)
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

func scaleImage(dst draw.Image, dr image.Rectangle, src image.Image, sr image.Rectangle) {
	xdraw.CatmullRom.Scale(dst, dr, src, sr, draw.Src, nil)
}
