package event

// AliveEvent is regularly sent by openHAB 3.4+ (API v5+)
type AliveEvent struct{}

func NewAliveEvent() AliveEvent {
	return AliveEvent{}
}

func (e AliveEvent) Topic() string {
	return ""
}

func (e AliveEvent) Type() Type {
	return TypeServerAlive
}

func (e AliveEvent) TypeName() string {
	return "Alive"
}

var _ Event = &AliveEvent{}
