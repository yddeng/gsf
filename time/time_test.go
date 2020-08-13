package time

import (
	"fmt"
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	nn1 := time.Now()
	n1 := Now()
	fmt.Println(nn1.String(), n1.String(), nn1.Sub(n1).String())

	t1, err := ParseTime("2020-8-14 10:40:00")
	if err != nil {
		fmt.Println(err)
		return
	}
	Init(t1)
	time.Sleep(time.Second)
	n2 := Now()
	nn2 := time.Now()
	fmt.Println(nn2.String(), n2.String(), nn2.Sub(n2).String())

}

func TestCalcLatestTimeAfter(t *testing.T) {
	t1, err := ParseTime("2020-8-13 10:40:00")
	if err != nil {
		fmt.Println(err)
		return
	}
	Init(t1)
	fmt.Println(Now().String())

	tt := CalcLatestTimeAfter(10, 30, 00)
	fmt.Println(tt.String())

	tt1 := CalcLatestTimeAfter(10, 50, 00)
	fmt.Println(tt1.String())
}

func TestCalcLatestWeekTimeAfter(t *testing.T) {
	t1, err := ParseTime("2020-8-13 10:40:00")
	if err != nil {
		fmt.Println(err)
		return
	}
	Init(t1)
	fmt.Println(Now().String())

	tt1 := CalcLatestWeekTimeAfter(time.Monday, 10, 00, 00)
	fmt.Println(tt1.String())

	tt2 := CalcLatestWeekTimeAfter(time.Thursday, 10, 30, 00)
	fmt.Println(tt2.String())

	tt3 := CalcLatestWeekTimeAfter(time.Thursday, 10, 50, 00)
	fmt.Println(tt3.String())

	tt4 := CalcLatestWeekTimeAfter(time.Friday, 10, 00, 00)
	fmt.Println(tt4.String())
}

func TestNone(t *testing.T) {
	now := time.Now()
	fmt.Println(now.UnixNano(), now.Unix())
}
