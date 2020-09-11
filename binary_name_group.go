package ganglia

type BinaryNameGroup []string

func (group BinaryNameGroup) findName(name string) int {
	index := -1

	for i, n := range group {
		if name == n {
			index = i
		}
	}

	return index
}

func (group BinaryNameGroup) Value(names []string) uint {
	var value uint
	bits := len(group)

	for _, name := range names {
		index := group.findName(name)

		if index > -1 {
			shift := (bits - index - 1)
			value += (1 << shift)
		}
	}

	return value
}
