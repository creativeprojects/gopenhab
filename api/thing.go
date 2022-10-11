package api

// Thing structure in the openHAB API
type Thing struct {
	UID        string          `json:"UID"`
	Label      string          `json:"label"`
	StatusInfo ThingStatusInfo `json:"statusInfo"`
}

type ThingStatusInfo struct {
	Status       string `json:"status"`
	StatusDetail string `json:"statusDetail"`
}
