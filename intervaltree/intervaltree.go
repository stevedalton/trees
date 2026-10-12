// Package intervaltree implements a left-leaning red-black tree
// that acts as an interval tree for inserting and searching.
package intervaltree

import (
	"fmt"
	"strings"
)

var s = strings.Repeat(" ", 80)

// Interval represents the closed interval [Lo, Hi].
type Interval struct {
	Lo, Hi int64
	ID     string
}

// Node is a node of a binary tree with Interval data payload.
type Node struct {
	parent      *Node
	left, right *Node
	interval    Interval
	max         int64 // Max element contained in any interval of tree rooted at this node.
	red         bool  // Is the node red?
}

// leftChild returns true if the node is the left child of its parent.
func (n *Node) leftChild() bool {
	if n == nil || n.parent == nil {
		return false
	}
	return n.parent.left == n
}

// fixMax updates the maximum for the specified node.
func (n *Node) fixMax() {
	if n == nil {
		return
	}
	m := n.interval.Hi
	if n.left != nil {
		m = max(m, n.left.max)
	}
	if n.right != nil {
		m = max(m, n.right.max)
	}
	n.max = m
}

// fixMax fixes maxes for all nodes from the passed node up to the root.
func fixMax(n *Node) {
	for n != nil {
		n.fixMax()
		n = n.parent
	}
}

// insert inserts the provided node as a child of the current node.
func (n *Node) insert(c *Node) {
	if n == nil || c == nil {
		return
	}
	var less bool
	// The new node is less than the current node if:
	// * new Lo is less than current Lo
	// * Lo is equal but new Hi is less than current Hi
	switch {
	case n.interval.Lo != c.interval.Lo:
		less = n.interval.Lo > c.interval.Lo
	default:
		less = n.interval.Hi > c.interval.Hi
	}
	switch {
	case less && n.left == nil:
		// Add as brand new left subtree.
		n.left = c
		c.parent = n
	case less:
		// New node belongs in existing left subtree.
		n.left.insert(c)
	case n.right == nil:
		// Add as brand new right subtree.
		n.right = c
		c.parent = n
	default:
		// New node belongs in existing right subtree.
		n.right.insert(c)
	}
}

// insertFix fixes the tree after an insert.
// It updates rotations and colorings to preserve red-black invariants.
// It returns a potential root for the tree.
func insertFix(n *Node) *Node {
	// We only need to fix from the current node up through its ancestors.
	// We worry about red parent-child relationships, as well as the same
	// number of red nodes on every path.
	// Note that this method is always called on a red node.
	switch uncle := n.Uncle(); {
	case n == nil:
		return nil
	case n.parent == nil:
		// Root node is black.
		n.red = false
		return n
	case !n.parent.red:
		// If the parent is not red, we are done fixing upward.
		return n.parent
	case uncle != nil && uncle.red:
		// Uncle is red and parent is red.
		// Thus paths through this node will have one more red
		// than paths through the uncle.
		// Make both black and continue to fix in their common parent.
		n.parent.red = false
		uncle.red = false
		n.parent.parent.red = true
		return insertFix(n.parent.parent)
	}
	// At this point we've fixed everything above us.
	// Now fix the tree rooted here. We need to perform at most 2 rotations.
	left := n.leftChild()
	if parentLeft := n.parent.leftChild(); left != parentLeft {
		n = n.parent
		if parentLeft {
			n.RotateLeft()
		} else {
			n.RotateRight()
		}
		left = !left
	}
	if left {
		n.parent.parent.RotateRight()
	} else {
		n.parent.parent.RotateLeft()
	}
	// The rotations above ensure we have a sibling.
	// Ensure we are both red, and parent is black.
	// This combined with the upward fix earlier finishes
	// fixing the tree.
	n.Sibling().red = true
	n.parent.red = false
	return n.parent
}

// Sibling returns the sibling of the current node, if it exists.
func (n *Node) Sibling() *Node {
	switch {
	case n == nil || n.parent == nil:
		return nil
	case n.leftChild():
		return n.parent.right
	default:
		return n.parent.left
	}
}

// Uncle returns the sibling of the parent node, if it exists.
func (n *Node) Uncle() *Node {
	if n == nil || n.parent == nil {
		return nil
	}
	return n.parent.Sibling()
}

// Find returns all intervals in the tree containing the specified point.
func (n *Node) Find(p int64) []Interval {
	if n == nil || p > n.max {
		return nil
	}
	res := n.left.Find(p)
	if p < n.interval.Lo {
		// Every node in the right subtree has Lo >= this node's Lo.
		return res
	}
	res = append(res, n.right.Find(p)...)
	if p <= n.interval.Hi {
		res = append(res, n.interval)
	}
	return res
}

// RotateLeft does a left rotation of the tree rooted at this node.
func (n *Node) RotateLeft() {
	if n == nil || n.right == nil {
		return
	}
	n.right.parent = n.parent
	if n.parent != nil {
		if n.leftChild() {
			n.parent.left = n.right
		} else {
			n.parent.right = n.right
		}
	}
	n.parent = n.right
	n.right = n.parent.left
	if n.right != nil {
		n.right.parent = n
	}
	n.parent.left = n
	n.fixMax()
	n.parent.fixMax()
}

// RotateRight does a right rotation of the tree rooted at this node.
func (n *Node) RotateRight() {
	if n == nil || n.left == nil {
		return
	}
	n.left.parent = n.parent
	if n.parent != nil {
		if n.leftChild() {
			n.parent.left = n.left
		} else {
			n.parent.right = n.left
		}
	}
	n.parent = n.left
	n.left = n.parent.right
	if n.left != nil {
		n.left.parent = n
	}
	n.parent.right = n
	n.fixMax()
	n.parent.fixMax()
}

// Print prints the tree rooted at this node.
// It uses square brackets ([,]) for black nodes, and
// curly brackets ({,}) for red nodes.
func (n *Node) Print(depth int) string {
	if n == nil {
		return ""
	}
	var builder strings.Builder
	b, e := "[", "]"
	if n.red {
		b, e = "{", "}"
	}
	if n.right != nil {
		builder.WriteString(n.right.Print(depth + 1))
	}
	builder.WriteString(fmt.Sprintf("%s%s%2d,%2d%s\n", s[:3*depth], b, n.interval.Lo, n.interval.Hi, e))
	if n.left != nil {
		builder.WriteString(n.left.Print(depth + 1))
	}
	return builder.String()
}

func (n *Node) String() string {
	if n == nil {
		return ""
	}
	return n.interval.ID
}

// Tree encapsulates a rooted interval tree.
type Tree struct {
	root *Node
}

// Find finds all intervals in the tree containing the point.
func (t *Tree) Find(p int64) []Interval {
	return t.root.Find(p)
}

// Insert inserts the specified interval into the tree.
func (t *Tree) Insert(i Interval) error {
	if t == nil {
		return fmt.Errorf("cannot insert in nil tree")
	}
	if i.Hi < i.Lo {
		return fmt.Errorf("[%d,%d] is not a valid interval", i.Lo, i.Hi)
	}
	n := &Node{
		interval: i,
		max:      i.Hi,
		red:      true,
	}

	if t.root == nil {
		t.root = n
	} else {
		t.root.insert(n)
	}
	// Fix maximums, as well as red-black properties.
	fixMax(n.parent)
	if n = insertFix(n); n.parent == nil {
		t.root = n
	}
	return nil
}

func (t *Tree) String() string {
	if t == nil || t.root == nil {
		return ""
	}
	return t.root.Print(0)
}
