package openhab

import (
	"net/http"
	"time"
)

type Config struct {
	// URL of your openHAB instance. It should detect automatically the REST API URL from the main URL.
	URL      string
	User     string
	Password string
	// APIToken takes precedence over User/Password authentication
	APIToken string
	// Client is optional. You can specify a custom *http.Client if you need, otherwise it's going to use the http.DefaultClient
	Client *http.Client
	// TimeoutHTTP is the maximum time allowed to send or receive commands through the openHAB API. Default is 5 seconds.
	//
	// This value is used when using the methods without passing a context (like GetItem()).
	// For all methods passing a context (like GetItemContext()), the deadline is taken from the context instead.
	TimeoutHTTP time.Duration
	// ReconnectionInitialBackoff represents how long to wait after the first failure before retrying.
	// If undefined, it defaults to 1 second
	ReconnectionInitialBackoff time.Duration
	// ReconnectionMultiplier is the factor with which to multiply backoff after a failed retry.
	// If undefined, it defaults to 2.0
	ReconnectionMultiplier float64
	// ReconnectionJitter represents by how much to randomize backoffs (+/-).
	// If undefined, it defaults to 0 (linear backoff)
	ReconnectionJitter time.Duration
	// ReconnectionMaxBackoff is the upper bound on backoff.
	// If undefined, it defaults to 1 minute
	ReconnectionMaxBackoff time.Duration
	// ReconnectionMinBackoff is the lower bound on backoff.
	// If undefined, it defaults to 1 second
	ReconnectionMinBackoff time.Duration
	// StableConnectionDuration is the time after which we consider the connection to openHAB to be stable (and resets the backoff timer).
	// If undefined, it defaults to 1 minute
	StableConnectionDuration time.Duration
	// CancellationTimeout is the time to wait for a rule to finish before sending a cancellation to its context.
	// This timeout is only used when the client is closing down.
	// If undefined, it defaults to 5 seconds
	CancellationTimeout time.Duration
	// Telemetry is used to send metrics to a monitoring system.
	// If undefined, it defaults to a no-op implementation.
	Telemetry Telemetry
}
