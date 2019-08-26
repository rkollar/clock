package clock

import (
	"time"
)

type Real struct{}

func NewReal() *Real {
	return &Real{}
}

func (self *Real) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (self *Real) AfterFunc(d time.Duration, f func()) *Timer {
	return newRealTimer(d, f)
}

func (self *Real) Now() time.Time {
	return time.Now()
}

func (self *Real) Since(t time.Time) time.Duration {
	return time.Since(t)
}

func (self *Real) Sleep(d time.Duration) {
	time.Sleep(d)
}

func (self *Real) Tick(d time.Duration) <-chan time.Time {
	return time.Tick(d)
}

func (self *Real) NewTicker(d time.Duration) *Ticker {
	return newRealTicker(d)
}

func (self *Real) NewTimer(d time.Duration) *Timer {
	return newRealTimer(d, nil)
}
