package serial

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/tarm/serial"
)

type Broker struct {
	buffer []byte
	port   *serial.Port
}

func Open(port string) (*Broker, error) {
	config := &serial.Config{
		Baud:        115200,
		Name:        port,
		ReadTimeout: time.Second,
	}

	connection, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	broker := new(Broker)
	broker.buffer = make([]byte, 3)
	broker.port = connection
	return broker, nil
}

func (broker *Broker) Close() error {
	return broker.port.Close()
}

func (broker *Broker) ReadAnalogValue(pin int) (uint, error) {
	broker.buffer[0] = 'A'
	broker.buffer[1] = byte(pin)

	if _, err := broker.port.Write(broker.buffer); err != nil {
		return 0, err
	}

	if _, err := broker.port.Read(broker.buffer); err != nil {
		return 0, err
	}

	if broker.buffer[0] == 'E' {
		return 0, fmt.Errorf("Failed to read value.")
	}

	value := binary.BigEndian.Uint16(broker.buffer[1:])
	return uint(value), nil
}

func (broker *Broker) ReadDigitalValue(pin int) (bool, error) {
	broker.buffer[0] = 'D'
	broker.buffer[1] = byte(pin)

	if _, err := broker.port.Write(broker.buffer); err != nil {
		return false, err
	}

	if _, err := broker.port.Read(broker.buffer); err != nil {
		return false, err
	}

	if broker.buffer[0] == 'E' {
		return false, fmt.Errorf("Failed to read value.")
	}

	return broker.buffer[1] != 0, nil
}

func (broker *Broker) WriteDigitalValue(pin int, value bool) error {
	binary := 0

	if value {
		binary = 1
	}

	broker.buffer[0] = 'W'
	broker.buffer[1] = byte(pin)
	broker.buffer[2] = byte(binary)

	if _, err := broker.port.Write(broker.buffer); err != nil {
		return err
	}

	if _, err := broker.port.Read(broker.buffer); err != nil {
		return err
	}

	if broker.buffer[0] == 'E' {
		return fmt.Errorf("Failed to write value.")
	}

	return nil
}
