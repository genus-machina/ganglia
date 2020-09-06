package mqtt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/genus-machina/ganglia"
)

const (
	AtMostOnce = iota
	AtLeastOnce
	ExactlyOnce
)

const (
	Offline = "offline"
	Online  = "online"
)

const (
	timeFormat = time.RFC3339
)

type AnalogEvent struct {
	T time.Time `json:"time"`
	V uint      `json:"value"`
}

func encodeAnalogEvent(event ganglia.AnalogEvent) []byte {
	encodableEvent := &AnalogEvent{event.Time(), uint(event.Value())}
	payload, _ := json.Marshal(encodableEvent)
	return payload
}

func parseAnalogEvent(message mqtt.Message) (*AnalogEvent, error) {
	event := new(AnalogEvent)
	err := json.Unmarshal(message.Payload(), event)
	return event, err
}

func (event *AnalogEvent) Encode() []byte {
	payload, _ := json.Marshal(event)
	return payload
}

func (event *AnalogEvent) Time() time.Time {
	return event.T
}

func (event *AnalogEvent) Value() ganglia.AnalogValue {
	return ganglia.AnalogValue(event.V)
}

type DigitalEvent struct {
	L bool      `json:"level"`
	T time.Time `json:"time"`
}

func encodeDigitalEvent(event ganglia.DigitalEvent) []byte {
	encodableEvent := &DigitalEvent{bool(event.Value()), event.Time()}
	payload, _ := json.Marshal(encodableEvent)
	return payload
}

func parseDigitalEvent(message mqtt.Message) (*DigitalEvent, error) {
	event := new(DigitalEvent)
	err := json.Unmarshal(message.Payload(), event)
	return event, err
}

func (event *DigitalEvent) Time() time.Time {
	return event.T
}

func (event *DigitalEvent) Value() ganglia.DigitalValue {
	return ganglia.DigitalValue(event.L)
}

func StatusMessage(status string) string {
	payload := fmt.Sprintf("{\"status\": \"%s\"}", status)
	return payload
}
