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
	// DelayBeforeReconnecting represents a delay between the time you lost the connection to openHAB and the time you try to reconnect.
	// If undefined, it defaults to 1 second
	DelayBeforeReconnecting time.Duration
}
