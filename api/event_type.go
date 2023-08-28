package api

const (
	EventItemAdded              = "ItemAddedEvent"              // An item has been added to the item registry.
	EventItemRemoved            = "ItemRemovedEvent"            // An item has been removed from the item registry.
	EventItemUpdated            = "ItemUpdatedEvent"            // An item has been updated in the item registry.
	EventItemCommand            = "ItemCommandEvent"            // A command is sent to an item via a channel.
	EventItemState              = "ItemStateEvent"              // The state of an item is updated.
	EventItemStateUpdated       = "ItemStateUpdatedEvent"       // The state of an item is updated - since OH 4.0
	EventItemStatePredicted     = "ItemStatePredictedEvent"     // The state of an item predicted to be updated.
	EventItemStateChanged       = "ItemStateChangedEvent"       // The state of an item has changed.
	EventGroupItemStateUpdated  = "GroupStateUpdatedEvent"      // The state of a group of items has been updated through a member.
	EventGroupItemStateChanged  = "GroupItemStateChangedEvent"  // The state of a group item has changed through a member.
	EventThingAdded             = "ThingAddedEvent"             // A thing has been added to the thing registry.
	EventThingRemoved           = "ThingRemovedEvent"           // A thing has been removed from the thing registry.
	EventThingUpdated           = "ThingUpdatedEvent"           // A thing has been updated in the thing registry.
	EventThingStatusInfo        = "ThingStatusInfoEvent"        // The status of a thing is updated.
	EventThingStatusInfoChanged = "ThingStatusInfoChangedEvent" // The status of a thing changed.
	EventInboxAdded             = "InboxAddedEvent"             // A discovery result has been added to the inbox.
	EventInboxRemoved           = "InboxRemovedEvent"           // A discovery result has been removed from the inbox.
	EventInboxUpdate            = "InboxUpdateEvent"            // A discovery result has been updated in the inbox.
	EventItemChannelLinkAdded   = "ItemChannelLinkAddedEvent"   // An item channel link has been added to the registry.
	EventItemChannelLinkRemoved = "ItemChannelLinkRemovedEvent" // An item channel link has been removed from the registry.
	EventChannelTriggered       = "ChannelTriggeredEvent"       // A channel has been triggered.
	// event added in API v5
	EventTypeAlive      = "ALIVE"           // API version >=5 sends ALIVE events (every minute or so)
	EventTypeStartlevel = "StartlevelEvent" // Event sent during server startup (typically from 30 to 100)
)

const (
	TopicEventAdded          = "added"          // item, thing, inbox, link
	TopicEventRemoved        = "removed"        // item, thing, inbox, link
	TopicEventUpdated        = "updated"        // item, thing, inbox
	TopicEventCommand        = "command"        // item
	TopicEventState          = "state"          // item
	TopicEventStateUpdated   = "stateupdated"   // item
	TopicEventStatePredicted = "statepredicted" // item
	TopicEventStateChanged   = "statechanged"   // item
	TopicEventStatus         = "status"         // thing
	TopicEventStatusChanged  = "statuschanged"  // thing
	TopicEventTriggered      = "triggered"      // channel
)
