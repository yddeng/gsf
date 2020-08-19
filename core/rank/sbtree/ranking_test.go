package sbtree

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestRanking(t *testing.T) {
	count := 100 * 10000
	r := NewRanking()
	nums := createNums(count)
	startTime := time.Now()
	for k, v := range nums {
		r.Set(k, v) // 乱序插入
	}
	fmt.Printf("write useTime:%.2fs\n", time.Now().Sub(startTime).Seconds())

	startTime = time.Now()
	for k, v := range nums {
		r.Set(k, v+1) // 乱序更新
	}
	fmt.Printf("reset useTime:%.2fs\n", time.Now().Sub(startTime).Seconds())

	startTime = time.Now()
	for k := range nums {
		should := count - int(k)
		if rank := r.Get(k); rank != should {
			t.Errorf("ranking fail %d should %d, but %d", k, should, rank)
		}
	}
	fmt.Printf("read useTime:%.2fs\n", time.Now().Sub(startTime).Seconds())
}

func TestRanking_Set(t *testing.T) {
	r := NewRanking()
	r.Set(0, 0)
	r.Set(1, 1)
	r.Set(2, 2)

	// reset, 得分相同的情况下，最新更新的 key 排名更靠前
	r.Set(1, 2)
	if r.Get(1) != 1 {
		t.Errorf("ranking set fail")
	}
}

func TestRanking_Copy_Walk(t *testing.T) {
	count := 10
	r := createRanking(count)

	r2 := r.Copy()
	r2.Walk(func(index int, key, score int64) {
		if key != score {
			t.Errorf("ranking copy fail, key:%d score:%d", key, score)
		} else if index != count-int(score) {
			t.Errorf("ranking copy fail, key:%d index:%d", key, index)
		}
	})
}

func TestRanking_GetN(t *testing.T) {
	count := 10
	r := createRanking(count)

	for i := 1; i < count+1; i++ {
		k, v := int64(count-i), int64(count-i)
		if key, score := r.GetN(i); key != k && score != v {
			t.Errorf("ranking getn fail, should be %d,%d, but %d,%d", k, v, key, score)
		}
	}
}

func TestRanking_GetRange(t *testing.T) {
	count := 10
	r := createRanking(count)

	ranges := []struct {
		from int
		to   int
	}{
		{1, 10},
		{1, 3},
		{2, 5},
		{2, 8},
		{5, 8},
		{6, 9},
		{1, 15},
		{-1, 8},
		{-1, 15},
		{5, 5},
	}

	for _, rang := range ranges {
		from, to := rang.from, rang.to
		keys, scores := r.GetRange(from, to)
		// 获取实际长度 //
		if to > r.Len() {
			to = r.Len()
		}
		if from < 1 {
			from = 1
		}
		if to < from {
			to = from
		}
		size := to - from + 1
		// /////////// //
		if len(keys) != size {
			t.Errorf("ranking get range fail, no enough data, should %d, but %d", size, len(keys))
		}
		for i, v := range keys {
			key, score := int64(count-from), int64(count-from)
			if v != key || scores[i] != score {
				t.Errorf("ranking get rank fail, should be %d,%df, but %d,%d", key, score, v, scores[i])
			}
			from++
		}
	}
}

func createRanking(count int) *Ranking {
	r := NewRanking()
	nums := createNums(count)
	for k, v := range nums {
		r.Set(k, v) // 乱序插入
	}
	return r
}

// map[i]i, ex: 1:1, 2:2
func createNums(count int) map[int64]int64 {
	nums := map[int64]int64{}
	for i := 0; i < count; i++ {
		nums[int64(i)] = int64(i)
	}
	return nums
}

func k2v(k string) int {
	i, _ := strconv.Atoi(k)
	return i
}

func getName(i int) string {
	return strconv.Itoa(i)
}
