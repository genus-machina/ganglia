package mqtt

import (
	"strings"
)

const (
	Borealis = "borealis"
	Dashi    = "dashi"
	Enrico   = "enrico"
	Tweak    = "tweak"

	DevicesPath = "devices"
	StatusPath  = "status"
)

func DeviceTopic(name string) string {
	return join("", DevicesPath, name)
}

func DeviceStatusTopic(name string) string {
	return StatusTopic(DeviceTopic(name))
}

func join(parts ...string) string {
	return strings.Join(parts, "/")
}

func StatusTopic(base string) string {
	return join(base, StatusPath)
}
