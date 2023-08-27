package internal

type RestIndex struct {
	APIVersion        string              `json:"version"`
	Locale            string              `json:"locale"`
	MeasurementSystem string              `json:"measurementSystem"`
	RuntimeInfo       RuntimeInfo         `json:"runtimeInfo"`
	Links             []map[string]string `json:"links"`
}

type RuntimeInfo struct {
	Version     string `json:"version"`
	BuildString string `json:"buildString"`
}
