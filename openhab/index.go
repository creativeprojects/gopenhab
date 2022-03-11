package openhab

type restIndex struct {
	APIVersion        string              `json:"version"`
	Locale            string              `json:"locale"`
	MeasurementSystem string              `json:"measurementSystem"`
	Links             []map[string]string `json:"links"`
}
