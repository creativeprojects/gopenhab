package event

import "fmt"

type subscription struct {
	id        int
	topic     string
	eventType Type
	callback  func(e Event)
}

func (s subscription) String() string {
	return fmt.Sprintf("id=%d; topic=%q, eventType=%q", s.id, s.topic, s.eventType)
}
