package time

import (
	"fmt"
	"github.com/yddeng/dutil/queue"
	"testing"
	"time"
)

func TestNewTimerMgr(t *testing.T) {
	t1, err := ParseTime("2020-08-14 10:40:00")
	if err != nil {
		fmt.Println(err)
		return
	}
	Init(t1)
	fmt.Println(Now().String())
	mgr := NewTimerMgr(t1)
	mgr.Once(time.Second, func(t *Timer, now time.Time) {
		fmt.Println("t1", now.String())
	})

	mgr.Once(time.Second*3, func(t *Timer, now time.Time) {
		fmt.Println("t2", now.String())
	})

	mgr.Repeat(time.Second, func(t *Timer, now time.Time) {
		fmt.Println("t3", now.String())
	})

	for {
		time.Sleep(time.Millisecond * 200)
		mgr.Loop(Now())
	}
}

func TestTimerMgr_Once(t *testing.T) {
	eveQue := queue.NewEventQueue(100)
	eveQue.Run(1)
	t1, err := ParseTime("2020-08-14 10:40:00")
	if err != nil {
		fmt.Println(err)
		return
	}
	Init(t1)
	fmt.Println(Now().String())
	mgr := NewTimerMgr(t1)

	StartLoopTask(time.Millisecond*200, func() {
		eveQue.Push(func() {
			mgr.Loop(Now())
		})
	})

	mgr.Once(time.Second, func(t *Timer, now time.Time) {
		fmt.Println("t1", now.String())
	}, eveQue)

	mgr.Once(time.Second*3, func(t *Timer, now time.Time) {
		fmt.Println("t2", now.String())
	}, eveQue)

	mgr.Repeat(time.Second, func(t *Timer, now time.Time) {
		fmt.Println("t3", now.String())
	}, eveQue)

	select {}
}
