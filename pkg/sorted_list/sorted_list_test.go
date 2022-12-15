package sortedlist

import (
	"fmt"
	"sort"
	"testing"

	"golang.org/x/exp/slices"
)

func TestSortedList(t *testing.T) {
	s := New(3)

	items := []float64{1.2, 2.3, 67.2, 12.6, 7.2, 6.1}
	for _, v := range items {
		s.MaybeAdd(v, v)
	}

	res := s.Results()

	castedResults := make([]float64, 3)
	for idx, r := range res {
		if r != nil {
			castedResults[idx] = r.(float64)
		} else {
			castedResults[idx] = 0.0
		}
	}

	sort.Sort(sort.Float64Slice(items))

	fin := make([]float64, 3)
	subItems := items[len(items)-3:]
	for i := 2; i >= 0; i-- {
		fmt.Println(i)
		fin[i] = subItems[i]
	}

	if slices.Compare(castedResults, fin) != 0 {
		fmt.Println(castedResults, fin)
		t.Errorf("slices do not equate")
	}
}
