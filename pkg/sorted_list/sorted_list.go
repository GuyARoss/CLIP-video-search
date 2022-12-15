package sortedlist

import (
	"math"
	"sync"
)

type SortableList interface {
	Results() []interface{}
	MaybeAdd(sortKey float64, item interface{})
	TotalEvaluated() int
}

type SortedList struct {
	lowestRated float64
	lowestIndex int
	evaluated   int

	items    []interface{}
	sortKeys []float64
	m        *sync.Mutex
}

func (s *SortedList) Results() []interface{} {
	return s.items
}

func (s *SortedList) addAtIndex(index int, sortKey float64, item interface{}) {
	s.sortKeys[index] = sortKey
	s.items[index] = item

	if index == 0 {
		s.lowestRated = sortKey
	}
}

func (s *SortedList) MaybeAdd(sortKey float64, item interface{}) {
	s.m.Lock()
	s.evaluated += 1

	if sortKey > s.lowestRated {
		for i := 0; i < len(s.sortKeys); i++ {
			k := s.sortKeys[i]

			if sortKey > k {
				if i > 0 {
					s.addAtIndex(i-1, k, s.items[i])
				}
				if i == len(s.sortKeys)-1 {
					s.addAtIndex(i, sortKey, item)
				}
			}

			if sortKey < k {
				s.addAtIndex(i-1, sortKey, item)

				break
			}
		}
	}
	s.m.Unlock()
}

func (s *SortedList) TotalEvaluated() int {
	return s.evaluated
}

func New(topN int) SortableList {
	return &SortedList{
		items:       make([]interface{}, topN),
		sortKeys:    make([]float64, topN),
		lowestRated: math.Inf(-1),
		m:           &sync.Mutex{},
	}
}
