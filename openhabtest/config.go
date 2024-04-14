package openhabtest

type Version int

const (
	V2  Version = iota // openHAB v2.*
	V3                 // openHAB v3.* & v4.0
	V41                // openHAB v4.1+
)

// Config is the configuration object for the mock openhab Server
type Config struct {
	// Log is sending debugging information to the test logger (t.Log(), t.Logf())
	Log Logger
	// SendEventsFromAPI will send all the events automatically on the event bus when
	// adding, updating, deleting items, things, channels, etc. from the REST API.
	// The default is off meaning the events are to be sent from the unit test
	SendEventsFromAPI bool
	// Version of the openHAB server to mock: valid values are V2, V3 or V41
	Version Version
}
