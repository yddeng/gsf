package util

import (
	"github.com/yddeng/dutil/task"
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

func WaitCondition(fn func() bool, taskq ...*task.TaskQueue) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	var taskQueue *task.TaskQueue
	if len(taskq) > 0 {
		taskQueue = taskq[0]
	}

	donefire := int32(0)

	if nil == taskQueue {
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
				taskQueue.Push(func() {
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
