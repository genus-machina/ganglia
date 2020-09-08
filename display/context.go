package display

type context struct {
	content Widget
	halted  bool
	updates chan *context
}

func createContext(content Widget, updates chan *context) *context {
	c := new(context)
	c.content = content
	c.updates = updates
	return c
}

func (c *context) Content() Widget {
	return c.content
}

func (c *context) Halt() {
	c.halted = true
}

func (c *context) Render() {
	if !c.halted {
		c.updates <- c
	}
}
