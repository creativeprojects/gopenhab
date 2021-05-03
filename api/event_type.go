package api

const (
	EventItemAdded              = "ItemAddedEvent"              // An item has been added to the item registry.
	EventItemRemoved            = "ItemRemovedEvent"            // An item has been removed from the item registry.
	EventItemUpdated            = "ItemUpdatedEvent"            // An item has been updated in the item registry.
	EventItemCommand            = "ItemCommandEvent"            // A command is sent to an item via a channel.
	EventItemState              = "ItemStateEvent"              // The state of an item is updated.
	EventItemStatePredicted     = "ItemStatePredictedEvent"     // The state of an item predicted to be updated.
	EventItemStateChanged       = "ItemStateChangedEvent"       // The state of an item has changed.
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
)

const (
	TopicEventAdded          = "added"          // item, thing, inbox, link
	TopicEventRemoved        = "removed"        // item, thing, inbox, link
	TopicEventUpdated        = "updated"        // item, thing, inbox
	TopicEventCommand        = "command"        // item
	TopicEventState          = "state"          // item
	TopicEventStatePredicted = "statepredicted" // item
	TopicEventStateChanged   = "statechanged"   // item
	TopicEventStatus         = "status"         // thing
	TopicEventStatusChanged  = "statuschanged"  // thing
	TopicEventTriggered      = "triggered"      // channel
)
