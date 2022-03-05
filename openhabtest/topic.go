package openhabtest

func topicPrefix(version Version) string {
	prefix := "smarthome/"
	if version >= V3 {
		prefix = "openhab/"
	}
	return prefix
}
