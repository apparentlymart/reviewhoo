package main

import (
	"sort"
)

type StringScore struct {
	String string
	Score  int
}

type StringRanking []*StringScore

func NewRanking(tally map[string]int) StringRanking {
	ret := make(StringRanking, 0, len(tally))
	for str, score := range tally {
		ret.Add(str, score)
	}
	ret.Sort()
	return ret
}

func (r *StringRanking) Add(s string, score int) {
	*r = append(*r, &StringScore{
		String: s,
		Score:  score,
	})
}

func (r StringRanking) Sort() {
	sort.Stable(r)
}

func (r StringRanking) Len() int {
	return len(r)
}

func (r StringRanking) Less(i, j int) bool {
	// We want a "highest score first" sort, so this
	// comparison is intentionally backwards.
	return r[i].Score > r[j].Score
}

func (r StringRanking) Swap(i, j int) {
	temp := r[j]
	r[j] = r[i]
	r[i] = temp
}

func (r StringRanking) Highest() int {
	if len(r) == 0 {
		return 0
	}
	return r[0].Score
}

func (r StringRanking) TopN(n int) StringSet {
	if n > len(r) {
		n = len(r)
	}

	s := make(StringSet)
	for i := 0; i < n; i++ {
		s.Add(r[i].String)
	}
	return s
}

// TopPercent returns a set of the authors in the ranking that have scores
// within p percent of the highest score.
//
// minScore is an absolute lower limit to avoid returning nonsense for
// rankings that have only a few very low scores.
//
// This assumes that the ranking has already been sorted.
func (r StringRanking) TopPercent(p int, minScore int) StringSet {
	if len(r) == 0 {
		return make(StringSet)
	}

	topScore := r[0].Score
	limitScore := (topScore * p) / 100

	s := make(StringSet)
	for _, score := range r {
		if score.Score < limitScore {
			break
		}
		if score.Score < minScore {
			break
		}
		s.Add(score.String)
	}
	return s
}
