package segmenttree

import (
	"fmt"
	"testing"
)

func TestSum(t *testing.T) {
	v := []int{1, 3, -2, 8, -7}
	tree := New(v, func(a, b int) int { return a + b })
	tests := []struct {
		lo, hi        int
		want, newWant int
	}{
		{
			lo:      2,
			hi:      5,
			want:    -1,
			newWant: 4,
		},
		{
			lo:      0,
			hi:      5,
			want:    3,
			newWant: 8,
		},
		{
			lo:      2,
			hi:      3,
			want:    -2,
			newWant: 3,
		},
		{
			lo:      3,
			hi:      3,
			want:    0,
			newWant: 0,
		},
		{
			lo:      10,
			hi:      20,
			want:    0,
			newWant: 0,
		},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("[%d:%d]", tc.lo, tc.hi), func(t *testing.T) {
			if got := tree.Aggregate(tc.lo, tc.hi); got != tc.want {
				t.Errorf("Aggregate(%d, %d) = %d, want %d", tc.lo, tc.hi, got, tc.want)
			}
		})
	}
	// Set v[2] = 3.
	tree.Update(2, 3)
	for _, tc := range tests {
		t.Run(fmt.Sprintf("[%d:%d]", tc.lo, tc.hi), func(t *testing.T) {
			if got := tree.Aggregate(tc.lo, tc.hi); got != tc.newWant {
				t.Errorf("Aggregate(%d, %d) = %d, want %d", tc.lo, tc.hi, got, tc.newWant)
			}
		})
	}
}

func TestMaxTimes(t *testing.T) {
	v := []int{1, 3, -2, 3, 1, 8, -7, 4, 11, 2, 0, 2, 11, 11}
	tree := New(newPairs(v), combinePairs)
	tests := []struct {
		lo, hi        int
		want, newWant [2]int
	}{
		{
			lo:      2,
			hi:      5,
			want:    [2]int{3, 1},
			newWant: [2]int{3, 2},
		},
		{
			lo:      0,
			hi:      5,
			want:    [2]int{3, 2},
			newWant: [2]int{3, 3},
		},
		{
			lo:      2,
			hi:      3,
			want:    [2]int{-2, 1},
			newWant: [2]int{3, 1},
		},
		{
			lo:      3,
			hi:      3,
			want:    [2]int{0, 0},
			newWant: [2]int{0, 0},
		},
		{
			lo:      10,
			hi:      20,
			want:    [2]int{11, 2},
			newWant: [2]int{11, 2},
		},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("[%d:%d]", tc.lo, tc.hi), func(t *testing.T) {
			if got := tree.Aggregate(tc.lo, tc.hi); got != tc.want {
				t.Errorf("Aggregate(%d, %d) = %d, want %d", tc.lo, tc.hi, got, tc.want)
			}
		})
	}
	// Set v[2] = 3.
	tree.Update(2, [2]int{3, 1})
	for _, tc := range tests {
		t.Run(fmt.Sprintf("[%d:%d]", tc.lo, tc.hi), func(t *testing.T) {
			if got := tree.Aggregate(tc.lo, tc.hi); got != tc.newWant {
				t.Errorf("Aggregate(%d, %d) = %d, want %d", tc.lo, tc.hi, got, tc.newWant)
			}
		})
	}
}

func TestKthPositive(t *testing.T) {
	v := []int{1, 3, -2, 3, 1, 8, -7, 4, 11, 2, 0, 2, 11, 11}
	positives, indices := countPositives(v)
	tree := New(positives, func(a, b int) int { return a + b })
	t.Logf("Tree: %v", tree)

	indices = append(indices, -1)
	for i, v := range indices {
		t.Run(fmt.Sprintf("%d positive", i+1), func(t *testing.T) {
			if got := findPositive(NewWalker(tree), i+1); got != v {
				t.Errorf("findPositive(%d) = %d, want %d", i+1, got, v)
			}
		})
	}
}

func newPairs(values []int) [][2]int {
	ret := make([][2]int, 0, len(values))
	for _, v := range values {
		ret = append(ret, [2]int{v, 1})
	}
	return ret
}

func combinePairs(a, b [2]int) [2]int {
	switch {
	case a[1] == 0:
		return b
	case b[1] == 0:
		return a
	case a[0] > b[0]:
		return a
	case a[0] < b[0]:
		return b
	default:
		return [2]int{a[0], a[1] + b[1]}
	}
}

func countPositives(values []int) ([]int, []int) {
	ret := make([]int, 0, len(values))
	var idxs []int
	for i, v := range values {
		if v > 0 {
			ret = append(ret, 1)
			idxs = append(idxs, i)
		} else {
			ret = append(ret, 0)
		}
	}
	return ret, idxs
}

func findPositive(walker *Walker[int], k int) int {
	if k > walker.RootValue() {
		return -1
	}
	if walker.Right()-walker.Left() == 1 {
		return walker.Left()
	}
	left := walker.LeftValue()
	if left >= k {
		walker.MoveLeft()
		return findPositive(walker, k)
	}
	walker.MoveRight()
	return findPositive(walker, k-left)
}
