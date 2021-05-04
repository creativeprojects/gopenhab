package event

import (
	"encoding/json"

	"github.com/creativeprojects/gopenhab/api"
)

type ItemReceivedCommand struct {
	topic       string
	payload     string
	CommandType string
	Command     string
}

func NewItemReceivedCommand(topic, payload string) (*ItemReceivedCommand, error) {
	data := api.EventCommand{}
	err := json.Unmarshal([]byte(payload), &data)
	if err != nil {
		return nil, err
	}
	return &ItemReceivedCommand{
		topic:       topic,
		payload:     payload,
		CommandType: data.Type,
		Command:     data.Value,
	}, nil
}

func (i ItemReceivedCommand) Topic() string {
	return i.topic
}

func (i ItemReceivedCommand) Type() Type {
	return ItemCommand
}

type ItemReceivedUpdate struct {
	item  string
	state string
}

func NewItemReceivedUpdate(item, state string) *ItemReceivedUpdate {
	return &ItemReceivedUpdate{
		item:  item,
		state: state,
	}
}

type ItemChanged struct {
	item string
	from string
	to   string
}

func NewItemChanged(item, from, to string) *ItemChanged {
	return &ItemChanged{
		item: item,
		from: from,
		to:   to,
	}
}
