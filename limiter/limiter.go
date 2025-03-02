package limiter

import (
	"context"
	"sort"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter interface {
	Wait(ctx context.Context) error
	Limit() rate.Limit
}

func Per(enventCount int, duration time.Duration) rate.Limit {
	return rate.Every(duration / time.Duration(enventCount))
}

func NewMultiLimit(limiters ...RateLimiter) *multiLimiter {
	byLimit := func(i, j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters, byLimit)

	return &multiLimiter{limiters}
}

type multiLimiter struct {
	limiters []RateLimiter
}

func (l *multiLimiter) Wait(ctx context.Context) error {
	for _, l := range l.limiters {
		if err := l.Wait(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (l *multiLimiter) Limit() rate.Limit {
	return l.limiters[0].Limit()
}
