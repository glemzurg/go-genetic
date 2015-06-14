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
