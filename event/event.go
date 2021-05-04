package event

type Event interface {
	Topic() string
	Type() Type
}
