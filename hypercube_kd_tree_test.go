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
	c.Assert(kdTree.root.maximumLeft, Equals, nil)

	c.Assert(kdTree.root.left.dimension, Equals, 1)
	c.Assert(kdTree.root.left.hypercube, Equals, &cube54)
	c.Assert(kdTree.root.left.maximumLeft, Equals, []float64{5.0, 3.0})

	c.Assert(kdTree.root.left.left.dimension, Equals, 0)
	c.Assert(kdTree.root.left.left.hypercube, Equals, &cube23)
	c.Assert(kdTree.root.left.left.maximumLeft, Equals, nil)
	c.Assert(kdTree.root.left.left.left, Equals, emptyNode)
	c.Assert(kdTree.root.left.left.right, Equals, emptyNode)

	c.Assert(kdTree.root.left.right.dimension, Equals, 0)
	c.Assert(kdTree.root.left.right.hypercube, Equals, &cube47)
	c.Assert(kdTree.root.left.right.maximumLeft, Equals, nil)
	c.Assert(kdTree.root.left.right.left, Equals, emptyNode)
	c.Assert(kdTree.root.left.right.right, Equals, emptyNode)

	c.Assert(kdTree.root.right.dimension, Equals, 1)
	c.Assert(kdTree.root.right.hypercube, Equals, &cube96)
	c.Assert(kdTree.root.right.maximumLeft, Equals, []float64{9.0, 1.0})

	c.Assert(kdTree.root.right.left.dimension, Equals, 0)
	c.Assert(kdTree.root.right.left.hypercube, Equals, &cube81)
	c.Assert(kdTree.root.right.left.maximumLeft, Equals, nil)
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

	// Determine whether it is meaningful for the cube to search down the left branch of the k-d tree.
	var isLeftSearchable bool

	// We need two hypercubes to compare.
	var searchingHypercube *specimenHypercube // The hypercube walking the tree.
	var nodeHypercube *specimenHypercube      // The hypercube in this node.

	// The nodes keep track of the dimensions of cubes going down the left branch.
	var maximumLeft []float64 = nil // Assume we don't know what's to the left so must explore.

	// Start with two hypercubes that overlap but no domination.
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, true) // Don't know what's down left branch, need to search.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 1.0}})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.0, 0.0}, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})

	// Reverse the inputs.
	searchingHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, true) // Don't know what's down left branch, need to search.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.0, 0.0}})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 1.0}, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})

	// Two hypercubes with indicator bases and overlap but no domination.
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.5, 1.5}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.5, 0.5}}
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, true) // Don't know what's down left branch, need to search.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.5, 1.5}})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.5, 0.5}, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})

	// Reverse the inputs.
	searchingHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.5, 0.5}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.5, 1.5}}
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, true) // Don't know what's down left branch, need to search.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.5, 0.5}})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.5, 1.5}, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})

	// One cube dominates the other.
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{0.5, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, true) // Don't know what's down left branch, need to search.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.5, 1.0}})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{0.5, 1.0}, isDominated: true, indicatorBase: nil})

	// Reverse the inputs.
	searchingHypercube = &specimenHypercube{dimensions: []float64{0.5, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, false) // Searching cube dominated, nothing more to search.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{0.5, 1.0}, isDominated: true, indicatorBase: nil})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.5, 1.0}, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})

	// One cube already dominated.
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, isDominated: true, indicatorBase: nil}
	nodeHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, false) // Searching cube dominated, nothing more can be done for hypervolume indicator.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, isDominated: true, indicatorBase: nil})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.0, 0.0}, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})

	// Reverse the inputs.
	searchingHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, isDominated: true, indicatorBase: nil}
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, true) // Don't know what's down left branch, need to search.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.0, 0.0}})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, isDominated: true, indicatorBase: nil})

	// Both cubes dominated.
	searchingHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, isDominated: true, indicatorBase: nil}
	nodeHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, isDominated: true, indicatorBase: nil}
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, false) // Searching cube dominated, nothing more to search.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, isDominated: true, indicatorBase: nil})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, isDominated: true, indicatorBase: nil})

	// The cubes are absolutely equal.
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, true) // Don't know what's down left branch, need to search.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, isDominated: true, indicatorBase: nil})

	// The cubes are the same cube.
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}}
	nodeHypercube = searchingHypercube
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, true) // When we don't compare cubes, assume searchable.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}})

	// Compare against a prior compare.
	nodeHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}}
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}, comparedWith: map[*specimenHypercube]bool{nodeHypercube: true}}
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, true) // When we don't compare cubes, assume searchable.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.0, 0.0}, comparedWith: map[*specimenHypercube]bool{nodeHypercube: true}})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.0, 0.0}})

	// Maximim left that is above the searching indicator base in at least one dimension.
	maximumLeft = []float64{0.75, 0.75}
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.5, 0.5}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.5, 0.5}}
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, true) // There is something down the left branch worth searching.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.5, 1.0}})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.0, 0.5}, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})

	// Maximim left that is equal to the searching idicator base in all dimensions.
	maximumLeft = []float64{0.5, 1.0}
	searchingHypercube = &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.5, 0.5}}
	nodeHypercube = &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{0.5, 0.5}}
	isLeftSearchable = kdCompareHypercubes(searchingHypercube, nodeHypercube, maximumLeft)
	c.Check(isLeftSearchable, Equals, false) // Nothing more for the searching hyper cube to fine regarding its indicator base.
	c.Check(searchingHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{1.0, 2.0}, indicatorBase: []float64{0.5, 1.0}})
	c.Assert(nodeHypercube, DeepEquals, &specimenHypercube{dimensions: []float64{2.0, 1.0}, indicatorBase: []float64{1.0, 0.5}, comparedWith: map[*specimenHypercube]bool{searchingHypercube: true}})
}

func (s *HypercubeKdTreeSuite) Test_CalculateNewMaximumLeft(c *C) {

	// First build up the maximum left.
	c.Check(calculateNewMaximumLeft(nil, 3, 0, 1.0), DeepEquals, []float64{1.0})
	c.Check(calculateNewMaximumLeft([]float64{1.0}, 3, 1, 2.0), DeepEquals, []float64{1.0, 2.0})
	c.Check(calculateNewMaximumLeft([]float64{1.0, 2.0}, 3, 2, 3.0), DeepEquals, []float64{1.0, 2.0, 3.0})
	c.Check(calculateNewMaximumLeft([]float64{1.0, 2.0, 3.0}, 3, 0, 0.5), DeepEquals, []float64{0.5, 2.0, 3.0})
	c.Check(calculateNewMaximumLeft([]float64{0.5, 2.0, 3.0}, 3, 1, 1.5), DeepEquals, []float64{0.5, 1.5, 3.0})
	c.Check(calculateNewMaximumLeft([]float64{0.5, 1.5, 3.0}, 3, 2, 2.5), DeepEquals, []float64{0.5, 1.5, 2.5})

	// Setting to the same value is ok.
	c.Check(calculateNewMaximumLeft([]float64{1.0, 2.0, 3.0}, 3, 0, 1.0), DeepEquals, []float64{1.0, 2.0, 3.0})
	c.Check(calculateNewMaximumLeft([]float64{1.0, 2.0, 3.0}, 3, 1, 2.0), DeepEquals, []float64{1.0, 2.0, 3.0})
	c.Check(calculateNewMaximumLeft([]float64{1.0, 2.0, 3.0}, 3, 2, 3.0), DeepEquals, []float64{1.0, 2.0, 3.0})

	// Invalid params.
	c.Check(func() { calculateNewMaximumLeft([]float64{1.0}, 3, 2, 3.0) }, Panics, "ASSERT: invalid maximum left append: 2 [1.000000]")
	c.Check(func() { calculateNewMaximumLeft([]float64{1.0, 2.0, 3.0}, 3, 1, 2.1) }, Panics, "ASSERT: invalid maximum left update: 2.100000 [1 2 3]")
}
