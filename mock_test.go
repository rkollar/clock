package clock

import (
	"container/heap"
	"fmt"
	"testing"
	"time"
)

func ASSERT(tb testing.TB, condition bool, msg string, v ...interface{}) {
	tb.Helper()

	if !condition {
		tb.Fatalf(fmt.Sprintf(msg, v...))
	}
}

func OK(tb testing.TB, err error) {
	tb.Helper()

	if err != nil {
		tb.Fatalf("Unexpected ERROR: %s", err)
	}
}

func TestTimer(t *testing.T) {
	t.Run("order", func(t *testing.T) {
		clock := NewMock()
		clock.Set(time.Now())

		_ = clock.NewTimer(10 * time.Second)
		_ = clock.NewTimer(30 * time.Second)
		_ = clock.NewTimer(20 * time.Second)

		var last time.Time
		for len(clock.timers.heap) > 0 {
			timerElem := heap.Pop(&(clock.timers)).(*timerElem)
			ASSERT(t, !timerElem.target().Before(last), "bad heap order")
			last = timerElem.target()
		}
	})
	t.Run("AfterFunc", func(t *testing.T) {
	})
	t.Run("NewTimer", func(t *testing.T) {
		t.Run("fire", func(t *testing.T) {
			clock := NewMock()
			clock.Set(time.Now())

			start := clock.Now()
			timer := clock.NewTimer(10 * time.Second)

			clock.Advance(5 * time.Second)
			select {
			case <-timer.C:
				ASSERT(t, false, "timer fired prematurely")
			default:
			}

			clock.Advance(5 * time.Second)
			var c time.Time
			select {
			case c = <-timer.C:
			default:
				ASSERT(t, false, "timer has not fired")
			}
			ASSERT(t, c == start.Add(10*time.Second), "bad time received: %s", c)
		})
		t.Run("stop", func(t *testing.T) {
			clock := NewMock()
			clock.Set(time.Now())

			timer := clock.NewTimer(10 * time.Second)

			clock.Advance(5 * time.Second)
			stopped := timer.Stop()
			ASSERT(t, stopped, "timer was not stopped")
			ASSERT(t, clock.timers.Len() == 0, "timer not removed")

			select {
			case <-timer.C:
				ASSERT(t, false, "ticker fired after being stopped")
			default:
			}
		})
		t.Run("reset", func(t *testing.T) {
			clock := NewMock()
			clock.Set(time.Now())

			timer := clock.NewTimer(10 * time.Second)

			clock.Advance(5 * time.Second)
			reset := timer.Reset(10 * time.Second)
			ASSERT(t, reset, "timer was not reset")

			clock.Advance(5 * time.Second)
			ASSERT(t, clock.timers.Len() == 1, "timer removed")

			clock.Advance(5 * time.Second)
			stopped := timer.Stop()
			ASSERT(t, !stopped, "timer stopped")
			ASSERT(t, clock.timers.Len() == 0, "timer not removed")
		})
		t.Run("reset after stop", func(t *testing.T) {
			t.Skip() //TODO
		})
	})
}

func TestTicker(t *testing.T) {
	t.Run("tick", func(t *testing.T) {
		clock := NewMock()
		clock.Set(time.Now())

		start := clock.Now()
		ticker := clock.NewTicker(10 * time.Second)

		clock.Advance(5 * time.Second)
		select {
		case <-ticker.C:
			ASSERT(t, false, "ticker fired prematurely")
		default:
		}

		clock.Advance(5 * time.Second)
		var c time.Time
		select {
		case c = <-ticker.C:
		default:
			ASSERT(t, false, "ticker has not fired")
		}
		ASSERT(t, c == start.Add(10*time.Second), "bad time received: %s", c)
	})
	t.Run("stop", func(t *testing.T) {
		clock := NewMock()
		clock.Set(time.Now())

		ticker := clock.NewTicker(10 * time.Second)

		clock.Advance(5 * time.Second)
		ticker.Stop()
		ASSERT(t, len(clock.tickers.elems) == 0, "ticker not removed")

		select {
		case <-ticker.C:
			ASSERT(t, false, "ticker fired after being stopped")
		default:
		}
	})
}
