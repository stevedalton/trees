package intervaltree

import (
	"cmp"
	"fmt"
	"slices"
	"testing"

	googcmp "github.com/google/go-cmp/cmp"
)

func cmpInterval(a, b Interval) int {
	return cmp.Compare(a.ID, b.ID)
}

func TestTree(t *testing.T) {
	tree := &Tree{}
	intervals := []Interval{
		{4, 6, "First"},
		{3, 5, "Second"},
		{8, 11, "Third"},
		{2, 20, "Fourth"},
		{1, 10, "Fifth"},
		{1, 9, "Sixth"},
		{1, 8, "Seventh"},
		{4, 4, "Eighth"},
		{1, 7, "Ninth"},
		{7, 40, "Tenth"},
	}
	for _, i := range intervals {
		if err := tree.Insert(i); err != nil {
			t.Errorf("Insert(%v) = %v, want nil", i, err)
		}
		if err := checkTree(tree); err != nil {
			t.Errorf("After insert of %v, tree is not valid: %v", i, err)
		}
	}
	t.Log(tree.String())

	points := []int64{2, 3, 5, 7, 1, 13, 17}
	for _, p := range points {
		got, want := tree.Find(p), fullFind(tree.root, p)
		slices.SortFunc(got, cmpInterval)
		slices.SortFunc(want, cmpInterval)
		if diff := googcmp.Diff(want, got); diff != "" {
			t.Errorf("Find(%d) gave wrong results (-want +got):\n%s\n", p, diff)
		}
	}
}

// fullFind does an exhaustive search of the tree to find all intervals.
func fullFind(n *Node, p int64) []Interval {
	if n == nil {
		return nil
	}
	res := fullFind(n.right, p)
	res = append(res, fullFind(n.left, p)...)
	if n.interval.Lo <= p && p <= n.interval.Hi {
		res = append(res, n.interval)
	}
	return res
}

// checkTree ensures the tree:
// * follows red-black invariant properties.
// * maximums are accurate.
// * is ordered correctly.
func checkTree(t *Tree) error {
	switch {
	case t == nil || t.root == nil:
		return nil
	case t.root.red:
		return fmt.Errorf("root of the tree is red")
	case t.root.parent != nil:
		return fmt.Errorf("node %s is not the root", t.root)
	default:
	}
	_, _, err := checkRecursive(t.root)
	return err
}

// checkRecursive recursively checks the red-black tree rooted at node.
// It watches for bad parent/child links, and also red-black invariants.
// It returns: count of black nodes; max of all intervals; invariant error.
func checkRecursive(n *Node) (int, int64, error) {
	switch {
	case n == nil:
		return 1, 0, nil
	case n.left != nil && n.left.parent != n:
		return 0, 0, fmt.Errorf("node %s has a left child %s with the wrong parent", n, n.left)
	case n.right != nil && n.right.parent != n:
		return 0, 0, fmt.Errorf("node %s has a right child %s with the wrong parent", n, n.right)
	case !n.red:
	case n.left != nil && n.left.red:
		return 0, 0, fmt.Errorf("red node %s has a red left child %s", n, n.left)
	case n.right != nil && n.right.red:
		return 0, 0, fmt.Errorf("red node %s has a red right child %s", n, n.right)
	}
	lBlack, lMax, err := checkRecursive(n.left)
	if err != nil {
		return 0, 0, err
	}
	rBlack, rMax, err := checkRecursive(n.right)
	if err != nil {
		return 0, 0, err
	}
	nMax := max(n.interval.Hi, lMax, rMax)
	switch {
	case nMax != n.max:
		return 0, 0, fmt.Errorf("node %s has max attribute %d that should be %d", n, n.max, nMax)
	case lBlack != rBlack:
		return 0, 0, fmt.Errorf("node %s has black paths on the left of length %d, but black paths on the right of length %d", n, lBlack, rBlack)
	case !n.red:
		lBlack++
	}
	return lBlack, nMax, nil
}
