// Package segmenttree implements a segment tree.
package segmenttree

// Tree represents a generic segment tree.
type Tree[T any] struct {
	nodes   []T
	length  int
	combine func(a, b T) T
}

// New constructs a segment tree from the provided slice, using the provided combine function.
func New[T any](values []T, combine func(a, b T) T) *Tree[T] {
	nodes := make([]T, 2*len(values))
	build(values, nodes, combine, 0, 0, len(values))
	return &Tree[T]{
		nodes:   nodes,
		length:  len(values),
		combine: combine,
	}
}

// Aggregate aggregates the values across the segment [lo:hi].
func (t *Tree[T]) Aggregate(lo, hi int) (ret T) {
	if lo > t.length || hi <= 0 {
		return
	}
	if lo < 0 {
		lo = 0
	}
	if hi > t.length {
		hi = t.length
	}
	return t.aggregate(0, 0, t.length, lo, hi)
}

// Update records an update corresponding to values[idx] = value.
func (t *Tree[T]) Update(idx int, value T) {
	t.update(0, 0, t.length, idx, value)
}

func build[T any](values []T, nodes []T, combine func(a, b T) T, index, left, right int) {
	if left >= right {
		return
	}
	if right-left == 1 {
		nodes[index] = values[left]
		return
	}
	lc, mid, rc := split(index, left, right)
	build(values, nodes, combine, lc, left, mid)
	build(values, nodes, combine, rc, mid, right)
	nodes[index] = combine(nodes[lc], nodes[rc])
}

func (t *Tree[T]) aggregate(index, left, right, lo, hi int) (ret T) {
	if lo >= hi {
		return
	}
	if lo == left && hi == right {
		return t.nodes[index]
	}
	lc, mid, rc := split(index, left, right)
	minimum := min(hi, mid)
	maximum := max(lo, mid)
	return t.combine(t.aggregate(lc, left, mid, lo, minimum), t.aggregate(rc, mid, right, maximum, hi))
}

func (t *Tree[T]) update(index, left, right, idx int, value T) {
	if left >= right {
		return
	}
	if right-left == 1 {
		t.nodes[index] = value
		return
	}
	lc, mid, rc := split(index, left, right)
	if idx < mid {
		t.update(lc, left, mid, idx, value)
	} else {
		t.update(rc, mid, right, idx, value)
	}
	t.nodes[index] = t.combine(t.nodes[lc], t.nodes[rc])
}

func split(index, left, right int) (lc, mid, rc int) {
	// We store nodes in Eulerian circuit order.
	mid = (left + right) / 2
	lc = index + 1
	rc = index + 2*(mid-left)
	return
}
