package api

// Thing structure in the openHAB API
type Thing struct {
	UID           string            `json:"UID"`
	Label         string            `json:"label"`
	StatusInfo    ThingStatusInfo   `json:"statusInfo"`
	BridgeUID     string            `json:"bridgeUID"`
	Configuration map[string]any    `json:"configuration"`
	Properties    map[string]string `json:"properties"`
	ThingTypeUID  string            `json:"thingTypeUID"`
}

type ThingStatusInfo struct {
	Status       string `json:"status"`
	StatusDetail string `json:"statusDetail"`
	Description  string `json:"description"`
}
