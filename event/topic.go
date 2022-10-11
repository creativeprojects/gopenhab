package event

import "strings"

// splitItemTopic returns the item name, triggering item (if any) and the event type
func splitItemTopic(topic string) (string, string, string) {
	return splitTopic(topic, "items")
}

// splitThingTopic returns the thing name and the event type
func splitThingTopic(topic string) (string, string) {
	name, _, evType := splitTopic(topic, "things")
	return name, evType
}

func splitTopic(topic, collection string) (string, string, string) {
	// "smarthome" was used in openHAB 2.x
	// "openhab" is used since openHAB 3.0
	parts := strings.Split(topic, "/")
	if len(parts) < 4 || len(parts) > 5 ||
		(parts[0] != "smarthome" && parts[0] != "openhab") ||
		parts[1] != collection {
		return "", "", ""
	}
	if len(parts) == 5 {
		return parts[2], parts[3], parts[4]
	}
	return parts[2], "", parts[3]
}
