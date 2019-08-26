package clock

import (
	"container/heap"
	"sort"
	"sync"
	"time"
)

type timerElem struct {
	*Timer

	Index int
}

type timers struct {
	heap []*timerElem
}

func (self *timers) Len() int {
	return len(self.heap)
}
func (self *timers) Swap(i, j int) {
	self.heap[i], self.heap[j] = self.heap[j], self.heap[i]
	self.heap[i].Index = i
	self.heap[j].Index = j
}
func (self *timers) Less(i, j int) bool {
	return self.heap[i].target().Before(self.heap[j].target())
}
func (self *timers) Push(x interface{}) {
	n := self.Len()
	elem := x.(*timerElem)
	elem.Index = n
	self.heap = append(self.heap, elem)
}

func (self *timers) Pop() interface{} {
	old := self.heap
	n := len(old)
	elem := old[n-1]
	old[n-1] = nil
	self.heap = old[0 : n-1]

	elem.Index = -1
	return elem
}
func (self *timers) add(elem *timerElem) {
	if elem.Index != -1 {
		return
	}
	heap.Push(self, elem)
}
func (self *timers) remove(elem *timerElem) {
	if elem.Index == -1 {
		return
	}
	heap.Remove(self, elem.Index)
}
func (self *timers) update(elem *timerElem) {
	if elem.Index == -1 {
		panic("updating non-existent timerElem")
	}
	heap.Fix(self, elem.Index)
}

type tickerElem struct {
	*Ticker

	Index int
}

type tickers struct {
	elems []*tickerElem
}

func (self *tickers) add(elem *tickerElem) {
	if elem.Index != -1 {
		return
	}

	index := sort.Search(len(self.elems), func(i int) bool {
		return self.elems[i].target().After(elem.target())
	})
	self.elems = append(self.elems, &tickerElem{})
	copy(self.elems[index+1:], self.elems[index:])
	self.elems[index] = elem
	elem.Index = index
}
func (self *tickers) remove(elem *tickerElem) {
	if elem.Index == -1 {
		return
	}

	i := elem.Index
	elem.Index = -1
	copy(self.elems[i:], self.elems[i+1:])
	self.elems = self.elems[:len(self.elems)-1]
}

type Mock struct {
	mu      sync.RWMutex
	now     time.Time
	timers  timers
	tickers tickers
}

func NewMock() *Mock {
	return &Mock{}
}

func (self *Mock) process(t time.Time) {
	var more bool
	for {
		more = self.fireNextTimer(t)
		if !more {
			break
		}
	}

	self.tick(t)
}
func (self *Mock) newTimer(d time.Duration, f func()) *Timer {
	elem := &timerElem{
		Index: -1,
	}
	add := func() {
		self.mu.Lock()
		if elem.Index == -1 {
			self.timers.add(elem)
		} else {
			self.timers.update(elem)
		}
		self.mu.Unlock()
	}
	remove := func() {
		self.mu.Lock()
		self.timers.remove(elem)
		self.mu.Unlock()
	}
	self.mu.Lock()
	elem.Timer = newMockTimer(self.now.Add(d), f, self.Now, add, remove)
	self.timers.add(elem)
	self.mu.Unlock()

	return elem.Timer
}

func (self *Mock) fireNextTimer(target time.Time) bool {
	self.mu.RLock()
	if self.timers.Len() == 0 {
		self.mu.RUnlock()
		return false
	}

	if self.timers.heap[0].target().After(target) {
		self.mu.RUnlock()

		return false
	}
	self.mu.RUnlock()

	self.mu.Lock()
	if self.timers.Len() == 0 {
		self.mu.Unlock()
		return false
	}

	if self.timers.heap[0].target().After(target) {
		self.mu.Unlock()
		return false
	}

	elem := heap.Pop(&(self.timers)).(*timerElem)
	timer := elem.Timer
	now := self.now
	self.mu.Unlock()

	timer.fire(now)
	return true
}

func (self *Mock) tick(target time.Time) {
	self.mu.RLock()

	for _, ticker := range self.tickers.elems {
		if ticker.target().After(target) {
			break
		}
		ticker.tick(target)
	}

	self.mu.RUnlock()
}

func (self *Mock) Advance(d time.Duration) {
	var t time.Time

	self.mu.Lock()
	t = self.now.Add(d)
	self.now = t
	self.mu.Unlock()

	self.process(t)
}
func (self *Mock) Set(t time.Time) {
	self.mu.Lock()
	self.now = t
	self.mu.Unlock()

	self.process(t)
}

func (self *Mock) After(d time.Duration) <-chan time.Time {
	return self.NewTimer(d).C
}
func (self *Mock) AfterFunc(d time.Duration, f func()) *Timer {
	return self.newTimer(d, f)
}
func (self *Mock) NewTimer(d time.Duration) *Timer {
	return self.newTimer(d, nil)
}
func (self *Mock) Now() time.Time {
	self.mu.RLock()
	defer self.mu.RUnlock()

	return self.now
}
func (self *Mock) Since(t time.Time) time.Duration {
	return self.Now().Sub(t)
}
func (self *Mock) Sleep(d time.Duration) {
	<-self.After(d)
}
func (self *Mock) Tick(d time.Duration) <-chan time.Time {
	if d <= 0 {
		return nil
	}
	return self.NewTicker(d).C
}

func (self *Mock) NewTicker(d time.Duration) *Ticker {
	elem := &tickerElem{
		Index: -1,
	}
	remove := func() {
		self.mu.Lock()
		self.tickers.remove(elem)
		self.mu.Unlock()
	}
	now := self.Now()

	self.mu.Lock()
	elem.Ticker = newMockTicker(d, now, remove)
	self.tickers.add(elem)
	self.mu.Unlock()

	return elem.Ticker
}
