package genetic

import (
	"log"
	"math"
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

// calculateHypervolumeIndicator searches through the whole tree and calculates the hypervolume indicator base point for the given hypercube.
// If the we learn the hypercube is dominated (wholely inside another hypercube) the indicator base is nil.
func (t *hypercubeKdTree) calculateHypervolumeIndicatorBase(hypercube *specimenHypercube) (isDominated bool, newIndicatorBase []float64) {

	// The initial base for a hypervolume indicator is (0.0, 0.0, ...)
	// Initially left branches are viable for all dimensions.
	var indicatorBase []float64 = make([]float64, len(hypercube.dimensions))
	var leftViable []bool = make([]bool, len(hypercube.dimensions))
	for i := range hypercube.dimensions {
		indicatorBase[i] = 0.0
		leftViable[i] = true
	}

	// Now just return what we find when we dive into the tree.
	return t.root.calculateHypervolumeIndicatorBase(hypercube, indicatorBase, leftViable)
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
func (n *hypercubeKdTreeNode) calculateHypervolumeIndicatorBase(hypercube *specimenHypercube, indicatorBase []float64, leftViable []bool) (isDominated bool, newIndicatorBase []float64) {

	// For the moment, the left branch is as viable as in the node above us.
	var newLeftViable []bool = leftViable

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

			// Caculate what the new indicator base is when this nodes hypercube encrouches on our hypervolume indicator and
			// "consumes" part of its volume.
			indicatorBase = moveIndicatorBase(hypercube.dimensions, indicatorBase, n.hypercube.dimensions)
		}

		// We may have discovered that in this dimension, to the left, we have nothing more to learn.
		newLeftViable = updateLeftViable(n.dimension, leftViable, indicatorBase, n.hypercube.dimensions)
	}

	// We have compared our seaching cube against the cube in this node (or we are the cube of this node) and were not dominated.
	// The searching cube's hypervolume indicator base may have moved up a little. Dive into the child branches to find more
	// cubes that shrink our hypervolume indicator.
	var isLeftDominated, isRightDominated bool
	var leftIndicatorBase, rightIndicatorBase []float64

	// Go down left branch if there is one, and there is still more to learn down that branch.
	var leftNotViable bool = !isLeftViable(newLeftViable)
	if n.left == nil || leftNotViable {
		isLeftDominated, leftIndicatorBase = false, indicatorBase
	} else {
		isLeftDominated, leftIndicatorBase = n.left.calculateHypervolumeIndicatorBase(hypercube, indicatorBase, newLeftViable) // The left branch may have just accrued a dimension no longer viable.
	}

	// Go down right branch if there is one.
	if n.right == nil {
		isRightDominated, rightIndicatorBase = false, indicatorBase
	} else {
		isRightDominated, rightIndicatorBase = n.right.calculateHypervolumeIndicatorBase(hypercube, indicatorBase, leftViable) // Keep original left viable.
	}

	// Dominated?
	if isLeftDominated || isRightDominated {
		return true, nil
	}

	// Not dominated, merge what we learned from both branches regarding the hypervolume indicator base.
	indicatorBase = moveIndicatorBase(hypercube.dimensions, leftIndicatorBase, rightIndicatorBase)
	return false, indicatorBase
}

// moveIndicatorBase moves the indicator base towards the limit based on the input. The limit is the point that defines
// the opposite corner of the hypervolume indicator. We cannot move the base past it in any dimension. The input is the
// corner of the other hypercube being considered, it is "consuming" volume from the hypervolume indicator by pushing
// the base in any dimension that it already "owns".
func moveIndicatorBase(limit []float64, base []float64, other []float64) (movedBase []float64) {
	movedBase = make([]float64, len(base))
	for i := range base {

		// If the other cube dominates the searching cube in a given dimension, it has no effect on the
		// indicator base. As long as the searching cube finds no hypercube that completely dominates
		// every dimension it is not dominiated. It has a hypervolume indicator (some contribution it uniquely makes).
		// The hypevolume indicator is only shrunk by dimensions that don't wholey erase the dimension.
		if other[i] < limit[i] {
			// What is highest value from indicator base and new point?
			movedBase[i] = math.Max(base[i], other[i])
		} else {
			// If this is a whole dominated dimension, ignore it. Not part of hypercube indicator volume.
			movedBase[i] = base[i]
		}
	}
	return movedBase
}

// updateLeftViable udpates our bookkeeping to explore the left-hand branch (where numbers get lower). The right-hand
// branch is already accounted for by checking if the searching cube is dominated.
func updateLeftViable(dimension int, leftViable []bool, base []float64, other []float64) []bool {
	// The numbers in the base always move higher. If the number in the current dimenions is less than or equal to the
	// base, that dimension is no longer of interest for the hypervolume indicator (at least in left branch searches
	// where numbers always go lower).
	var newLeftViable []bool = make([]bool, len(leftViable))
	copy(newLeftViable, leftViable)
	if base[dimension] >= other[dimension] {
		// The base is already equal to or higher in this dimension.
		// We cannot improve the base by checking left-hand branches.
		newLeftViable[dimension] = false
	}
	return newLeftViable
}

// isLeftViable checks out bookkeeping data, if any dimension is still viable in a left hand branch keep digging.
// The moment no dimensino is of interested in the left hand branches, stop going down them.
func isLeftViable(leftViable []bool) bool {
	for _, viable := range leftViable {
		if viable {
			return true
		}
	}
	return false
}

// kdCompareHypercubes compares two hypercubes at a node in the k-d tree and learns everything it can from them, altering them as needed to
// to capture the learnings. Besides altering the hypercubes themselves, we learn if there is more to learn for the searching hypercube
// down either or both of the branches.
func kdCompareHypercubes(searchingHypercube *specimenHypercube, nodeHypercube *specimenHypercube, maximumLeft []float64) (isLeftSearchable bool) {

	// If the cubes are the same cube, just bail.
	if searchingHypercube == nodeHypercube {
		// Assume we need to keep searching left.
		return true
	}

	// If we have already compared these cubes, bail.
	var alreadyCompared bool
	if _, alreadyCompared = searchingHypercube.comparedWith[nodeHypercube]; alreadyCompared {
		// Assume we need to keep searching left.
		return true
	}

	// Assume the left branch of this node is not searchabe. Discover if it is.
	isLeftSearchable = false

	// If we have no sense of what's down the left branch, we need to search it.
	if maximumLeft == nil {
		isLeftSearchable = true
	}

	// Either cube may already be dominated before we compare.
	var isSearchingPriorDominated bool = searchingHypercube.isDominated
	var isNodePriorDominated bool = nodeHypercube.isDominated

	// Compare the searching hypercube agaisnt the node hypercube. Each of the cubes can shrink the other's hypervolume indicator by
	// moving the base of the hypervolume indicator cube closer to the dimension point that defines the specimen's hypercube.
	// Compute a lot of related information from comparing the cubes.
	var searchingBaseIndicator, nodeBaseIndicator []float64 // Calculate new indicator bases.
	var isSearchingDominated bool = true                    // Determine if cubes are dominated.
	var isNodeDominated bool = true                         // Determine if cubes are dominated.
	var isEqual bool = true                                 // Determine if the cubes are not identical.
	for i := 0; i < len(searchingHypercube.dimensions); i++ {

		// What are the values at this dimension?
		var searchingDimensionN float64 = searchingHypercube.dimensions[i]
		var nodeDimensionN float64 = nodeHypercube.dimensions[i]

		// Different hypercubes?
		if searchingDimensionN != nodeDimensionN {
			isEqual = false
		}

		// Break domination, if the cubes are no prior dominated.
		if !isSearchingPriorDominated && searchingDimensionN > nodeDimensionN {
			// At least one dimension is not within the other cube.
			isSearchingDominated = false
		}
		if !isNodePriorDominated && nodeDimensionN > searchingDimensionN {
			// At least one dimension is not within the other cube.
			isNodeDominated = false
		}

		// Move the base indicator of the node hypercube if it is not dominated.
		if !isNodePriorDominated {

			// Ignore values greater than the cube's dimension, they don't define the hypervolume indicator.
			var nodeIndicatorBaseN float64 = nodeHypercube.indicatorBase[i]
			if nodeIndicatorBaseN < searchingDimensionN && searchingDimensionN < nodeDimensionN {
				// The searching cube's dimension is between the node cube's dimension and indicator base.
				// Move the indicator base up to the new dimension.
				nodeIndicatorBaseN = searchingDimensionN
			}
			nodeBaseIndicator = append(nodeBaseIndicator, nodeIndicatorBaseN)
		}

		// Move the base indicator of the searching hypercube if it is not dominated.
		if !isSearchingPriorDominated {

			// Ignore values greater than the cube's dimension, they don't define the hypervolume indicator.
			var searchingIndicatorBaseN float64 = searchingHypercube.indicatorBase[i]
			if searchingIndicatorBaseN < nodeDimensionN && nodeDimensionN < searchingDimensionN {
				// The node cube's dimension is between the searching cube's dimension and indicator base.
				// Move the indicator base up to the new dimension.
				searchingIndicatorBaseN = nodeDimensionN
			}
			searchingBaseIndicator = append(searchingBaseIndicator, searchingIndicatorBaseN)

			// We may have a sense of the biggest dimensions to the left.
			// If we do have a sense of them, only search left branches if cubes down
			// those branches can tell us more about the searching cube's hypervolume indicator.
			if maximumLeft != nil && maximumLeft[i] > searchingIndicatorBaseN {
				// Something down the left branch could push push our hypervolume indicator base higher.
				isLeftSearchable = true
			}
		}
	}

	// If the two cubes are equal, the searching cube dominates the node cube.
	if isEqual {
		isSearchingDominated = false
		isNodeDominated = true
	}

	// Remember our two indicators.
	if isSearchingDominated || isSearchingPriorDominated {
		searchingHypercube.isDominated = true
		searchingHypercube.indicatorBase = nil
		// If the searching cube is dominated there is nothing more to the left to check.
		isLeftSearchable = false
	} else {
		searchingHypercube.indicatorBase = searchingBaseIndicator
	}

	// Remember our two indicators.
	if isNodeDominated || isNodePriorDominated {
		nodeHypercube.isDominated = true
		nodeHypercube.indicatorBase = nil
		// If the node cube is dominated, there still may be cubes to the left that could move the indicator.
	} else {
		nodeHypercube.indicatorBase = nodeBaseIndicator
		// Remember on the node cube that we compared these cubes.
		// When its the node cube's time to search it can skip the tests against the current searching cube.
		if nodeHypercube.comparedWith == nil {
			nodeHypercube.comparedWith = map[*specimenHypercube]bool{}
		}
		nodeHypercube.comparedWith[searchingHypercube] = true
	}

	// // For the moment, the left branch is as viable as in the node above us.
	// var newLeftViable []bool = leftViable
	//
	// // Only examine this node's hypercube if it's not the one we're currently checking. This is a process of checking other hypercubes
	// // but the one we're checking is somewhere in the k-d tree.
	// if hypercube != n.hypercube {
	// 	// This node's hypercube is not us.
	//
	// 	// First, just check for simple domination regarding the hypercube of this node.
	// 	switch {
	//
	// 	case hypercube.equals(n.hypercube):
	// 		// We run into a potential issue here where two hypercubes have the exact same dimensions.
	// 		// Only one of these hypercubes will have the hypervolume indicator. The other will be considered a non-contributor.
	// 		// When two hypercubes have exact same dimensions, we want the current searching hypercube to be the "winner".
	// 		// The other hypercube hasn't been searched yet. If it had then one of these cubes would already be dominated.
	// 		// As-in the current searching cube would have been dominated, in which case it would not be in the search right now.
	// 		n.hypercube.isDominated = true // This is a pointer manipulation altering the top-level hypercube list.
	//
	// 	case hypercube.isDominatedBy(n.hypercube):
	// 		// We're not equal and this hypercube exists wholely inside another cube. We don't know whether some other cube
	// 		// has or will dominate the other cube, be we are definitely a non-contributor to the population's hypervolume.
	// 		hypercube.isDominated = true // This is a pointer manipulation altering the top-level hypercube list.
	//
	// 	case n.hypercube.isDominatedBy(hypercube):
	// 		// We're not equal and this nodes hypercube wholely exists inside the searching hypercube. It is dominated.
	// 		n.hypercube.isDominated = true
	// 	}
	//
	// 	// If we are dominated, short circuit out of the search. We have the answer we need for this hypercube.
	// 	if hypercube.isDominated {
	// 		return true, nil
	// 	}
	//
	// 	// We're not dominated. If this node's hypercube is also not dominated, it can tell us something about our hypervolume
	// 	// indicator. If this node's hypercube is dominated, what it has to tell us is irrelevant. The hypercube that dominated
	// 	//  it will be the one that shapes our hypervolume indicator (or it was the current searching cube).
	// 	if !n.hypercube.isDominated {
	//
	// 		// Caculate what the new indicator base is when this nodes hypercube encrouches on our hypervolume indicator and
	// 		// "consumes" part of its volume.
	// 		indicatorBase = moveIndicatorBase(hypercube.dimensions, indicatorBase, n.hypercube.dimensions)
	// 	}
	//
	// 	// We may have discovered that in this dimension, to the left, we have nothing more to learn.
	// 	newLeftViable = updateLeftViable(n.dimension, leftViable, indicatorBase, n.hypercube.dimensions)
	// }
	//
	// // We have compared our seaching cube against the cube in this node (or we are the cube of this node) and were not dominated.
	// // The searching cube's hypervolume indicator base may have moved up a little. Dive into the child branches to find more
	// // cubes that shrink our hypervolume indicator.
	// var isLeftDominated, isRightDominated bool
	// var leftIndicatorBase, rightIndicatorBase []float64
	//
	// // Go down left branch if there is one, and there is still more to learn down that branch.
	// var leftNotViable bool = !isLeftViable(newLeftViable)
	// if n.left == nil || leftNotViable {
	// 	isLeftDominated, leftIndicatorBase = false, indicatorBase
	// } else {
	// 	isLeftDominated, leftIndicatorBase = n.left.calculateHypervolumeIndicatorBase(hypercube, indicatorBase, newLeftViable) // The left branch may have just accrued a dimension no longer viable.
	// }
	//
	// // Go down right branch if there is one.
	// if n.right == nil {
	// 	isRightDominated, rightIndicatorBase = false, indicatorBase
	// } else {
	// 	isRightDominated, rightIndicatorBase = n.right.calculateHypervolumeIndicatorBase(hypercube, indicatorBase, leftViable) // Keep original left viable.
	// }
	//
	// // Dominated?
	// if isLeftDominated || isRightDominated {
	// 	return true, nil
	// }
	//
	// // Not dominated, merge what we learned from both branches regarding the hypervolume indicator base.
	// indicatorBase = moveIndicatorBase(hypercube.dimensions, leftIndicatorBase, rightIndicatorBase)

	return isLeftSearchable
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
