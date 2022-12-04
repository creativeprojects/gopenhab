package event

import (
	"fmt"
	"strings"

	"github.com/creativeprojects/gopenhab/api"
)

const (
	itemTopicPrefix    = "items/"
	thingTopicPrefix   = "things/"
	channelTopicPrefix = "channels/"
)

type Type int

const (
	TypeUnknown Type = iota
	TypeClientStarted
	TypeClientConnected
	TypeClientConnectionStable
	TypeClientDisconnected
	TypeClientStopped
	TypeClientError
	TypeServerAlive            // API version >=5 sends ALIVE messages (every minute or so)
	TypeServerStartlevel       // API version >=5 sends Startlevel events during startup (typically from 30 to 100)
	TypeTimeCron               // On a specific date and/or time
	TypeItemAdded              // An item has been added to the item registry.
	TypeItemRemoved            // An item has been removed from the item registry.
	TypeItemUpdated            // An item has been updated in the item registry.
	TypeItemCommand            // A command is sent to an item via a channel.
	TypeItemState              // The state of an item is updated.
	TypeItemStatePredicted     // The state of an item predicted to be updated.
	TypeItemStateChanged       // The state of an item has changed.
	TypeGroupItemStateChanged  // The state of a group item has changed through a member.
	TypeThingAdded             // A thing has been added to the thing registry.
	TypeThingRemoved           // A thing has been removed from the thing registry.
	TypeThingUpdated           // A thing has been updated in the thing registry.
	TypeThingStatusInfo        // The status of a thing is updated.
	TypeThingStatusInfoChanged // The status of a thing changed.
	TypeInboxAdded             // A discovery result has been added to the inbox.
	TypeInboxRemoved           // A discovery result has been removed from the inbox.
	TypeInboxUpdate            // A discovery result has been updated in the inbox.
	TypeItemChannelLinkAdded   // An item channel link has been added to the registry.
	TypeItemChannelLinkRemoved // An item channel link has been removed from the registry.
	TypeChannelTriggered       // A channel has been triggered.
)

// Match returns true if the name matches the topic
func (t Type) Match(topic, name string) bool {
	switch t {
	case TypeUnknown, TypeClientStarted, TypeClientConnected, TypeClientConnectionStable,
		TypeClientDisconnected, TypeClientStopped, TypeClientError, TypeTimeCron:
		return true
	case TypeItemAdded:
		return topic == itemTopicPrefix+name+"/"+api.TopicEventAdded
	case TypeItemRemoved:
		return topic == itemTopicPrefix+name+"/"+api.TopicEventRemoved
	case TypeItemUpdated:
		return topic == itemTopicPrefix+name+"/"+api.TopicEventUpdated
	case TypeItemCommand:
		return topic == itemTopicPrefix+name+"/"+api.TopicEventCommand
	case TypeItemState:
		return topic == itemTopicPrefix+name+"/"+api.TopicEventState
	case TypeItemStateChanged:
		return topic == itemTopicPrefix+name+"/"+api.TopicEventStateChanged
	case TypeGroupItemStateChanged:
		return strings.HasPrefix(topic, itemTopicPrefix+name+"/") &&
			strings.HasSuffix(topic, "/"+api.TopicEventStateChanged)
	case TypeThingStatusInfo:
		return topic == thingTopicPrefix+name+"/"+api.TopicEventStatus
	default:
		panic(fmt.Sprintf("event.Type %d Match undefined", t))
	}
}
