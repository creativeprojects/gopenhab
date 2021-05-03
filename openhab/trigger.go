package openhab

type Trigger interface {
	// activate the trigger for func() in the context of a *Client
	activate(client *Client, run func()) error
}
