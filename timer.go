package clock

import (
	"sync"
	"time"
)

type Timer struct {
	C <-chan time.Time

	// real
	timer *time.Timer

	// mock
	f       func()           // AfterFunc function
	now     func() time.Time // fuction to get the current time
	c       chan<- time.Time // send-only channel to access the receive-only channel C
	add     func()           // callback when timer becomes active
	remove  func()           // callback when timer becomes inactive
	mu      sync.Mutex
	target_ time.Time // target time for firing
	stopped bool      // true if the timer has fired or has been stopped
}

func newMockTimer(target time.Time, f func(), now func() time.Time, add func(), remove func()) *Timer {
	if now == nil {
		panic("nil now")
	}

	timer := &Timer{
		target_: target,
		f:       f,
		now:     now,
		add:     add,
		remove:  remove,
	}
	if f == nil {
		c := make(chan time.Time, 1)
		timer.C = c
		timer.c = c
	}

	return timer
}
func newRealTimer(d time.Duration, f func()) *Timer {
	var realTimer *time.Timer
	if f != nil {
		realTimer = time.AfterFunc(d, f)
	} else {
		realTimer = time.NewTimer(d)
	}
	return &Timer{
		C:     realTimer.C,
		timer: realTimer,
	}
}

func (self *Timer) target() time.Time {
	self.mu.Lock()
	defer self.mu.Unlock()

	return self.target_
}
func (self *Timer) stop() bool {
	self.mu.Lock()
	if self.stopped {
		self.mu.Unlock()
		return false
	}
	self.stopped = true
	self.mu.Unlock()

	if self.remove != nil {
		self.remove()
	}

	return true
}
func (self *Timer) fire(t time.Time) bool {
	if !self.stop() {
		return false
	}

	if self.f != nil {
		self.f()
	} else {
		select {
		case self.c <- t:
		default:
		}
	}

	return true
}

func (self *Timer) Stop() bool {
	if self.timer != nil {
		return self.timer.Stop()
	}

	if !self.stop() {
		return false
	}

	return true
}
func (self *Timer) Reset(d time.Duration) bool {
	if self.timer != nil {
		return self.timer.Reset(d)
	}

	var active bool

	target := self.now().Add(d)

	self.mu.Lock()
	if !self.stopped {
		active = true
	}
	self.target_ = target
	self.mu.Unlock()

	if self.add != nil {
		self.add()
	}

	return active
}
