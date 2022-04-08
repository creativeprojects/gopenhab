package openhab

type ClientState int

const (
	StateStarting ClientState = iota
	StateConnecting
	StateConnected
	StateDisconnected
)
