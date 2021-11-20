package util

import (
	"errors"
	"time"
)

type RateLimiter interface {
	Run(fn func() error) error
}

type rateLimiter struct {
	ticker     *time.Ticker
	startTime  time.Time
	limit      int
	totalCalls int
}

func NewRateLimiter(limit int) *rateLimiter {
	return &rateLimiter{
		ticker:     time.NewTicker(time.Duration(60/limit) * time.Second),
		startTime:  time.Now(),
		limit:      limit,
		totalCalls: 0,
	}
}

func (r *rateLimiter) Run(fn func() error) error {
	err := make(chan error, 1)
	defer close(err)

	for range r.ticker.C {
		ok := r.increment()
		if !ok {
			err <- errors.New("reach the rate limit number")
			break
		}

		err <- fn()
		break
	}

	return <-err
}

func (r *rateLimiter) increment() bool {
	if time.Now().Minute() > r.startTime.Minute() {
		r.startTime = time.Now()
		r.totalCalls = 0
	} else if r.totalCalls == r.limit {
		return false
	}

	r.totalCalls++
	return true
}
