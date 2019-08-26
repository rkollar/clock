package clock

import (
	"sync"
	"time"
)

type Ticker struct {
	C <-chan time.Time

	// real
	ticker *time.Ticker

	// mock
	c        chan<- time.Time
	interval time.Duration
	remove   func()
	mu       sync.Mutex
	stopped  bool
	target_  time.Time
}

func newMockTicker(d time.Duration, now time.Time, remove func()) *Ticker {
	c := make(chan time.Time, 1)
	return &Ticker{
		C:        c,
		c:        c,
		target_:  now.Add(d),
		interval: d,
		remove:   remove,
	}
}
func newRealTicker(d time.Duration) *Ticker {
	realTicker := time.NewTicker(d)
	return &Ticker{
		C:      realTicker.C,
		ticker: realTicker,
	}
}

func (self *Ticker) target() time.Time {
	self.mu.Lock()
	defer self.mu.Unlock()

	return self.target_
}

func (self *Ticker) tick(target time.Time) {
	var val time.Time

	self.mu.Lock()
	if self.target_.After(target) {
		self.mu.Unlock()
		return
	}

	val = self.target_
	self.target_ = val.Add(self.interval)

	select {
	case self.c <- val:
	default:
		self.mu.Unlock()
		return
	}
	self.mu.Unlock()

}

func (self *Ticker) Stop() {
	if self.ticker != nil {
		self.ticker.Stop()
		return
	}

	self.mu.Lock()
	if self.stopped {
		self.mu.Unlock()
		return
	}
	self.mu.Unlock()

	if self.remove != nil {
		self.remove()
	}
}
