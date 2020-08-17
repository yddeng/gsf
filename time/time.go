package time

import (
	"errors"
	"fmt"
	"time"
)

var (
	startTime = time.Now()
	offsetDur = time.Duration(0)
)

// 将当前时间重置到 startTime 。
func Init(startTime_ time.Time) {
	startTime = startTime_
	offsetDur = startTime.Sub(time.Now())
}

func SetOffset(offset time.Duration) {
	offsetDur = offset
}

func StartTime() time.Time {
	return startTime
}

// 通过该接口调用，则获取时间为偏移后的时间。
func Now() time.Time {
	now := time.Now()
	if offsetDur != 0 {
		now = now.Add(offsetDur)
	}
	return now
}

// 解析"yyyy-mm-dd hh:mm:ss"格式的日期字符串
func ParseTime(s string) (time.Time, error) {
	var year, month, day, hour, minute, sec int
	n, err := fmt.Sscanf(s, "%d-%d-%d %d:%d:%d", &year, &month, &day, &hour, &minute, &sec)
	if n < 5 || err != nil ||
		month < 1 || month > 12 ||
		day < 1 || day > 31 ||
		hour < 0 || hour > 23 ||
		minute < 0 || minute > 59 ||
		sec < 0 || sec > 59 {
		return time.Time{}, errors.New("invalid time format")
	}
	return time.Date(year, time.Month(month), day, hour, minute, sec, 0, time.Local), nil

	//time.Parse("2006-01-02 15:04:05", str)
}

func getBaseTime(bases ...time.Time) time.Time {
	if len(bases) > 0 {
		return bases[0]
	}
	return Now()
}

// 日刷新时间，返回最近一个时间
func CalcLatestTimeAfter(hour, min, sec int, base ...time.Time) time.Time {
	now := getBaseTime(base...)

	t := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, 0, time.Local)
	if now.Before(t) {
		return t
	}

	return t.AddDate(0, 0, 1)
}

// 周刷新时间，返回最近一个时间。 weekday (0-6)
func CalcLatestWeekTimeAfter(weekday time.Weekday, hour, min, sec int, base ...time.Time) time.Time {
	now := getBaseTime(base...)

	day := now.Day()
	offset := weekday - now.Weekday()
	if offset != 0 {
		day += int(offset)
	}

	rt := time.Date(now.Year(), now.Month(), day, hour, min, sec, 0, time.Local)
	if now.Before(rt) {
		return rt
	}

	return rt.AddDate(0, 0, 7)
}
