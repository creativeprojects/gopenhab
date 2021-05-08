package event

type Type int

const (
	Unknown Type = iota
	ClientConnected
	ClientDisconnected
	TimeCron
	ItemAdded              // An item has been added to the item registry.
	ItemRemoved            // An item has been removed from the item registry.
	ItemUpdated            // An item has been updated in the item registry.
	ItemCommand            // A command is sent to an item via a channel.
	ItemState              // The state of an item is updated.
	ItemStatePredicted     // The state of an item predicted to be updated.
	ItemStateChanged       // The state of an item has changed.
	GroupItemStateChanged  // The state of a group item has changed through a member.
	ThingAdded             // A thing has been added to the thing registry.
	ThingRemoved           // A thing has been removed from the thing registry.
	ThingUpdated           // A thing has been updated in the thing registry.
	ThingStatusInfo        // The status of a thing is updated.
	ThingStatusInfoChanged // The status of a thing changed.
	InboxAdded             // A discovery result has been added to the inbox.
	InboxRemoved           // A discovery result has been removed from the inbox.
	InboxUpdate            // A discovery result has been updated in the inbox.
	ItemChannelLinkAdded   // An item channel link has been added to the registry.
	ItemChannelLinkRemoved // An item channel link has been removed from the registry.
	ChannelTriggered       // A channel has been triggered.
)
