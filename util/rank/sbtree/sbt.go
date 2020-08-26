package sbtree

import "fmt"

type tree struct {
	left  *tree
	right *tree
	size  int
	key   int64
	score int64
}

func add(t *tree, key, score int64) *tree {
	if t == nil {
		return &tree{key: key, score: score, size: 1}
	} else if score < t.score {
		t.left = add(t.left, key, score)
	} else {
		t.right = add(t.right, key, score)
	}
	t.size++
	t = maintain(t, score > t.score)
	return t
}

func del(t *tree, key, score int64) *tree {
	if t == nil {
		return nil
	} else if t.key == key {
		if t.left == nil {
			return t.right
		} else if t.right == nil {
			return t.left
		} else {
			first := getFirst(t.right)
			t.key, t.score, first.key, first.score = first.key, first.score, t.key, t.score
			t.right = del(t.right, key, score)
		}
	} else if score <= t.score { // 注意此处是 <= 而不是 <，由于在插入时使用的是 < 导致相同 score 的节点在左边，此处需要使用 <= 来检查左方数据
		t.left = del(t.left, key, score)
		t = maintain(t, true)
	} else {
		t.right = del(t.right, key, score)
		t = maintain(t, false)
	}
	t.size--
	return t
}

// 由小到大，从 1 开始，0 表示无排序
func rank(t *tree, key, score int64) int {
	if t == nil {
		return 0
	}
	if key == t.key {
		return size(t.right) + 1
	} else if score >= t.score {
		return rank(t.right, key, score)
	} else {
		return size(t.right) + rank(t.left, key, score) + 1
	}
}

func maintain(t *tree, flag bool) *tree {
	if !flag {
		if t.left != nil && size(t.left.left) > size(t.right) {
			t = rotateRight(t)
		} else if t.left != nil && size(t.left.right) > size(t.right) {
			t.left = rotateLeft(t.left)
			t = rotateRight(t)
		} else {
			return t
		}
	} else {
		if t.right != nil && size(t.right.right) > size(t.left) {
			t = rotateLeft(t)
		} else if t.right != nil && size(t.right.left) > size(t.left) {
			t.right = rotateRight(t.right)
			t = rotateLeft(t)
		} else {
			return t
		}
	}
	t.left = maintain(t.left, false)
	t.right = maintain(t.right, true)
	t = maintain(t, false)
	t = maintain(t, true)
	return t
}

func rotateRight(t *tree) *tree {
	if left := t.left; left != nil {
		t.left = left.right
		left.right = t
		left.size = t.size
		t.size = size(t.left) + size(t.right) + 1
		return left
	} else {
		return nil
	}
}

func rotateLeft(t *tree) *tree {
	if right := t.right; right != nil {
		t.right = right.left
		right.left = t
		right.size = t.size
		t.size = size(t.left) + size(t.right) + 1
		return right
	} else {
		return nil
	}
}

func getFirst(t *tree) *tree {
	for t.left != nil {
		t = t.left
	}
	return t
}

func check(t *tree) bool {
	if t.right != nil && (size(t.right.left) > size(t.left) || size(t.right.right) > size(t.left)) {
		return false
	} else if t.left != nil && (size(t.left.left) > size(t.right) || size(t.left.right) > size(t.right)) {
		return false
	} else if t.left != nil && !check(t.left) {
		return false
	} else if t.right != nil && !check(t.right) {
		return false
	} else {
		return true
	}
}

func size(t *tree) int {
	if t == nil {
		return 0
	}
	return t.size
}

func copyTree(t *tree) *tree {
	if t == nil {
		return nil
	}
	return &tree{
		left:  copyTree(t.left),
		right: copyTree(t.right),
		size:  t.size,
		key:   t.key,
		score: t.score,
	}
}

func walk(t *tree, index int, handler func(index int, key, score int64)) int {
	if t == nil {
		return index
	}
	index = walk(t.right, index, handler)
	handler(index, t.key, t.score)
	index++
	return walk(t.left, index, handler)
}

func getN(t *tree, n int) *tree {
	if t == nil || n <= 0 || n > t.size {
		return nil
	}
	if tIndex := size(t.right) + 1; n == tIndex {
		return t
	} else if n > tIndex {
		return getN(t.left, n-tIndex)
	} else {
		return getN(t.right, n)
	}
}

// from<=to, return [from,to]
func getRange(t *tree, from, to int, keys, scores []int64) ([]int64, []int64) {
	if t != nil && from <= to && from <= t.size {
		if tIndex := size(t.right) + 1; from > tIndex {
			keys, scores = getRange(t.left, from-tIndex, to-tIndex, keys, scores)
		} else if from < tIndex {
			keys, scores = getRange(t.right, from, to, keys, scores)
			if to >= tIndex {
				keys, scores = append(keys, t.key), append(scores, t.score)
				keys, scores = getRange(t.left, 1, to-tIndex, keys, scores)
			}
		} else {
			keys, scores = append(keys, t.key), append(scores, t.score)
			if to > tIndex {
				keys, scores = getRange(t.left, 1, to-tIndex, keys, scores)
			}
		}
	}
	return keys, scores
}

func (t *tree) print() {
	t.debug("")
}

func (t *tree) debug(pre string) {
	if t == nil {
		return
	}
	fmt.Println(pre, t.key, t.score)
	pre += "-"
	if t.left != nil {
		fmt.Println("left:")
		t.left.debug(pre)
	}
	if t.right != nil {
		fmt.Println("right:")
		t.right.debug(pre)
	}
}
