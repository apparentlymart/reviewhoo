package main

import (
	"math/rand"
	"sort"
)

type StringSet map[string]struct{}

func (s StringSet) Add(name string) {
	s[name] = struct{}{}
}

func (s StringSet) Remove(name string) {
	delete(s, name)
}

func (s StringSet) Has(name string) bool {
	_, ok := s[name]
	return ok
}

func (s StringSet) Union(t StringSet) StringSet {
	u := make(StringSet)
	for a := range s {
		u.Add(a)
	}
	for a := range t {
		u.Add(a)
	}
	return u
}

func (s StringSet) Subtract(t StringSet) StringSet {
	r := make(StringSet)
	for a := range s {
		if !t.Has(a) {
			r.Add(a)
		}
	}
	return r
}

// Shuffled returns the members of the set as a slice where the
// ordering is random.
func (s StringSet) Shuffled() []string {
	ret := make([]string, len(s))
	perm := rand.Perm(len(s))
	pi := 0
	for a := range s {
		ret[perm[pi]] = a
		pi++
	}
	return ret
}

func (s StringSet) Sorted() []string {
	ret := make([]string, len(s))
	i := 0
	for a := range s {
		ret[i] = a
		i++
	}
	sort.Strings(ret)
	return ret
}
