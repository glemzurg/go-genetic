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
		root: newHypercubeKdTreeNode(hypercubes, dimensions, 0, nil),
	}
}

// calculateHypervolumeIndicator searches through the whole tree and calculates the hypervolume indicator base point for the given hypercube.
// If the we learn the hypercube is dominated (wholely inside another hypercube) the indicator base is nil.
func (t *hypercubeKdTree) calculateHypervolumeIndicatorBase(hypercube *specimenHypercube) {

	// Walk the tree with the hypercube, updating its hypervolume indicator base.
	t.root.calculateHypervolumeIndicatorBase(hypercube)
}

// hypercubeKdTreeNode is a single node in the search tree. Some portion of the hypercube population is
// in this node, split over a particular dimension into a left and right branch.
type hypercubeKdTreeNode struct {
	dimension   int                  // The dimension this node was split on.
	hypercube   *specimenHypercube   // The hypercube at this node.
	maximumLeft []float64            // The maximum hypercube that could exist in a left branch (the branch may be nil though). Nil if we're not deep enough in tree to know.
	isDominated bool                 // Are all the cubes under this branch dominated by another cube?
	left        *hypercubeKdTreeNode // The hypercubes branch of this node with a dimensional value less than this node's hypercube.
	right       *hypercubeKdTreeNode // The hypercubes branch of this node with a dimensional value greater than this node's hypercube.
}

// newHypercubeKdTreeNode creates a new node recusively from remaining hypercubes, recursively.
func newHypercubeKdTreeNode(hypercubes []*specimenHypercube, dimensions int, depth int, maximumLeft []float64) hypercubeKdTreeNode {

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

		// We have cubes to the left, compute the the maximum possible hypercube to the left.
		// What is the highest value in this dimension?
		var lastLeftHypercube *specimenHypercube = leftHypercubes[len(leftHypercubes)-1]
		var maximumLeftValue float64 = lastLeftHypercube.dimensions[splitDimension]

		// Add this information to the maximum left value.
		var newMaximumLeft []float64 = calculateNewMaximumLeft(maximumLeft, dimensions, splitDimension, maximumLeftValue)

		var left hypercubeKdTreeNode = newHypercubeKdTreeNode(leftHypercubes, dimensions, depth+1, newMaximumLeft)
		node.left = &left

		// Only remember the maximum left if it is complete hypercube, otherwise nil means the node doens't know
		// enough about what's on the left-hand branch.
		if len(newMaximumLeft) == dimensions {
			node.maximumLeft = newMaximumLeft
		}
	}

	// Anything for the right branch?
	if len(rightHypercubes) > 0 {

		// We have cubes to the right, compute the the maximum possible hypercube to the right.
		// What is the highest value in this dimension?
		var lastRightHypercube *specimenHypercube = rightHypercubes[len(rightHypercubes)-1]
		var maximumRightValue float64 = lastRightHypercube.dimensions[splitDimension]

		// Add this information to the maximum left value.
		var newMaximumLeft []float64 = calculateNewMaximumLeft(maximumLeft, dimensions, splitDimension, maximumRightValue)

		var right hypercubeKdTreeNode = newHypercubeKdTreeNode(rightHypercubes, dimensions, depth+1, newMaximumLeft)
		node.right = &right
	}

	// Return our finished node and its subtrees.
	return node
}

// calculateNewMaximumLeft updates the maximum left with a new value indicating the highest value in a specific
// dimension down the next left-hand branch. If the full maximum left isn't filled out yet, just keep building it.
func calculateNewMaximumLeft(maximumLeft []float64, dimensions int, splitDimension int, maximumLeftValue float64) []float64 {
	var newMaximumLeft []float64 = make([]float64, len(maximumLeft))
	copy(newMaximumLeft, maximumLeft)

	// Update the new maximum left value.
	// If we have yet to get a set of dimensions, keep building.
	if len(newMaximumLeft) < dimensions {

		// Sanity check, if we are building up the list the length should be equal to the split dimension
		// which would be one less than the dimension we are working with (dimension 1 is split dimension index 0)
		if len(newMaximumLeft) != splitDimension {
			log.Panicf("ASSERT: invalid maximum left append: %d %f", splitDimension, newMaximumLeft)
		}

		// Add the next dimension.
		newMaximumLeft = append(newMaximumLeft, maximumLeftValue)

	} else {

		// Sanity check, if we are updating an existing value, we can only ever make it smaller.
		if maximumLeftValue > newMaximumLeft[splitDimension] {
			log.Panicf("ASSERT: invalid maximum left update: %f %v", maximumLeftValue, newMaximumLeft)
		}

		// Update the value in the given position.
		newMaximumLeft[splitDimension] = maximumLeftValue
	}

	return newMaximumLeft
}

// calculateHypervolumeIndicator searches this node and its children for hypercubes that would increase the size of the hypervolume indicator's base.
// the final hypervolume indicator will be the volume of the hypercube created by the source hypercube's defining point in one corner and the indicator's
// base in the opposite corner. As we discover hypercubes that overlap our current indicator, we move the indicator base "higher", closer to the defining
// point, and so shrink the eventual hypervolume indicator.
func (n *hypercubeKdTreeNode) calculateHypervolumeIndicatorBase(hypercube *specimenHypercube) {

	// Update the searching hypercube and the node hypercube's hypervolume indicators by comparing them.
	var isLeftSearchable, isLeftDominated bool = kdCompareHypercubes(hypercube, n.hypercube, n.maximumLeft)

	// Are all the left branch cubes dominated? They may still move the hypervolume indicator so have to be searched.
	// But since we know, update them accordingly.
	if isLeftDominated {
		n.left.dominateWholeBranch()
	}

	// Once we discover the searching hypercube is dominated, there is nothing more to learn about its
	// hypervolume indicator... it doesn't have one. Only continue if we are not dominated.
	if !hypercube.isDominated {

		// Let's look at larger hypercubes first. We may discover we're dominated.
		if n.right != nil {
			n.right.calculateHypervolumeIndicatorBase(hypercube)
		}

		// Still not dominated?
		if !hypercube.isDominated {

			// We may have discovered there is nothing more to learn going down the left branch.
			// Also, if the left branch is dominated, it cannot shape hypervolume indicators.
			if n.left != nil && isLeftSearchable && !isLeftDominated {
				n.left.calculateHypervolumeIndicatorBase(hypercube)
			}
		}
	}
}

// dominateWholeBranch flags every cube in this branch of the k-d tree as dominated (existing inside another cube).
func (n *hypercubeKdTreeNode) dominateWholeBranch() {
	if n != nil {
		if !n.isDominated {
			n.hypercube.setDominated()
			n.left.dominateWholeBranch()
			n.right.dominateWholeBranch()
			n.isDominated = true
		}
	}
}

// kdCompareHypercubes compares two hypercubes at a node in the k-d tree and learns everything it can from them, altering them as needed to
// to capture the learnings. Besides altering the hypercubes themselves, we learn if there is more to learn for the searching hypercube
// down the left branch (smaller cubes) as well as whether all the left branch cubes are dominated by this.
func kdCompareHypercubes(searchingHypercube *specimenHypercube, nodeHypercube *specimenHypercube, maximumLeft []float64) (isLeftSearchable bool, isLeftDominated bool) {

	log.Println("=================")
	log.Println(searchingHypercube)
	log.Println(nodeHypercube)

	// If the cubes are the same cube, just bail.
	if searchingHypercube == nodeHypercube {
		// Assume we need to keep searching left and nothing to the left is dominated.
		return true, false
	}

	// If either cube is dominated, there is nothing to learn about the hypervolume indicators from comparing them.
	if searchingHypercube.isDominated {
		// Search cube is dominated so there is nothing more to search to the left. Don't know whether left branch is dominated.
		return false, false
	}
	if nodeHypercube.isDominated {
		// Assume we need to keep searching left and nothing to the left is dominated.
		return true, false
	}

	// If we have already compared these cubes, bail.
	var alreadyCompared bool
	if _, alreadyCompared = searchingHypercube.comparedWith[nodeHypercube]; alreadyCompared {
		// Assume we need to keep searching left and nothing to the left is dominated.
		return true, false
	}

	// Assume the left branch of this node is not searchabe. Discover if it is.
	// Asseme the left branch of this node only includes dominated cubes. Discover if it is not true.
	isLeftSearchable = false
	isLeftDominated = true

	// If we have no sense of what's down the left branch, we need to search it.
	// Can't know if the left branch cubes are dominated.
	if maximumLeft == nil {
		isLeftSearchable = true
		isLeftDominated = false
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
			if maximumLeft != nil {

				// Only search left branches if cubes down those branches can tell us more about the
				// searching cube's hypervolume indicator.
				if maximumLeft[i] > searchingIndicatorBaseN {
					// Something down the left branch could push push our hypervolume indicator base higher.
					isLeftSearchable = true
				}

				// Check to see if the left branched cubes are dominated. The maximum left is the biggest
				// possible hypercube down the left-hand branches.
				if maximumLeft[i] > searchingDimensionN {
					// The left-hand branch can contain a cube that is not be dominated.
					isLeftDominated = false
				}
			}
		}
	}

	// If the two cubes are equal, the searching cube dominates the node cube.
	if isEqual {
		isSearchingDominated = false
		isNodeDominated = true
	}

	// Regardless of what we just learned about our cubes, we cannot undo domination.
	if isSearchingPriorDominated {
		isSearchingDominated = true
	}
	if isNodePriorDominated {
		isNodeDominated = true
	}

	// Remember our two indicators.
	if isSearchingDominated {
		searchingHypercube.setDominated()
		// If the searching cube is dominated there is nothing more to the left to check.
		isLeftSearchable = false
	} else {
		// If the node cube is dominated it cannot tell us anything about our hypervolume indicator base.
		// Only cubes that define the outer edge of the hypercube populate can shape the hypervolume indicator.
		// Consider, that if hidden cubes could then a mass of little shapes would overlay and shrink our
		// indicator to something very small, not representative of how much volume this hypercube has which
		// doesn't overlay other cubes.
		if !isNodeDominated {
			searchingHypercube.indicatorBase = searchingBaseIndicator
		}
	}

	// Remember our two indicators.
	if isNodeDominated {
		nodeHypercube.setDominated()
		// If the node cube is dominated, there still may be cubes to the left that could move the indicator.
	} else {
		// If the other cube is not dominated remember our indicator.
		// If hte other cube is dominated, it can't influence indicators.
		if !isSearchingDominated {
			nodeHypercube.indicatorBase = nodeBaseIndicator
			// Remember on the node cube that we compared these cubes.
			// When its the node cube's time to search it can skip the tests against the current searching cube.
			if nodeHypercube.comparedWith == nil {
				nodeHypercube.comparedWith = map[*specimenHypercube]bool{}
			}
			nodeHypercube.comparedWith[searchingHypercube] = true
		}
	}

	return isLeftSearchable, isLeftDominated
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
