package mqtt

import (
	"strings"
)

const (
	DevicesPath = "devices"
	StatusPath  = "status"
)

func DeviceTopic(name string) string {
	return JoinPaths("", DevicesPath, name)
}

func DeviceStatusTopic(name string) string {
	return StatusTopic(DeviceTopic(name))
}

func JoinPaths(parts ...string) string {
	return strings.Join(parts, "/")
}

func StatusTopic(base string) string {
	return JoinPaths(base, StatusPath)
}
