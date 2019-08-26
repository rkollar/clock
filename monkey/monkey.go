package monkey

import (
	"time"

	"bou.ke/monkey"
)

func Panic() {
	monkey.Patch(time.Now, func() time.Time {
		panic("real time call of Now")
	})
	monkey.Patch(time.AfterFunc, func(d time.Duration, f func()) *time.Timer {
		panic("real time call of AfterFunc")
	})
	monkey.Patch(time.NewTicker, func(d time.Duration) *time.Ticker {
		panic("real time call of NewTicker")
	})
	monkey.Patch(time.NewTimer, func(d time.Duration) *time.Timer {
		panic("real time call of NewTimer")
	})
	monkey.Patch(time.Since, func(t time.Time) time.Duration {
		panic("real time call of Since")
	})
	monkey.Patch(time.Sleep, func(d time.Duration) {
		panic("real time call of Sleep")
	})
	monkey.Patch(time.After, func(d time.Duration) <-chan time.Time {
		panic("real time call of After")
	})
}

func Replace(mock clock.Mock) {
	monkey.Patch(time.Now, func() time.Time {
		return mock.Now()
	})
	monkey.Patch(time.AfterFunc, func(d time.Duration, f func()) *time.Timer {
		return mock.AfterFunc(d, f)
	})
	monkey.Patch(time.NewTicker, func(d time.Duration) *time.Ticker {
		return mock.NewTicker(d)
	})
	monkey.Patch(time.NewTimer, func(d time.Duration) *time.Timer {
		return mock.NewTimer(d)
	})
	monkey.Patch(time.Since, func(t time.Time) time.Duration {
		return mock.Since(t)
	})
	monkey.Patch(time.Sleep, func(d time.Duration) {
		mock.Sleep(d)
	})
	monkey.Patch(time.After, func(d time.Duration) <-chan time.Time {
		return mock.After(d)
	})
}
