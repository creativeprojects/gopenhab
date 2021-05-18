package openhab

import (
	"net/http"
	"time"
)

type Config struct {
	// URL of you openHAB instance. It should detect automatically the REST API URL from the main URL.
	URL                 string
	User                string
	Password            string
	BasicAuthentication bool
	APIToken            string
	// Client is optional. You can specify a custom *http.Client if you need, otherwise it's going to use the http.DefaultClient
	Client *http.Client
	// TimeoutHTTP is the maximum time allowed to send or receive commands through the openHAB API. Default is 5 seconds.
	TimeoutHTTP time.Duration
	// ReconnectionInitialBackoff represents how long to wait after the first failure before retrying.
	// If undefined, it defaults to 1 second
	ReconnectionInitialBackoff time.Duration
	// ReconnectionMultiplier is the factor with which to multiply backoff after a failed retry.
	// If undefined, it defaults to 2.0
	ReconnectionMultiplier float64
	// ReconnectionJitter represents by how much to randomize backoffs.
	// If undefined, it defaults to 0 (linear backoff)
	ReconnectionJitter float64
	// ReconnectionMaxBackoff is the upper bound on backoff.
	// If undefined, it defaults to 1 minute
	ReconnectionMaxBackoff time.Duration
	// StableConnectionDuration is the time after which we consider the connection to openHAB to be stable (and resets the backoff timer).
	// If undefined, it defaults to 1 minute
	StableConnectionDuration time.Duration
}
