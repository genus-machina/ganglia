package widgets

import (
	"testing"

	"github.com/genus-machina/ganglia/widgets/test"
)

func TestBoxColumnFixedAlignBottom(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(10)

	box := NewBox().
		Align(Bottom).
		Append(NewText(face, "one")).
		Append(NewText(face, "two")).
		Append(NewText(face, "three")).
		Direction(Column).
		Display(Fixed)

	display.Render(box)
}

func TestBoxColumnFixedAlignMiddle(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(10)

	box := NewBox().
		Align(Middle).
		Append(NewText(face, "one")).
		Append(NewText(face, "two")).
		Append(NewText(face, "three")).
		Direction(Column).
		Display(Fixed)

	display.Render(box)
}

func TestBoxColumnFixedAlignTop(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(10)

	box := NewBox().
		Align(Top).
		Append(NewText(face, "one")).
		Append(NewText(face, "two")).
		Append(NewText(face, "three")).
		Direction(Column).
		Display(Fixed)

	display.Render(box)
}

func TestBoxColumnFixedJustifyLeft(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(10)

	box := NewBox().
		Append(NewText(face, "dog")).
		Append(NewText(face, "bird")).
		Append(NewText(face, "elephant")).
		Direction(Column).
		Display(Fixed).
		Justify(Left)

	display.Render(box)
}

func TestBoxColumnFixedJustifyCenter(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(10)

	box := NewBox().
		Append(NewText(face, "dog")).
		Append(NewText(face, "bird")).
		Append(NewText(face, "elephant")).
		Direction(Column).
		Display(Fixed).
		Justify(Center)

	display.Render(box)
}

func TestBoxColumnFixedJustifyRight(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(10)

	box := NewBox().
		Append(NewText(face, "dog")).
		Append(NewText(face, "bird")).
		Append(NewText(face, "elephant")).
		Direction(Column).
		Display(Fixed).
		Justify(Right)

	display.Render(box)
}

func TestBoxColumnFlexAlignTop(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(10)

	box := NewBox().
		Align(Top).
		Append(NewText(face, "one")).
		Append(NewText(face, "two")).
		Append(NewText(face, "three")).
		Direction(Column).
		Display(Flex)

	display.Render(box)
}

func TestBoxColumnFlexAlignBottom(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(10)

	box := NewBox().
		Align(Bottom).
		Append(NewText(face, "one")).
		Append(NewText(face, "two")).
		Append(NewText(face, "three")).
		Direction(Column).
		Display(Flex)

	display.Render(box)
}

func TestBoxColumnFlexAlignMiddle(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(10)

	box := NewBox().
		Align(Middle).
		Append(NewText(face, "one")).
		Append(NewText(face, "two")).
		Append(NewText(face, "three")).
		Direction(Column).
		Display(Flex)

	display.Render(box)
}

func TestBoxColumnOverflow(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(20)

	box := NewBox().
		Append(NewText(face, "one")).
		Append(NewText(face, "two")).
		Append(NewText(face, "three")).
		Direction(Column)

	display.Render(box)
}

func TestBoxRowOverflow(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(10)

	box := NewBox().
		Append(NewText(face, "Hello")).
		Append(NewText(face, "Amazing")).
		Append(NewText(face, "World")).
		Direction(Row)

	display.Render(box)
}

func TestBoxRowFixedAlignTop(t *testing.T) {
	display := test.NewTestDisplay(t)

	box := NewBox().
		Align(Top).
		Append(NewText(BuildFontFace(8), "A")).
		Append(NewText(BuildFontFace(16), "B")).
		Append(NewText(BuildFontFace(24), "C")).
		Direction(Row).
		Display(Fixed)

	display.Render(box)
}

func TestBoxRowFixedAlignBottom(t *testing.T) {
	display := test.NewTestDisplay(t)

	box := NewBox().
		Align(Bottom).
		Append(NewText(BuildFontFace(8), "A")).
		Append(NewText(BuildFontFace(16), "B")).
		Append(NewText(BuildFontFace(24), "C")).
		Direction(Row).
		Display(Fixed)

	display.Render(box)
}

func TestBoxRowFixedAlignMiddle(t *testing.T) {
	display := test.NewTestDisplay(t)

	box := NewBox().
		Align(Middle).
		Append(NewText(BuildFontFace(8), "A")).
		Append(NewText(BuildFontFace(16), "B")).
		Append(NewText(BuildFontFace(24), "C")).
		Direction(Row).
		Display(Fixed)

	display.Render(box)
}

func TestBoxRowFixedJustifyLeft(t *testing.T) {
	display := test.NewTestDisplay(t)

	box := NewBox().
		Append(NewText(BuildFontFace(8), "A")).
		Append(NewText(BuildFontFace(16), "B")).
		Append(NewText(BuildFontFace(24), "C")).
		Direction(Row).
		Display(Fixed).
		Justify(Left)

	display.Render(box)
}

func TestBoxRowFixedJustifyRight(t *testing.T) {
	display := test.NewTestDisplay(t)

	box := NewBox().
		Append(NewText(BuildFontFace(8), "A")).
		Append(NewText(BuildFontFace(16), "B")).
		Append(NewText(BuildFontFace(24), "C")).
		Direction(Row).
		Display(Fixed).
		Justify(Right)

	display.Render(box)
}

func TestBoxRowFixedJustifyCenter(t *testing.T) {
	display := test.NewTestDisplay(t)

	box := NewBox().
		Append(NewText(BuildFontFace(8), "A")).
		Append(NewText(BuildFontFace(16), "B")).
		Append(NewText(BuildFontFace(24), "C")).
		Direction(Row).
		Display(Fixed).
		Justify(Center)

	display.Render(box)
}

func TestBoxRowFlexJustifyLeft(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(10)

	box := NewBox().
		Append(NewText(face, "A")).
		Append(NewText(face, "B")).
		Append(NewText(face, "C")).
		Direction(Row).
		Display(Flex).
		Justify(Left)

	display.Render(box)
}

func TestBoxRowFlexJustifyCenter(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(10)

	box := NewBox().
		Append(NewText(face, "A")).
		Append(NewText(face, "B")).
		Append(NewText(face, "C")).
		Direction(Row).
		Display(Flex).
		Justify(Center)

	display.Render(box)
}

func TestBoxRowFlexJustifyRight(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(10)

	box := NewBox().
		Append(NewText(face, "A")).
		Append(NewText(face, "B")).
		Append(NewText(face, "C")).
		Direction(Row).
		Display(Flex).
		Justify(Right)

	display.Render(box)
}

func TestBoxDimensions(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(20)

	box := NewBox().
		Align(Middle).
		Append(NewText(face, "Hello")).
		Height(10).
		Justify(Center).
		Width(10)

	display.Render(box)
}

func TestBoxDimensionsHeight(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(20)

	box := NewBox().
		Append(NewText(face, "Hello")).
		Height(10).
		Justify(Center)

	display.Render(box)
}

func TestBoxDimensionsWidth(t *testing.T) {
	display := test.NewTestDisplay(t)
	face := BuildFontFace(20)

	box := NewBox().
		Align(Middle).
		Append(NewText(face, "Hello")).
		Justify(Center).
		Width(10)

	display.Render(box)
}
