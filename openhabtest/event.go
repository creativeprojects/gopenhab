package openhabtest

import (
	"encoding/json"
	"fmt"

	"github.com/creativeprojects/gopenhab/api"
	"github.com/creativeprojects/gopenhab/event"
)

// EventString returns the event topic and event string to send to the event bus
func EventString(e event.Event, prefix string) (string, string) {
	if e == nil {
		return "", ""
	}
	switch ev := e.(type) {
	case event.ItemReceivedCommand:
		rawPayload, err := json.Marshal(api.EventCommand{
			Type:  ev.CommandType,
			Value: ev.Command,
		})
		if err != nil {
			panic(err)
		}
		topic := prefix + ev.Topic()
		rawEvent, err := json.Marshal(api.EventMessage{
			Topic:   topic,
			Payload: string(rawPayload),
			Type:    api.EventItemCommand,
		})
		if err != nil {
			panic(err)
		}
		return topic, string(rawEvent)

	case event.ItemReceivedState:
		rawPayload, err := json.Marshal(api.EventState{
			Type:  ev.StateType,
			Value: ev.State,
		})
		if err != nil {
			panic(err)
		}
		topic := prefix + ev.Topic()
		rawEvent, err := json.Marshal(api.EventMessage{
			Topic:   topic,
			Payload: string(rawPayload),
			Type:    api.EventItemState,
		})
		if err != nil {
			panic(err)
		}
		return topic, string(rawEvent)

	case event.ItemStateChanged:
		rawPayload, err := json.Marshal(api.EventStateChanged{
			Type:     ev.NewStateType,
			Value:    ev.NewState,
			OldType:  ev.PreviousStateType,
			OldValue: ev.PreviousState,
		})
		if err != nil {
			panic(err)
		}
		topic := prefix + ev.Topic()
		rawEvent, err := json.Marshal(api.EventMessage{
			Topic:   topic,
			Payload: string(rawPayload),
			Type:    api.EventItemStateChanged,
		})
		if err != nil {
			panic(err)
		}
		return topic, string(rawEvent)

	default:
		panic(fmt.Sprintf("event type %d not (yet) handled", e.Type()))
	}
}
