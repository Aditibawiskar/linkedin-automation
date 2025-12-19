package stealth

import "time"

// ActionLimiter enforces minimum delay between actions
type ActionLimiter struct {
	last time.Time
	wait time.Duration
}

func NewLimiter(wait time.Duration) *ActionLimiter {
	return &ActionLimiter{wait: wait}
}

func (l *ActionLimiter) Wait() {
	if !l.last.IsZero() {
		time.Sleep(l.wait)
	}
	l.last = time.Now()
}
