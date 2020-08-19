package sbtree

import (
	"sync"
)

type Ranking struct {
	lock     sync.RWMutex
	tree     *tree
	scoreMap map[int64]int64
}

func NewRanking() *Ranking {
	return &Ranking{scoreMap: map[int64]int64{}}
}

// Set 设置 key 对应的分数
// 排名由大到小，得分相同情况下，最新更新的 key 排名更靠前
func (r *Ranking) Set(key, score int64) {
	r.lock.Lock()
	defer r.lock.Unlock()
	if oldScore, ok := r.scoreMap[key]; ok {
		if oldScore == score {
			return
		}
		r.tree = del(r.tree, key, oldScore)
		r.scoreMap[key] = score
		r.tree = add(r.tree, key, score)
	} else {
		r.tree = add(r.tree, key, score)
		r.scoreMap[key] = score
	}
}

// 由大到小排序，最大元素的排序为 1
func (r *Ranking) Get(key int64) int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if score, ok := r.scoreMap[key]; ok {
		return rank(r.tree, key, score)
	}
	return 0
}

// 获取排名区间数据，返回区间内 key,score 集合，序号从 1 开始，如数据长度不足，则返回所有有效数据
// [from,to], if to<from, to=from
func (r *Ranking) GetRange(from, to int) ([]int64, []int64) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if to < from {
		to = from
	}
	return getRange(r.tree, from, to, []int64{}, []int64{})
}

// 获取排名为 n 的 key,score
// n>0 && n<=ranking.len，否则返回 "",0
func (r *Ranking) GetN(n int) (int64, int64) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	if t := getN(r.tree, n); t != nil {
		return t.key, t.score
	}
	return 0, 0
}

func (r *Ranking) Len() int {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return len(r.scoreMap)
}

// index 从 1 开始
func (r *Ranking) Walk(handler func(index int, key, score int64)) {
	r.lock.RLock()
	defer r.lock.RUnlock()
	walk(r.tree, 1, handler)
}

func (r *Ranking) Copy() *Ranking {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return &Ranking{
		tree:     copyTree(r.tree),
		scoreMap: copyMap(r.scoreMap),
	}
}

func copyMap(m map[int64]int64) map[int64]int64 {
	m2 := map[int64]int64{}
	for k, v := range m {
		m2[k] = v
	}
	return m2
}
