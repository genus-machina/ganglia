package widgets

import (
	"image"
	"image/draw"

	"github.com/genus-machina/ganglia"
)

// direction
const (
	Column = iota
	Row
)

// display
const (
	Fixed = iota
	Flex
)

// justify
const (
	Left = iota
	Center
	Right
)

// align
const (
	Top = iota
	Middle
	Bottom
)

type Box struct {
	children []ganglia.Widget

	align     int
	direction int
	display   int
	justify   int

	height int
	width  int
}

func NewBox() *Box {
	box := new(Box)
	return box
}

func (box *Box) advance(remaining image.Rectangle, rendered image.Image) image.Rectangle {
	var delta image.Point

	switch box.direction {
	case Column:
		delta = image.Pt(remaining.Min.X, rendered.Bounds().Max.Y+1)
	case Row:
		delta = image.Pt(rendered.Bounds().Max.X+1, remaining.Min.Y)
	}

	return remaining.Intersect(image.Rectangle{delta, remaining.Max})
}

func (box *Box) Align(align int) *Box {
	box.align = align
	return box
}

func (box *Box) Append(child ganglia.Widget) *Box {
	box.children = append(box.children, child)
	return box
}

func (box *Box) bounds(requested image.Rectangle) image.Rectangle {
	var dx, dy int

	height := requested.Dy()
	width := requested.Dx()

	if box.height > 0 {
		height = minInt(box.height, requested.Dy())
	}

	if box.width > 0 {
		width = minInt(box.width, requested.Dx())
	}

	switch box.align {
	case Bottom:
		dy = requested.Dy() - height
	case Middle:
		dy = (requested.Dy() - height) / 2
	case Top:
	}

	switch box.justify {
	case Center:
		dx = (requested.Dx() - width) / 2
	case Left:
	case Right:
		dx = requested.Dx() - width
	}

	return image.Rectangle{
		requested.Min.Add(image.Pt(dx, dy)),
		requested.Min.Add(image.Pt(dx+width, dy+height)),
	}
}

func (box *Box) Direction(direction int) *Box {
	box.direction = direction
	return box
}

func (box *Box) Display(display int) *Box {
	box.display = display
	return box
}

func (box *Box) fixedColumnPlacement(item image.Image, total, remaining image.Rectangle) image.Rectangle {
	var dx, dy int
	bounds := item.Bounds()

	switch box.align {
	case Bottom:
		dy = remaining.Dy()
	case Middle:
		dy = remaining.Dy() / 2
	case Top:
	}

	switch box.justify {
	case Center:
		dx = (total.Dx() - bounds.Dx()) / 2
	case Left:
	case Right:
		dx = total.Dx() - bounds.Dx()
	}

	return bounds.Add(image.Pt(dx, dy))
}

func (box *Box) fixedRowPlacement(item image.Image, total, remaining image.Rectangle) image.Rectangle {
	var dx, dy int
	bounds := item.Bounds()

	switch box.align {
	case Bottom:
		dy = total.Dy() - bounds.Dy()
	case Middle:
		dy = (total.Dy() - bounds.Dy()) / 2
	case Top:
	}

	switch box.justify {
	case Center:
		dx = remaining.Dx() / 2
	case Left:
	case Right:
		dx = remaining.Dx()
	}

	return bounds.Add(image.Pt(dx, dy))
}

func (box *Box) fixedPlacement(item image.Image, total, remaining image.Rectangle) image.Rectangle {
	var placement image.Rectangle

	switch box.direction {
	case Column:
		placement = box.fixedColumnPlacement(item, total, remaining)
	case Row:
		placement = box.fixedRowPlacement(item, total, remaining)
	}

	return placement
}

func (box *Box) flexBox(bounds image.Rectangle, index int) image.Rectangle {
	var flexBox image.Rectangle
	count := len(box.children)

	switch box.direction {
	case Column:
		dy := bounds.Dy() / count
		flexBox = image.Rectangle{
			bounds.Min.Add(image.Pt(0, index*dy)),
			bounds.Max.Add(image.Pt(0, (index-count+1)*dy)),
		}
	case Row:
		dx := bounds.Dx() / count
		flexBox = image.Rectangle{
			bounds.Min.Add(image.Pt(index*dx, 0)),
			bounds.Max.Add(image.Pt((index-count+1)*dx, 0)),
		}
	}

	return flexBox
}

func (box *Box) flexPlacement(item image.Image, bounds image.Rectangle) image.Rectangle {
	var dx, dy int
	placement := item.Bounds()

	switch box.align {
	case Bottom:
		dy = bounds.Dy() - placement.Dy()
	case Middle:
		dy = (bounds.Dy() - placement.Dy()) / 2
	case Top:
	}

	switch box.justify {
	case Center:
		dx = (bounds.Dx() - placement.Dx()) / 2
	case Left:
	case Right:
		dx = bounds.Dx() - placement.Dx()
	}

	return placement.Add(image.Pt(dx, dy))
}

func (box *Box) Height(height int) *Box {
	box.height = height
	return box
}

func (box *Box) Justify(justify int) *Box {
	box.justify = justify
	return box
}

func (box *Box) Render(bounds image.Rectangle, rerender ganglia.Trigger) image.Image {
	buffer := image.NewNRGBA(box.bounds(bounds))

	switch box.display {
	case Fixed:
		box.renderFixed(buffer, rerender)
	case Flex:
		box.renderFlex(buffer, rerender)
	}

	return buffer
}

func (box *Box) renderFixed(buffer draw.Image, rerender ganglia.Trigger) {
	var children []image.Image
	remaining := buffer.Bounds()
	total := remaining

	for _, child := range box.children {
		rendered := child.Render(remaining, rerender)
		remaining = box.advance(remaining, rendered)
		children = append(children, rendered)
	}

	for _, child := range children {
		placement := box.fixedPlacement(child, total, remaining)
		draw.Draw(buffer, placement, child, child.Bounds().Min, draw.Src)
	}
}

func (box *Box) renderFlex(buffer draw.Image, rerender ganglia.Trigger) {
	for index, child := range box.children {
		bounds := box.flexBox(buffer.Bounds(), index)
		rendered := child.Render(bounds, rerender)
		placement := box.flexPlacement(rendered, bounds)
		draw.Draw(buffer, placement, rendered, rendered.Bounds().Min, draw.Src)
	}
}

func (box *Box) Width(width int) *Box {
	box.width = width
	return box
}
