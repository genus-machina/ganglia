package ganglia

type DigitalOutput chan<- DigitalValue

func (output DigitalOutput) Close() {
	close(output)
}

func (output DigitalOutput) Invert() DigitalOutput {
	inverted := make(chan DigitalValue)
	go output.invert(inverted)
	return inverted
}

func (output DigitalOutput) invert(inverted <-chan DigitalValue) {
	defer close(output)

	for value := range inverted {
		output <- !value
	}
}

func (output DigitalOutput) Write(value DigitalValue) {
	output <- value
}
