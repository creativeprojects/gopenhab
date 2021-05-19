package openhabtest

// Config is the configuration object for the mock openhab Server
type Config struct {
	// Log is sending debugging information to the test logger (t.Log())
	Log Logger
	// SendEvents will send all the events automatically on the event bus when adding, updating, deleting items, things, channels, etc.
	// The default is off meaning the events are to be sent from the unit test
	SendEvents bool
}
