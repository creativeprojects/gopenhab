package openhab

import (
	"math/rand"
	"time"
)

func nextBackoff(backoff time.Duration, config Config) time.Duration {
	backoff = time.Duration(float64(backoff) * config.ReconnectionMultiplier)
	if config.ReconnectionJitter > 0 {
		backoff += time.Duration(rand.Int63n(int64(config.ReconnectionJitter)*2) - int64(config.ReconnectionJitter))
	}
	if backoff > config.ReconnectionMaxBackoff {
		backoff = config.ReconnectionMaxBackoff
	}
	if backoff < config.ReconnectionMinBackoff {
		backoff = config.ReconnectionMinBackoff
	}
	return backoff
}
