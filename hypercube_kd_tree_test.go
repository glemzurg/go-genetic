package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type HypercubeKdTreeSuite struct{}

var _ = Suite(&HypercubeKdTreeSuite{})

// Add the tests.

func (s *HypercubeKdTreeSuite) Test_NewHypecubeKdTree(c *C) {

	// Use the example from: https://en.wikipedia.org/wiki/K-d_tree
	//
	// [(2,3), (5,4), (9,6), (4,7), (8,1), (7,2)]
	//
	var cube23 specimenHypercube = specimenHypercube{dimensions: []float64{2.0, 3.0}}
	var cube54 specimenHypercube = specimenHypercube{dimensions: []float64{5.0, 4.0}}
	var cube96 specimenHypercube = specimenHypercube{dimensions: []float64{9.0, 6.0}}
	var cube47 specimenHypercube = specimenHypercube{dimensions: []float64{4.0, 7.0}}
	var cube81 specimenHypercube = specimenHypercube{dimensions: []float64{8.0, 1.0}}
	var cube72 specimenHypercube = specimenHypercube{dimensions: []float64{7.0, 2.0}}
	var hypercubes []*specimenHypercube = []*specimenHypercube{&cube23, &cube54, &cube96, &cube47, &cube81, &cube72}

	// Build the tree.
	var kdTree hypercubeKdTree = newHypecubeKdTree(hypercubes)

	var emptyNode *hypercubeKdTreeNode // Leave unallocated.

	// Computes to this tree from: https://en.wikipedia.org/wiki/K-d_tree
	//
	// ((7, 2),
	//   ((5, 4),
	//      ((2, 3), None, None),
	//      ((4, 7), None, None)),
	//   ((9, 6),
	//      ((8, 1), None, None), None))
	//
	c.Assert(kdTree.root.dimension, Equals, 0)
	c.Assert(kdTree.root.hypercube, Equals, &cube72)

	c.Assert(kdTree.root.left.dimension, Equals, 1)
	c.Assert(kdTree.root.left.hypercube, Equals, &cube54)

	c.Assert(kdTree.root.left.left.dimension, Equals, 0)
	c.Assert(kdTree.root.left.left.hypercube, Equals, &cube23)
	c.Assert(kdTree.root.left.left.left, Equals, emptyNode)
	c.Assert(kdTree.root.left.left.right, Equals, emptyNode)

	c.Assert(kdTree.root.left.right.dimension, Equals, 0)
	c.Assert(kdTree.root.left.right.hypercube, Equals, &cube47)
	c.Assert(kdTree.root.left.right.left, Equals, emptyNode)
	c.Assert(kdTree.root.left.right.right, Equals, emptyNode)

	c.Assert(kdTree.root.right.dimension, Equals, 1)
	c.Assert(kdTree.root.right.hypercube, Equals, &cube96)

	c.Assert(kdTree.root.right.left.dimension, Equals, 0)
	c.Assert(kdTree.root.right.left.hypercube, Equals, &cube81)
	c.Assert(kdTree.root.right.left.left, Equals, emptyNode)
	c.Assert(kdTree.root.right.left.right, Equals, emptyNode)

	c.Assert(kdTree.root.right.right, Equals, emptyNode)

	// Invalid parameters.
	c.Assert(func() { newHypecubeKdTree(nil) }, Panics, `ERROR: newHypecubeKdTree called with no hypercubes`)
	c.Assert(func() { newHypecubeKdTree([]*specimenHypercube{}) }, Panics, `ERROR: newHypecubeKdTree called with no hypercubes`)
	c.Assert(func() { newHypecubeKdTree([]*specimenHypercube{&specimenHypercube{dimensions: nil}}) }, Panics, `ERROR: newHypecubeKdTree called with hypercubes that have no dimensions`)
	c.Assert(func() { newHypecubeKdTree([]*specimenHypercube{&specimenHypercube{dimensions: []float64{}}}) }, Panics, `ERROR: newHypecubeKdTree called with hypercubes that have no dimensions`)
}

func (s *HypercubeKdTreeSuite) Test_MoveIndicatorBase(c *C) {

	// Moving all values.
	c.Assert(moveIndicatorBase([]float64{9.0, 9.0, 9.0}, []float64{1.0, 2.0, 3.0}, []float64{3.0, 2.0, 1.0}), DeepEquals, []float64{3.0, 2.0, 3.0})
	c.Assert(moveIndicatorBase([]float64{9.0, 9.0, 9.0}, []float64{3.0, 2.0, 1.0}, []float64{1.0, 2.0, 3.0}), DeepEquals, []float64{3.0, 2.0, 3.0})

	// For completeness, test negative values.
	c.Assert(moveIndicatorBase([]float64{9.0, 9.0, 9.0}, []float64{-1.0, -2.0, -3.0}, []float64{-3.0, -2.0, -1.0}), DeepEquals, []float64{-1.0, -2.0, -1.0})
	c.Assert(moveIndicatorBase([]float64{9.0, 9.0, 9.0}, []float64{-3.0, -2.0, -1.0}, []float64{-1.0, -2.0, -3.0}), DeepEquals, []float64{-1.0, -2.0, -1.0})

	// Move value past limit have no effect.
	c.Assert(moveIndicatorBase([]float64{9.0, 9.0, 9.0}, []float64{1.0, 2.0, 3.0}, []float64{13.0, 2.0, 1.0}), DeepEquals, []float64{1.0, 2.0, 3.0})
	c.Assert(moveIndicatorBase([]float64{9.0, 9.0, 9.0}, []float64{3.0, 2.0, 1.0}, []float64{1.0, 12.0, 3.0}), DeepEquals, []float64{3.0, 2.0, 3.0})
}

func (s *HypercubeKdTreeSuite) Test_UpdateLeftViable(c *C) {

	// A dimension turns off.
	c.Assert(updateLeftViable(0, []bool{true, true, true}, []float64{1.0, 2.0, 3.0}, []float64{3.0, 1.0, 1.0}), DeepEquals, []bool{true, true, true})
	c.Assert(updateLeftViable(0, []bool{true, true, true}, []float64{1.0, 2.0, 3.0}, []float64{1.0, 1.0, 4.0}), DeepEquals, []bool{false, true, true})

	// A dimension stays off.
	c.Assert(updateLeftViable(0, []bool{false, true, true}, []float64{1.0, 2.0, 3.0}, []float64{3.0, 1.0, 1.0}), DeepEquals, []bool{false, true, true})
	c.Assert(updateLeftViable(0, []bool{false, true, true}, []float64{1.0, 2.0, 3.0}, []float64{1.0, 1.0, 4.0}), DeepEquals, []bool{false, true, true})

	// Changing a different dimension.
	c.Assert(updateLeftViable(1, []bool{true, true, true}, []float64{1.0, 2.0, 3.0}, []float64{3.0, 1.0, 1.0}), DeepEquals, []bool{true, false, true})
}

func (s *HypercubeKdTreeSuite) Test_IsLeftViable(c *C) {

	// At least one viable makes it viable.
	c.Assert(isLeftViable([]bool{true, true, true}), Equals, true)
	c.Assert(isLeftViable([]bool{true, false, true}), Equals, true)
	c.Assert(isLeftViable([]bool{true, false, false}), Equals, true)
	c.Assert(isLeftViable([]bool{false, false, false}), Equals, false)
}

func (s *HypercubeKdTreeSuite) Test_KdCompareHypercubes(c *C) {

	// Determine whether it is meaningful for the cube to search down the left or right branches of the k-d tree.
	var isLeftSearchable, isRightSearchable bool

	// We need two hypercubes to compare.
	var searchingHypercube *specimenHypercube
	var nodeHypercube *specimenHypercube

	// The nodes keep track of the dimensions of cubes going out of them.
	var maximumLeft []float64 = []float64{0.0, 0.0}
	var minimumRight []float64 = []float64{100.0, 100.0}

	// Start with two hypercubes that overlap but no domination.
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	isLeftSearchable, isRightSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft, minimumRight)
	c.Check(isLeftSearchable, Equals, true)  // Still searchable.
	c.Check(isRightSearchable, Equals, true) // Still searchable.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{1.0, 1.0}})
	c.Assert(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.0, 1.0}, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})

	// Reverse the inputs.
	searchingHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	isLeftSearchable, isRightSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft, minimumRight)
	c.Check(isLeftSearchable, Equals, true)  // Still searchable.
	c.Check(isRightSearchable, Equals, true) // Still searchable.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.0, 1.0}})
	c.Assert(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{1.0, 1.0}, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})

	// Left branch not searchable, indicator base too high.
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{2.0, 1.0}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{3.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	isLeftSearchable, isRightSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, []float64{2.0, 1.0}, minimumRight)
	c.Check(isLeftSearchable, Equals, false) // Nothing to the left that will move our indicator base higher.
	c.Check(isRightSearchable, Equals, true) // Still searchable.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{2.0, 1.0}})
	c.Assert(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{3.0, 2.0}, indicatorBase: []float64{1.0, 1.0}, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})

	// One cube dominates the other.
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{0.5, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	isLeftSearchable, isRightSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft, minimumRight)
	c.Check(isLeftSearchable, Equals, false) // Only smaller values would be to the left, the indicator base won't get larger.
	c.Check(isRightSearchable, Equals, true) // Still searchable.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.5, 1.0}})
	c.Assert(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{0.5, 1.0}, isDominated: true, indicatorBase: nil, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})

	// Reverse the inputs.
	searchingHypercube = &specimenHypercube{dimensions: []float64{0.5, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	isLeftSearchable, isRightSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft, minimumRight)
	c.Check(isLeftSearchable, Equals, false)  // Searching cube dominated, nothing more to search.
	c.Check(isRightSearchable, Equals, false) // Searching cube dominated, nothing more to search.
	// The searching cube has new indicator base.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{0.5, 1.0}, isDominated: true, indicatorBase: nil})
	// The node cube has new indicator base, but also remembers the comparision.
	c.Assert(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.5, 1.0}, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})

	// One cube already dominated.
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, isDominated: true, indicatorBase: nil}
	nodeHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	isLeftSearchable, isRightSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft, minimumRight)
	c.Check(isLeftSearchable, Equals, false)  // Searching cube dominated, nothing more can be done for hypervolume indicator.
	c.Check(isRightSearchable, Equals, false) // Searching cube dominated, nothing more can be done for hypervolume indicator.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, isDominated: true, indicatorBase: nil})
	c.Assert(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.0, 1.0}})

	// Reverse the inputs.
	searchingHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, isDominated: true, indicatorBase: nil}
	isLeftSearchable, isRightSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft, minimumRight)
	c.Check(isLeftSearchable, Equals, false) // The node cube is dominated, all cubes to the left are dominated too, only cubes to the right have more to share.
	c.Check(isRightSearchable, Equals, true) // Still searchable.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.0, 1.0}})
	c.Assert(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, isDominated: true, indicatorBase: nil})

	// The cubes are absolutely equal.
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	isLeftSearchable, isRightSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft, minimumRight)
	c.Check(isLeftSearchable, Equals, true)  // Still searchable, smaller valueas to the left but we need at least one to give us an indicator.
	c.Check(isRightSearchable, Equals, true) // Still searchable.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}})
	c.Assert(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, isDominated: true, indicatorBase: nil, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})

	// The cubes are the same cube.
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = searchingHypercube
	isLeftSearchable, isRightSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft, minimumRight)
	c.Check(isLeftSearchable, Equals, true)  // Still searchable, only comparing against cubes that are not us can tell us there is no more to find.
	c.Check(isRightSearchable, Equals, true) // Still searchable.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}})
	c.Assert(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}})

	// Compare against a prior compare.
	nodeHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}, comparedWith: map[*specimenHypercube]bool{nodeHypercube: true}}
	isLeftSearchable, isRightSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft, minimumRight)
	c.Check(isLeftSearchable, Equals, true)  // Still searchable.
	c.Check(isRightSearchable, Equals, true) // Still searchable.
	// No changes made to either hypercube.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}, comparedWith: map[*specimenHypercube]bool{nodeHypercube: true}})
	c.Assert(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}})

	// Compare against a prior compare, left branch not searchable.
	nodeHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{2.0, 1.0}, comparedWith: map[*specimenHypercube]bool{nodeHypercube: true}}
	isLeftSearchable, isRightSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, []float64{2.0, 1.0}, minimumRight)
	c.Check(isLeftSearchable, Equals, false) // Nothing to the left that will move our indicator base higher.
	c.Check(isRightSearchable, Equals, true) // Still searchable.
	// No changes made to either hypercube.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}, comparedWith: map[*specimenHypercube]bool{nodeHypercube: true}})
	c.Assert(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}})

	// Compare against a prior compare, right branch not searchable.
	nodeHypercube = &specimenHypercube{dimensions: []float64{2.0, 3.0}, indicatorBase: []float64{0.0, 0.0}}
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}, comparedWith: map[*specimenHypercube]bool{nodeHypercube: true}}
	isLeftSearchable, isRightSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft, []float64{1.0, 2.0})
	c.Check(isLeftSearchable, Equals, true)  // Still searchable.
	c.Check(isRightSearchable, Equals, true) // Only things to the right are dominating cubes, but there may be no cubes at all.
	// No changes made to either hypercube.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}, comparedWith: map[*specimenHypercube]bool{nodeHypercube: true}})
	c.Assert(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 3.0}, indicatorBase: []float64{0.0, 0.0}})

	// Compare against a prior compare, neither branch searchable.
	nodeHypercube = &specimenHypercube{dimensions: []float64{2.0, 3.0}, indicatorBase: []float64{0.0, 0.0}}
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{2.0, 3.0}, comparedWith: map[*specimenHypercube]bool{nodeHypercube: true}}
	isLeftSearchable, isRightSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, []float64{2.0, 1.0}, []float64{1.0, 2.0})
	c.Check(isLeftSearchable, Equals, false) // Nothing to the left that will move our indicator base higher.
	c.Check(isRightSearchable, Equals, true) // Only things to the right are dominating cubes, but there may be no cubes at all.
	// No changes made to either hypercube.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}, comparedWith: map[*specimenHypercube]bool{nodeHypercube: true}})
	c.Assert(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 3.0}, indicatorBase: []float64{0.0, 0.0}})

}
