package config

import "errors"

// ErrRateLimitExceeded rate limited exceeded error.
var ErrRateLimitExceeded = errors.New("rate limit exceeded")
