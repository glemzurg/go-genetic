package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type HypervolumeIndicatorSuite struct{}

var _ = Suite(&HypervolumeIndicatorSuite{})

// Add the tests.

func (s *HypervolumeIndicatorSuite) Test_DimensionMax_Stretch(c *C) {

	var dimension hypercubeDimension

	// Simple dimensionimize.
	dimension = newHypercubeDimension(-5.0, true, 1.0)
	c.Assert(dimension, Equals, hypercubeDimension{first: -5.0, second: -5.0, isMaximize: true, weight: 1.0})
	dimension.stretch(-6.0) // Don't stretch, less than base.
	c.Assert(dimension, Equals, hypercubeDimension{first: -5.0, second: -5.0, isMaximize: true, weight: 1.0})
	dimension.stretch(-5.0) // Don't stretch, equal to base.
	c.Assert(dimension, Equals, hypercubeDimension{first: -5.0, second: -5.0, isMaximize: true, weight: 1.0})
	dimension.stretch(-4.0) // Stretch dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: -4.0, second: -5.0, isMaximize: true, weight: 1.0})
	dimension.stretch(11.0) // Stretch dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: 11.0, second: -4.0, isMaximize: true, weight: 1.0})
	dimension.stretch(12.0) // Stretch dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: 12.0, second: 11.0, isMaximize: true, weight: 1.0})
	dimension.stretch(11.5) // Stretch second best dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: 12.0, second: 11.5, isMaximize: true, weight: 1.0})
	dimension.stretch(11.5) // Doesn't stretch second best dimension. Equal.
	c.Assert(dimension, Equals, hypercubeDimension{first: 12.0, second: 11.5, isMaximize: true, weight: 1.0})
	dimension.stretch(11.0) // Doesn't stretch dimension. Less than.
	c.Assert(dimension, Equals, hypercubeDimension{first: 12.0, second: 11.5, isMaximize: true, weight: 1.0})
	dimension.stretch(12.0) // Stretch second best dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: 12.0, second: 12.0, isMaximize: true, weight: 1.0})

	// Simple minimize.
	dimension = newHypercubeDimension(5.0, false, 1.0)
	c.Assert(dimension, Equals, hypercubeDimension{first: 5.0, second: 5.0, isMaximize: false, weight: 1.0})
	dimension.stretch(6.0) // Don't stretch, greater than base.
	c.Assert(dimension, Equals, hypercubeDimension{first: 5.0, second: 5.0, isMaximize: false, weight: 1.0})
	dimension.stretch(5.0) // Don't stretch, equal to base.
	c.Assert(dimension, Equals, hypercubeDimension{first: 5.0, second: 5.0, isMaximize: false, weight: 1.0})
	dimension.stretch(4.0) // Stretch dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: 4.0, second: 5.0, isMaximize: false, weight: 1.0})
	dimension.stretch(-11.0) // Stretch dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: -11.0, second: 4.0, isMaximize: false, weight: 1.0})
	dimension.stretch(-12.0) // Stretch dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: -12.0, second: -11.0, isMaximize: false, weight: 1.0})
	dimension.stretch(-11.5) // Stretch second best dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: -12.0, second: -11.5, isMaximize: false, weight: 1.0})
	dimension.stretch(-11.5) // Doesn't stretch second best dimension. Equal.
	c.Assert(dimension, Equals, hypercubeDimension{first: -12.0, second: -11.5, isMaximize: false, weight: 1.0})
	dimension.stretch(-11.0) // Doesn't stretch dimension. Less than.
	c.Assert(dimension, Equals, hypercubeDimension{first: -12.0, second: -11.5, isMaximize: false, weight: 1.0})
	dimension.stretch(-12.0) // Stretch second best dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: -12.0, second: -12.0, isMaximize: false, weight: 1.0})
}

func (s *HypervolumeIndicatorSuite) Test_StrechDimensions(c *C) {

	// Start with some hypercube.
	var hypercube []hypercubeDimension = []hypercubeDimension{
		newHypercubeDimension(-5.0, true, 1.0),
		newHypercubeDimension(5.0, false, 1.0),
		newHypercubeDimension(-5.0, true, 1.0),
	}
	var expectedHypercube []hypercubeDimension = []hypercubeDimension{
		hypercubeDimension{first: -5.0, second: -5.0, isMaximize: true, weight: 1.0},
		hypercubeDimension{first: 5.0, second: 5.0, isMaximize: false, weight: 1.0},
		hypercubeDimension{first: -5.0, second: -5.0, isMaximize: true, weight: 1.0},
	}
	c.Assert(hypercube, DeepEquals, expectedHypercube)

	// Stretch all dimensions.
	expectedHypercube = []hypercubeDimension{
		hypercubeDimension{first: -4.0, second: -5.0, isMaximize: true, weight: 1.0},
		hypercubeDimension{first: 4.0, second: 5.0, isMaximize: false, weight: 1.0},
		hypercubeDimension{first: -3.0, second: -5.0, isMaximize: true, weight: 1.0},
	}
	hypercube = stretchDimensions(hypercube, []float64{-4.0, 4.0, -3.0})
	c.Assert(hypercube, DeepEquals, expectedHypercube)

	// Stretch all dimensions.
	expectedHypercube = []hypercubeDimension{
		hypercubeDimension{first: 100.0, second: -4.0, isMaximize: true, weight: 1.0},
		hypercubeDimension{first: -100.0, second: 4.0, isMaximize: false, weight: 1.0},
		hypercubeDimension{first: 99.0, second: -3.0, isMaximize: true, weight: 1.0},
	}
	hypercube = stretchDimensions(hypercube, []float64{100.0, -100.0, 99.0})
	c.Assert(hypercube, DeepEquals, expectedHypercube)

	// Invalid parameters.
	c.Assert(func() { stretchDimensions(hypercube, []float64{100.0, -100.0}) }, Panics, `stretchDimensions expects 3 dimensions, but outcomes have 2 dimensions`)

}
