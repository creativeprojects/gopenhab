package openhab

import (
	"net/http"
)

type Config struct {
	URL                 string
	User                string
	Password            string
	BasicAuthentication bool
	APIToken            string
	Client              *http.Client
}
