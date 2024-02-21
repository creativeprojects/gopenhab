package event

import "fmt"

type subscription struct {
	id        int
	name      string
	eventType Type
	callback  func(e Event)
	once      bool
}

func (s subscription) String() string {
	return fmt.Sprintf("id=%d; name=%q, eventType=%q, once=%t", s.id, s.name, s.eventType, s.once)
}
