package genetic

import (
	"log"
	"sort"
)

// hypercubeKdTree is a K-D tree, a search tree useful for searching across multiple dimensions.
// Ror specimen hypercubes, it simplifies searching for overlapping hypercubes.
// Reference: https://en.wikipedia.org/wiki/K-d_tree
type hypercubeKdTree struct {
	root hypercubeKdTreeNode // The root of the search tree.
}

// newHypecubeKdTree creates a new multi-dimensional search tree for hypercubes.
func newHypecubeKdTree(hypercubes []*specimenHypercube) hypercubeKdTree {

	// Validate parameters.
	if len(hypercubes) == 0 {
		log.Panic("ERROR: newHypecubeKdTree called with no hypercubes")
	}

	// How many dimensions are there to these hypercube?
	// All the hypercubes have the same dimensions, just examine the first.
	var dimensions int = len(hypercubes[0].dimensions)

	// There must be at least some dimension.
	if dimensions == 0 {
		log.Panic("ERROR: newHypecubeKdTree called with hypercubes that have no dimensions")
	}

	// Build the tree.
	return hypercubeKdTree{
		root: newHypercubeKdTreeNode(hypercubes, dimensions, 0),
	}
}

// hypercubeKdTreeNode is a single node in the search tree. Some portion of the hypercube population is
// in this node, split over a particular dimension into a left and right branch.
type hypercubeKdTreeNode struct {
	dimension int                  // The dimension this node was split on.
	hypercube *specimenHypercube   // The hypercube at this node.
	left      *hypercubeKdTreeNode // The hypercubes branch of this node with a dimensional value less than this node's hypercube.
	right     *hypercubeKdTreeNode // The hypercubes branch of this node with a dimensional value greater than this node's hypercube.
}

// newHypercubeKdTreeNode creates a new node recusively from remaining hypercubes, recursively.
func newHypercubeKdTreeNode(hypercubes []*specimenHypercube, dimensions int, depth int) hypercubeKdTreeNode {

	// For each level the tree, sort and split the remaining hypercubes on a new dimension,
	// eventually starting over at the first dimension again.
	var splitDimension int = depth % dimensions

	// Is their only one hypercube left?
	if len(hypercubes) == 1 {
		return hypercubeKdTreeNode{
			dimension: splitDimension,
			hypercube: hypercubes[0],
		}
	}

	// More than one hypercube, we need to split them on the dimension in question.

	// Sort the hypercubes on the dimension we will split them by.
	var hypercubeSort byDimensionHypercubeSort = byDimensionHypercubeSort{
		dimension:  splitDimension,
		hypercubes: hypercubes,
	}
	sort.Sort(hypercubeSort)
	var sortedHypercubes []*specimenHypercube = hypercubeSort.hypercubes

	// Find the median hypercube, that is the hypercube for this node.
	// Integer division drops the remainder so the median index will calculate like this:
	//
	//    5 / 2 = 0 1 (2) 3 4
	//    4 / 2 = 0 1 (2) 3
	//    3 / 2 = 0 (1) 2
	//    2 / 2 = 0 (1) <--- here one child branch of the node will be empty.
	//    1 / 2 = (0) <-- This case already has exited this function, but both branches would be empty.
	//
	var medianIndex int = len(sortedHypercubes) / 2

	// Build the details of this node.
	var node hypercubeKdTreeNode = hypercubeKdTreeNode{
		dimension: splitDimension,
		hypercube: sortedHypercubes[medianIndex],
	}

	// What remaining hypercubes go into each branch?
	var leftHypercubes []*specimenHypercube = sortedHypercubes[:medianIndex]
	var rightHypercubes []*specimenHypercube = sortedHypercubes[medianIndex+1:]

	// Make a sanity check to ensure we have done our splitting right.
	if len(hypercubes) != len(leftHypercubes)+1+len(rightHypercubes) {
		log.Panic("ASSERT: Lost hypercubes on split, bug in code.")
	}

	// Anything for the left branch?
	if len(leftHypercubes) > 0 {
		var left hypercubeKdTreeNode = newHypercubeKdTreeNode(leftHypercubes, dimensions, depth+1)
		node.left = &left
	}

	// Anything for the right branch?
	if len(rightHypercubes) > 0 {
		var right hypercubeKdTreeNode = newHypercubeKdTreeNode(rightHypercubes, dimensions, depth+1)
		node.right = &right
	}

	// Return our finished node and its subtrees.
	return node
}

// calculateHypervolumeIndicator searches this node and its children for hypercubes that would increase the size of the hypervolume indicator's base.
// the final hypervolume indicator will be the volume of the hypercube created by the source hypercube's defining point in one corner and the indicator's
// base in the opposite corner. As we discover hypercubes that overlap our current indicator, we move the indicator base "higher", closer to the defining
// point, and so shrink the eventual hypervolume indicator.
func (n *hypercubeKdTreeNode) calculateHypervolumeIndicatorBase(hypercube *specimenHypercube, indicatorBase []float64) (isDominated bool, newIndicatorBase []float64) {

	// Only examine this node's hypercube if it's not the one we're currently checking. This is a process of checking other hypercubes
	// but the one we're checking is somewhere in the k-d tree.
	if hypercube != n.hypercube {
		// This node's hypercube is not us.

		// First, just check for simple domination regarding the hypercube of this node.
		switch {

		case hypercube.equals(n.hypercube):
			// We run into a potential issue here where two hypercubes have the exact same dimensions.
			// Only one of these hypercubes will have the hypervolume indicator. The other will be considered a non-contributor.
			// When two hypercubes have exact same dimensions, we want the current searching hypercube to be the "winner".
			// The other hypercube hasn't been searched yet. If it had then one of these cubes would already be dominated.
			// As-in the current searching cube would have been dominated, in which case it would not be in the search right now.
			n.hypercube.isDominated = true // This is a pointer manipulation altering the top-level hypercube list.

		case hypercube.isDominatedBy(n.hypercube):
			// We're not equal and this hypercube exists wholely inside another cube. We don't know whether some other cube
			// has or will dominate the other cube, be we are definitely a non-contributor to the population's hypervolume.
			hypercube.isDominated = true // This is a pointer manipulation altering the top-level hypercube list.

		case n.hypercube.isDominatedBy(hypercube):
			// We're not equal and this nodes hypercube wholely exists inside the searching hypercube. It is dominated.
			n.hypercube.isDominated = true
		}

		// If we are dominated, short circuit out of the search. We have the answer we need for this hypercube.
		if hypercube.isDominated {
			return true, nil
		}

		// We're not dominated. If this node's hypercube is also not dominated, it can tell us something about our hypervolume
		// indicator. If this node's hypercube is dominated, what it has to tell us is irrelevant. The hypercube that dominated
		//  it will be the one that shapes our hypervolume indicator (or it was the current searching cube).
		if !n.hypercube.isDominated {

		}

		// Does this node

	}

	return false, nil
}

// Maximize dimensions takes two sets of dimensions and makes a new one with the highest value for each axis.
func maximizeDimensions(a []float64, b []float64) (maximized []float64) {
	maximized = make([]float64, len(a))
	for i := range a {
		if a[i] > b[i] {
			maximized[i] = a[i]
		} else {
			maximized[i] = b[i]
		}
	}
	return maximized
}

// byDimensionHypercubeSort implements sort.Interface to sort hypercubes ascending by a particular dimension.
type byDimensionHypercubeSort struct {
	dimension  int                  // What dimension are we sorting on?
	hypercubes []*specimenHypercube // What hypercubes are we sorting?
}

// The sort.Interface.
func (a byDimensionHypercubeSort) Len() int { return len(a.hypercubes) }
func (a byDimensionHypercubeSort) Swap(i, j int) {
	a.hypercubes[i], a.hypercubes[j] = a.hypercubes[j], a.hypercubes[i]
}
func (a byDimensionHypercubeSort) Less(i, j int) bool {
	return a.hypercubes[i].dimensions[a.dimension] < a.hypercubes[j].dimensions[a.dimension]
}
