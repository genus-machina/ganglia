package ganglia

type DigitalOutputGroup []DigitalOutput

func (group DigitalOutputGroup) Write(value uint) {
	count := len(group)

	for index, output := range group {
		if value&(1<<(count-index-1)) > 0 {
			output <- High
		} else {
			output <- Low
		}
	}
}
