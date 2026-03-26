package main

import (
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
}

type Limiter struct {
	clock      Clock
	ratePerSec float64
	burst      int
	tokens     float64
	mu         sync.Mutex
	lastUpdate time.Time
}

func NewLimiter(clock Clock, ratePerSec float64, burst int) *Limiter {
	return &Limiter{
		clock:      clock,
		ratePerSec: ratePerSec,
		burst:      burst,
		tokens:     float64(burst),
		lastUpdate: clock.Now(),
	}
}

func (l *Limiter) Allow() bool {
	if l.burst <= 0 {
		return false
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	currentTime := l.clock.Now()
	passedTime := currentTime.Sub(l.lastUpdate)
	l.tokens = min(float64(l.burst), l.tokens+l.ratePerSec*passedTime.Seconds())

	l.lastUpdate = currentTime
	if l.tokens < 1 {
		return false
	}
	l.tokens -= 1
	return true
}
