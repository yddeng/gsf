package util

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInterval(b1, b2 int32) int32 {
	if b1 == b2 {
		return b1
	}

	min, max := b1, b2
	if min > max {
		min, max = max, min
	}
	return rand.Int31n(max-min+1) + min
}

func RandIntervalN(b1, b2 int32, n int) []int32 {
	if b1 == b2 {
		return []int32{b1}
	}

	min, max := b1, b2
	if min > max {
		min, max = max, min
	}
	l := max - min + 1
	if int32(n) > l {
		n = int(l)
	}

	r := make([]int32, n)
	m := make(map[int32]int32)
	for i := 0; i < n; i++ {
		v := rand.Int31n(l) + min

		if mv, ok := m[v]; ok {
			r[i] = mv
		} else {
			r[i] = v
		}

		lv := l - 1 + min
		if v != lv {
			if mv, ok := m[lv]; ok {
				m[v] = mv
			} else {
				m[v] = lv
			}
		}

		l--
	}

	return r
}
