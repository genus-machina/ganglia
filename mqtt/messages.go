package mqtt

import (
	"fmt"
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

func StatusMessage(status string) string {
	payload := fmt.Sprintf("{\"status\": \"%s\"}", status)
	return payload
}
