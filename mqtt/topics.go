package mqtt

import (
	"strings"
)

const (
	CommandPath = "command"
	DevicesPath = "devices"
	StatusPath  = "status"
)

func CommandTopic(base ...string) string {
	parts := concat(base, []string{CommandPath})
	return JoinPaths(parts...)
}

func concat(parts ...[]string) []string {
	var result []string

	for _, part := range parts {
		result = append(result, part...)
	}

	return result
}

func DeviceTopic(name string) string {
	return JoinPaths("", DevicesPath, name)
}

func DeviceStatusTopic(name string) string {
	return StatusTopic(DeviceTopic(name))
}

func JoinPaths(parts ...string) string {
	return strings.Join(parts, "/")
}

func StatusTopic(base ...string) string {
	parts := concat(base, []string{StatusPath})
	return JoinPaths(parts...)
}
