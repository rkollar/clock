package clock

import (
	"time"
)

type Clock interface {
	After(d time.Duration) <-chan time.Time
	Sleep(d time.Duration)
	Tick(d time.Duration) <-chan time.Time
	Since(t time.Time) time.Duration
	NewTicker(d time.Duration) *Ticker
	Now() time.Time
	AfterFunc(d time.Duration, f func()) *Timer
	NewTimer(d time.Duration) *Timer
}
