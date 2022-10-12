package openhab

// see openhab documentation: https://www.openhab.org/docs/concepts/things.html#thing-status
type ThingStatus string

const (
	ThingStatusAny           ThingStatus = ""              // Trigger an event for any status
	ThingStatusUninitialized ThingStatus = "UNINITIALIZED" // This is the initial status of a Thing when it is added or the framework is being started. This status is also assigned if the initializing process failed or the binding is not available. Commands sent to Channels will not be processed.
	ThingStatusInitializing  ThingStatus = "INITIALIZING"  // This status is assigned while the binding initializes the Thing. It depends on the binding how long the initializing process takes. Commands sent to Channels will not be processed.
	ThingStatusUnknown       ThingStatus = "UNKNOWN"       // The handler is fully initialized but due to the nature of the represented device/service it cannot really tell yet whether the Thing is ONLINE or OFFLINE. Therefore the Thing potentially might be working correctly already and may or may not process commands. But the framework is allowed to send commands, because some radio-based devices may go ONLINE if a command is sent to them. The handler should take care to switch the Thing to ONLINE or OFFLINE as soon as possible.
	ThingStatusOnline        ThingStatus = "ONLINE"        // The device/service represented by a Thing is assumed to be working correctly and can process commands.
	ThingStatusOffline       ThingStatus = "OFFLINE"       // The device/service represented by a Thing is assumed to be not working correctly and may not process commands. But the framework is allowed to send commands, because some radio-based devices may go back to ONLINE if a command is sent to them.
	ThingStatusRemoving      ThingStatus = "REMOVING"      // The device/service represented by a Thing should be removed, but the binding has not confirmed the deletion yet. Some bindings need to communicate with the device to unpair it from the system. Thing is probably not working and commands cannot be processed.
	ThingStatusRemoved       ThingStatus = "REMOVED"       // This status indicates that the device/service represented by a Thing was removed from the external system after the REMOVING was initiated by the framework. Usually this status is an intermediate status because the Thing gets removed from the database after this status was assigned.
)

// see openhab documentation https://www.openhab.org/docs/concepts/things.html#status-details
type ThingStatusDetail string

const (
	ThingStatusDetailNone                        ThingStatusDetail = "NONE"                          // No further status details available.
	ThingStatusDetailHandlerMissingError         ThingStatusDetail = "HANDLER_MISSING_ERROR"         // The handler cannot be initialized because the responsible binding is not available or started.
	ThingStatusDetailHandlerRegisteringError     ThingStatusDetail = "HANDLER_REGISTERING_ERROR"     // The handler failed in the service registration phase.
	ThingStatusDetailHandlerConfigurationPending ThingStatusDetail = "HANDLER_CONFIGURATION_PENDING" // The handler is registered but cannot be initialized because of missing configuration parameters.
	ThingStatusDetailHandlerInitializingError    ThingStatusDetail = "HANDLER_INITIALIZING_ERROR"    // The handler failed in the initialization phase.
	ThingStatusDetailBridgeUninitialized         ThingStatusDetail = "BRIDGE_UNINITIALIZED"          // The bridge associated with this Thing is not initialized.
	ThingStatusDetailDisabled                    ThingStatusDetail = "DISABLED"                      // The thing was explicitly disabled.
	ThingStatusDetailConfigurationPending        ThingStatusDetail = "CONFIGURATION_PENDING"         // The Thing is waiting to transfer configuration information to a device. Some bindings need to communicate with the device to make sure the configuration is accepted.
	ThingStatusDetailCommunicationError          ThingStatusDetail = "COMMUNICATION_ERROR"           // Error communicating with the device. This may be only a temporary error.
	ThingStatusDetailConfigurationError          ThingStatusDetail = "CONFIGURATION_ERROR"           // An issue with the configuration of a Thing prevents communication with the represented device or service. This issue might be solved by reconfiguring the Thing.
	ThingStatusDetailBridgeOffline               ThingStatusDetail = "BRIDGE_OFFLINE"                // Assuming the Thing to be offline because the corresponding bridge is offline.
	ThingStatusDetailFirmwareUpdating            ThingStatusDetail = "FIRMWARE_UPDATING"             // The Thing is currently undergoing a firmware update.
	ThingStatusDetailDutyCycle                   ThingStatusDetail = "DUTY_CYCLE"                    // The Thing is currently in DUTY_CYCLE state, which means it is blocked for further usage.
	ThingStatusDetailGone                        ThingStatusDetail = "GONE"                          // The Thing has been removed from the bridge or the network to which it belongs and is no longer available for use. The user can now remove the Thing from the system.
)
