package api

type EventMessage struct {
	Topic   string `json:"topic"`
	Payload string `json:"payload"`
	Type    string `json:"type"`
}

type EventCommand struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type EventState struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type EventStateChanged struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	OldType  string `json:"oldType"`
	OldValue string `json:"oldValue"`
}

type EventStatePredicted struct {
	PredictedType  string `json:"predictedType"`
	PredictedValue string `json:"predictedValue"`
	IsConfirmation bool   `json:"isConfirmation"`
}

type EventStatus struct {
	Status       string `json:"status"`
	StatusDetail string `json:"statusDetail"`
}

type EventTriggered struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}
