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

	log.Println(hypercubes)
	for _, cube := range hypercubes {
		log.Println(*cube)
	}

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
