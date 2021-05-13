package event

import "fmt"

type subscription struct {
	id        int
	name      string
	eventType Type
	callback  func(e Event)
}

func (s subscription) String() string {
	return fmt.Sprintf("id=%d; name=%q, eventType=%q", s.id, s.name, s.eventType)
}
