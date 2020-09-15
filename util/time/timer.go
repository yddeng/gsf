package time

import (
	"github.com/yddeng/dutil/heap"
	"github.com/yddeng/dutil/queue"
	"sync"
	"sync/atomic"
	"time"
)

type Timer struct {
	d      time.Duration
	et     time.Time
	f      func(t *Timer, now time.Time)
	eveQue *queue.EventQueue
	repeat bool
	stop   int32
	mgr    *TimerMgr
}

// 时间偏移后的时间戳
func whenTime(d time.Duration) time.Time {
	return Now().Add(d)
}

func newTimer(d time.Duration, f func(t *Timer, now time.Time), eveQue *queue.EventQueue, repeat bool) *Timer {
	return &Timer{et: whenTime(d), d: d, f: f, eveQue: eveQue, repeat: repeat}
}

func (t *Timer) Less(e heap.Element) bool {
	return t.et.Before(e.(*Timer).et)
}

func (t *Timer) do(now *time.Time) {
	if atomic.LoadInt32(&t.stop) == 1 {
		return
	}

	if t.eveQue == nil {
		go t.f(t, *now)
	} else {
		t.eveQue.Push(func() {
			t.f(t, *now)
		})
	}
	//repeat
	if t.repeat {
		if atomic.LoadInt32(&t.stop) == 1 {
			return
		}
		t.et = t.et.Add(t.d)
		t.mgr.addTimer(t)
	}

}

func (t *Timer) Stop() {
	atomic.StoreInt32(&t.stop, 1)
}

type TimerMgr struct {
	minHeap  *heap.Heap
	lastTime *time.Time
	*sync.Mutex
}

func (mgr *TimerMgr) addTimer(t *Timer) bool {
	mgr.Lock()
	lastTime := mgr.lastTime
	mgr.Unlock()
	// 已经过期
	if lastTime != nil && t.et.Before(*lastTime) {
		t.do(lastTime)
		return false
	}
	mgr.Lock()
	mgr.minHeap.Push(t)
	mgr.Unlock()
	t.mgr = mgr
	return true
}

func NewTimerMgr(now time.Time) *TimerMgr {
	return &TimerMgr{
		minHeap:  heap.NewHeap(),
		lastTime: &now,
		Mutex:    new(sync.Mutex),
	}
}

// 传入当前时间 单位秒
func (mgr *TimerMgr) Loop(now time.Time) {
	mgr.Lock()
	mgr.lastTime = &now
	mgr.Unlock()
	var e heap.Element
	for {
		mgr.Lock()
		e = mgr.minHeap.Peek()
		if nil != e && now.After(e.(*Timer).et) {
			t := e.(*Timer)
			mgr.minHeap.Pop()
			mgr.Unlock()
			t.do(&now)
		} else {
			mgr.Unlock()
			break
		}
	}
}

func (mgr *TimerMgr) Once(d time.Duration, f func(t *Timer, now time.Time), eventQue ...*queue.EventQueue) *Timer {
	var eveQue *queue.EventQueue
	if len(eventQue) > 0 {
		eveQue = eventQue[0]
	}
	t := newTimer(d, f, eveQue, false)
	mgr.addTimer(t)
	return t
}

func (mgr *TimerMgr) Repeat(d time.Duration, f func(t *Timer, now time.Time), eventQue ...*queue.EventQueue) *Timer {
	var eveQue *queue.EventQueue
	if len(eventQue) > 0 {
		eveQue = eventQue[0]
	}
	t := newTimer(d, f, eveQue, true)
	mgr.addTimer(t)
	return t
}

// 循环任务
func StartLoopTask(t time.Duration, f func()) *time.Ticker {
	timeTicker := time.NewTicker(t)
	go func() {
		for {
			select {
			case <-timeTicker.C:
				f()
			}
		}
	}()
	return timeTicker
}
