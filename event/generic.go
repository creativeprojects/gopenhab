package event

// GenericEvent is used when its type is unknown
type GenericEvent struct {
	typeName string
	topic    string
	payload  string
}

func (e GenericEvent) Topic() string {
	return e.topic
}

func (e GenericEvent) Type() Type {
	return Unknown
}

func (e GenericEvent) Payload() string {
	return e.payload
}

func (e GenericEvent) TypeName() string {
	return e.typeName
}

var _ Event = &GenericEvent{}
