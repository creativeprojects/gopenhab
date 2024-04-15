package event

import (
	"fmt"
	"strings"
)

// StartlevelEvent is triggered by the openhab server on startup (typically from 30 to 100).
// This event was introduced in API version 5.
type StartlevelEvent struct {
	topic string
	level int
}

func NewStartlevelEvent(topic string, startlevel int) StartlevelEvent {
	topic = strings.TrimPrefix(topic, "openhab/")
	return StartlevelEvent{
		topic: topic,
		level: startlevel,
	}
}

func (e StartlevelEvent) Topic() string {
	return e.topic
}

func (e StartlevelEvent) Type() Type {
	return TypeServerStartlevel
}

func (e StartlevelEvent) Level() int {
	return e.level
}

func (e StartlevelEvent) String() string {
	return fmt.Sprintf("Received start level %d from server", e.level)
}

var _ Event = &StartlevelEvent{}
