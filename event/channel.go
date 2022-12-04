package event

import "github.com/creativeprojects/gopenhab/api"

type ChannelTriggered struct {
	topic       string
	ChannelName string
	Event       string
}

func NewChannelTriggered(channelName, event string) ChannelTriggered {
	topic := channelTopicPrefix + channelName + "/" + api.TopicEventTriggered
	return ChannelTriggered{
		topic:       topic,
		ChannelName: channelName,
		Event:       event,
	}
}

func (i ChannelTriggered) Topic() string {
	return i.topic
}

func (i ChannelTriggered) Type() Type {
	return TypeChannelTriggered
}

// Verify interface
var _ Event = ChannelTriggered{}
