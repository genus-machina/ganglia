package i2c

import (
	"github.com/genus-machina/ganglia"
)

type displayContext struct {
	content ganglia.Widget
	halted  bool
	updates chan *displayContext
}

func createDisplayContext(content ganglia.Widget, updates chan *displayContext) *displayContext {
	c := new(displayContext)
	c.content = content
	c.updates = updates
	return c
}

func (c *displayContext) Content() ganglia.Widget {
	return c.content
}

func (c *displayContext) Halt() {
	c.halted = true
}

func (c *displayContext) Render() {
	if !c.halted {
		select {
		case c.updates <- c:
		default:
		}
	}
}
