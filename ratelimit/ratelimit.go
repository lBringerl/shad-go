//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	newRequestCh chan context.Context
	maxCount     int
	interval     time.Duration
	stopped      bool
	abortCh      chan struct{}
}

var ErrStopped = errors.New("limiter stopped")

func StartRequestMonitor(limiter *Limiter) {
	rQueue := make([]<-chan time.Time, 0)
	var lastReqTimer <-chan time.Time

	for {
		if len(rQueue) >= limiter.maxCount {
			select {
			case <-lastReqTimer:
				rQueue = rQueue[1:]
				if len(rQueue) > 0 {
					lastReqTimer = rQueue[0]
				}
				continue
			case <-limiter.abortCh:
				close(limiter.abortCh)
				return
			}
		}

		select {
		case <-limiter.newRequestCh:
			rQueue = append(rQueue, time.NewTimer(limiter.interval).C)
			lastReqTimer = rQueue[0]
		case <-lastReqTimer:
			rQueue = rQueue[1:]
			if len(rQueue) > 0 {
				lastReqTimer = rQueue[0]
			}
		case <-limiter.abortCh:
			close(limiter.abortCh)
			return
		}
	}
}

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	limiter := &Limiter{
		newRequestCh: make(chan context.Context),
		maxCount:     maxCount,
		interval:     interval,
		stopped:      false,
		abortCh:      make(chan struct{}),
	}

	if maxCount > 0 {
		go StartRequestMonitor(limiter)
	}

	return limiter
}

func (l *Limiter) Acquire(ctx context.Context) error {
	if l.stopped {
		return ErrStopped
	}

	select {
	case l.newRequestCh <- ctx:
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

func (l *Limiter) Stop() {
	if !l.stopped {
		l.stopped = true
		l.abortCh <- struct{}{}
	}
}
