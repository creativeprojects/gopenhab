package openhab

import (
	"math/rand/v2"
	"time"
)

func nextBackoff(backoff time.Duration, config Config) time.Duration {
	if backoff == 0 {
		backoff = config.ReconnectionInitialBackoff
	} else {
		backoff = time.Duration(float64(backoff) * config.ReconnectionMultiplier)
	}
	if config.ReconnectionJitter > 0 {
		backoff += time.Duration(rand.Int64N(int64(config.ReconnectionJitter)*2) - int64(config.ReconnectionJitter))
	}
	if backoff > config.ReconnectionMaxBackoff {
		backoff = config.ReconnectionMaxBackoff
	}
	if backoff < config.ReconnectionMinBackoff {
		backoff = config.ReconnectionMinBackoff
	}
	return backoff
}
