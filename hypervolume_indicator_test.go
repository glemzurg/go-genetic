package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type HypervolumeIndicatorSuite struct{}

var _ = Suite(&HypervolumeIndicatorSuite{})

// Add the tests.

func (s *HypervolumeIndicatorSuite) Test_DimensionMax_Stretch(c *C) {

	var max dimensionMax

	// Simple maximize.
	max = newDimensionMax(-5.0, true)
	c.Assert(max, Equals, dimensionMax{first: -5.0, second: -5.0, isMaximize: true})
	max.stretch(-6.0) // Don't stretch, less than base.
	c.Assert(max, Equals, dimensionMax{first: -5.0, second: -5.0, isMaximize: true})
	max.stretch(-5.0) // Don't stretch, equal to base.
	c.Assert(max, Equals, dimensionMax{first: -5.0, second: -5.0, isMaximize: true})
	max.stretch(-4.0) // Stretch dimension.
	c.Assert(max, Equals, dimensionMax{first: -4.0, second: -5.0, isMaximize: true})
	max.stretch(11.0) // Stretch dimension.
	c.Assert(max, Equals, dimensionMax{first: 11.0, second: -4.0, isMaximize: true})
	max.stretch(12.0) // Stretch dimension.
	c.Assert(max, Equals, dimensionMax{first: 12.0, second: 11.0, isMaximize: true})
	max.stretch(11.5) // Stretch second best dimension.
	c.Assert(max, Equals, dimensionMax{first: 12.0, second: 11.5, isMaximize: true})
	max.stretch(11.5) // Doesn't stretch second best dimension. Equal.
	c.Assert(max, Equals, dimensionMax{first: 12.0, second: 11.5, isMaximize: true})
	max.stretch(11.0) // Doesn't stretch dimension. Less than.
	c.Assert(max, Equals, dimensionMax{first: 12.0, second: 11.5, isMaximize: true})
	max.stretch(12.0) // Stretch second best dimension.
	c.Assert(max, Equals, dimensionMax{first: 12.0, second: 12.0, isMaximize: true})

	// Simple minimize.
	max = newDimensionMax(5.0, false)
	c.Assert(max, Equals, dimensionMax{first: 5.0, second: 5.0, isMaximize: false})
	max.stretch(6.0) // Don't stretch, greater than base.
	c.Assert(max, Equals, dimensionMax{first: 5.0, second: 5.0, isMaximize: false})
	max.stretch(5.0) // Don't stretch, equal to base.
	c.Assert(max, Equals, dimensionMax{first: 5.0, second: 5.0, isMaximize: false})
	max.stretch(4.0) // Stretch dimension.
	c.Assert(max, Equals, dimensionMax{first: 4.0, second: 5.0, isMaximize: false})
	max.stretch(-11.0) // Stretch dimension.
	c.Assert(max, Equals, dimensionMax{first: -11.0, second: 4.0, isMaximize: false})
	max.stretch(-12.0) // Stretch dimension.
	c.Assert(max, Equals, dimensionMax{first: -12.0, second: -11.0, isMaximize: false})
	max.stretch(-11.5) // Stretch second best dimension.
	c.Assert(max, Equals, dimensionMax{first: -12.0, second: -11.5, isMaximize: false})
	max.stretch(-11.5) // Doesn't stretch second best dimension. Equal.
	c.Assert(max, Equals, dimensionMax{first: -12.0, second: -11.5, isMaximize: false})
	max.stretch(-11.0) // Doesn't stretch dimension. Less than.
	c.Assert(max, Equals, dimensionMax{first: -12.0, second: -11.5, isMaximize: false})
	max.stretch(-12.0) // Stretch second best dimension.
	c.Assert(max, Equals, dimensionMax{first: -12.0, second: -12.0, isMaximize: false})
}

func (s *HypervolumeIndicatorSuite) Test_StrechDimensions(c *C) {

	// Start with some maxes.
	var maxes []dimensionMax = []dimensionMax{
		newDimensionMax(-5.0, true),
		newDimensionMax(5.0, false),
		newDimensionMax(-5.0, true),
	}
	var expectedMaxes []dimensionMax = []dimensionMax{
		dimensionMax{first: -5.0, second: -5.0, isMaximize: true},
		dimensionMax{first: 5.0, second: 5.0, isMaximize: false},
		dimensionMax{first: -5.0, second: -5.0, isMaximize: true},
	}
	c.Assert(maxes, DeepEquals, expectedMaxes)

	// Stretch all dimensions.
	expectedMaxes = []dimensionMax{
		dimensionMax{first: -4.0, second: -5.0, isMaximize: true},
		dimensionMax{first: 4.0, second: 5.0, isMaximize: false},
		dimensionMax{first: -3.0, second: -5.0, isMaximize: true},
	}
	maxes = stretchDimensions(maxes, []float64{-4.0, 4.0, -3.0})
	c.Assert(maxes, DeepEquals, expectedMaxes)

	// Stretch all dimensions.
	expectedMaxes = []dimensionMax{
		dimensionMax{first: 100.0, second: -4.0, isMaximize: true},
		dimensionMax{first: -100.0, second: 4.0, isMaximize: false},
		dimensionMax{first: 99.0, second: -3.0, isMaximize: true},
	}
	maxes = stretchDimensions(maxes, []float64{100.0, -100.0, 99.0})
	c.Assert(maxes, DeepEquals, expectedMaxes)

	// Invalid parameters.
	c.Assert(func() { stretchDimensions(maxes, []float64{100.0, -100.0}) }, Panics, `stretchDimensions expects 3 dimensions, but outcomes have 2 dimensions`)

}
