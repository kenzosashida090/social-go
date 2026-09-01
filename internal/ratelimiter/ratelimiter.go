package ratelimiter

import "time"

type Limiter interface {
	Allow(count int64, limit int64) (bool, time.Duration)
}

type Config struct {
	RequestPerTimeFrame int
	TimeFrame           time.Duration
	Enabled             bool
}
