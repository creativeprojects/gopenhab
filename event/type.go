package event

type Type int

const (
	TypeUnknown Type = iota
	TypeClientConnected
	TypeClientDisconnected
	TypeTimeCron
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
