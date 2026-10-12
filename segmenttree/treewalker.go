package segmenttree

// Walker allows walking the segment tree.
type Walker[T any] struct {
	tree        *Tree[T]
	nodeIdx     int
	left, right int
}

// NewWalker creates a walker anchored at the root of the provided tree.
func NewWalker[T any](tree *Tree[T]) *Walker[T] {
	return &Walker[T]{
		tree:    tree,
		nodeIdx: 0,
		left:    0,
		right:   tree.length,
	}
}

// Index returns the current node index of the root.
func (w *Walker[T]) Index() int { return w.nodeIdx }

// Left returns the left endpoint of the interval covered by the current root.
func (w *Walker[T]) Left() int { return w.left }

// Right returns the right endpoint of the interval covered by the current root.
func (w *Walker[T]) Right() int { return w.right }

// MoveLeft moves the walker to the left subtree.
func (w *Walker[T]) MoveLeft() {
	w.nodeIdx, w.right, _ = split(w.nodeIdx, w.left, w.right)
}

// MoveRight moves the walker to the right subtree.
func (w *Walker[T]) MoveRight() {
	_, w.left, w.nodeIdx = split(w.nodeIdx, w.left, w.right)
}

// RootValue returns the value at the current tree root.
func (w *Walker[T]) RootValue() T {
	return w.tree.nodes[w.nodeIdx]
}

// LeftValue is the value at the left child, or the default value if there is no child.
func (w *Walker[T]) LeftValue() (ret T) {
	if w.right-w.left <= 1 {
		return
	}
	lc, _, _ := split(w.nodeIdx, w.left, w.right)
	return w.tree.nodes[lc]
}

// RightValue is the value at the right child, or the default value if there is no child.
func (w *Walker[T]) RightValue() (ret T) {
	if w.right-w.left <= 1 {
		return
	}
	_, _, rc := split(w.nodeIdx, w.left, w.right)
	return w.tree.nodes[rc]
}
