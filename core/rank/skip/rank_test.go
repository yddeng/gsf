package skip

//go test -covermode=count -v -coverprofile=coverage.out -run=. -cpuprofile=rank.p
//go tool cover -html=coverage.out
//go tool pprof rank.p
//go test -v -run=^$ -bench BenchmarkRank -count 10

import (
	"fmt"
	"github.com/schollz/progressbar"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"testing"
	"time"
)

func TestBenchmarkRank2(t *testing.T) {
	var r *Rank = NewRank()
	testCount := 50000000
	idRange := 10000000
	{
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			item := r.getRankItem(uint64(idx))
			var score int
			if nil == item {
				score = rand.Int() % 1000000
			} else {
				score = item.value + rand.Int()%10000
			}
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg), len(r.id2Item))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
	}
}

func TestBenchmarkRank1(t *testing.T) {
	var r *Rank = NewRank()
	fmt.Println("TestBenchmarkRank")

	testCount := 10000000
	idRange := 10000000

	{
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			score := rand.Int()%1000000 + 1
			//fmt.Println(i, idx, score)
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
		assert.Equal(t, true, r.Check())
	}

	{
		testCount := 10000000
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			score := rand.Int()%1000000 + 1
			//fmt.Println(idx, score)
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
		assert.Equal(t, true, r.Check())
	}

	{

		testCount := 10000000

		bar := progressbar.New(int(testCount))
		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := (rand.Int() % len(r.id2Item)) + 1
			item := r.id2Item[uint64(idx)]
			score := rand.Int()%10000 + 1
			score = item.value + score
			//fmt.Println(idx, score)
			r.UpdateScore(uint64(idx), score)
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
		assert.Equal(t, true, r.Check())
	}

	{

		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			r.GetRankPersent(uint64(idx))
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
	}

	{
		bar := progressbar.New(int(testCount))

		beg := time.Now()
		for i := 0; i < testCount; i++ {
			idx := i%idRange + 1
			r.GetRank(uint64(idx))
			bar.Add(1)
		}
		fmt.Println(time.Now().Sub(beg))
		fmt.Println(len(r.spans), len(r.id2Item)/len(r.spans))
	}
}

func TestRank(t *testing.T) {
	fmt.Println("TestRank")

	var r *Rank = NewRank()
	fmt.Println("TestBenchmarkRank")

	testCount := 200
	idRange := 1000

	for i := 0; i < testCount; i++ {
		idx := i%idRange + 1
		score := rand.Int() % 10000
		r.UpdateScore(uint64(idx), score)
	}

	r.Show()

}

func BenchmarkRank(b *testing.B) {
	var r *Rank = NewRank()
	for i := 0; i < b.N; i++ {
		idx := (i % 1000000) + 1
		score := rand.Int()
		r.UpdateScore(uint64(idx), score)
	}
}

func TestRank_GetTopN(t *testing.T) {
	var r *Rank = NewRank()
	testCount := 200
	idRange := 1000

	for i := 0; i < testCount; i++ {
		idx := i%idRange + 1
		score := rand.Int() % 10000
		r.UpdateScore(uint64(idx), score)
	}

	r.Show()

	ids := r.GetTopN(50)
	fmt.Println(ids, "--", len(ids))

	ids = r.GetTopN(220)
	fmt.Println(ids, "--", len(ids))

}

func TestRank_GetScoreByIdx(t *testing.T) {
	var r *Rank = NewRank()
	testCount := 200
	idRange := 1000

	for i := 0; i < testCount; i++ {
		idx := i%idRange + 1
		score := rand.Int() % 10000
		r.UpdateScore(uint64(idx), score)
	}

	r.Show()

	fmt.Println(r.GetScoreByIdx(5))
	fmt.Println(r.GetScoreByIdx(200))
	fmt.Println(r.GetScoreByIdx(220))
}
