package util

import (
	"github.com/yddeng/dutil/queue"
	"sync"
	"sync/atomic"
	"time"
)

type TimeTaskFunc func()

func LoopTask(t time.Duration, fun TimeTaskFunc) *time.Ticker {
	timeTicker := time.NewTicker(t)
	go func() {
		for {
			select {
			case <-timeTicker.C:
				fun()
			}
		}
	}()
	return timeTicker
}

func WaitCondition(fn func() bool, eventq ...*queue.EventQueue) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	var eventQueue *queue.EventQueue
	if len(eventq) > 0 {
		eventQueue = eventq[0]
	}

	donefire := int32(0)

	if nil == eventQueue {
		go func() {
			for {
				time.Sleep(time.Millisecond * 100)
				if fn() {
					if atomic.LoadInt32(&donefire) == 0 {
						atomic.StoreInt32(&donefire, 1)
						wg.Done()
					}
					break
				}
			}
		}()
	} else {
		go func() {
			stoped := int32(0)
			for atomic.LoadInt32(&stoped) == 0 {
				time.Sleep(time.Millisecond * 100)
				eventQueue.Push(func() {
					if fn() {
						if atomic.LoadInt32(&donefire) == 0 {
							atomic.StoreInt32(&donefire, 1)
							wg.Done()
						}
						atomic.StoreInt32(&stoped, 1)
					}
				})
			}
		}()
	}

	wg.Wait()
}

func Must(i interface{}, e error) interface{} {
	if e != nil {
		panic(e)
	}
	return i
}
