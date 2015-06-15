package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type HypervolumeIndicatorSuite struct{}

var _ = Suite(&HypervolumeIndicatorSuite{})

// Add the tests.

func (s *HypervolumeIndicatorSuite) Test_HypervolumeIndicator(c *C) {

	// // Start with some hypercube.
	// var hypercube []hypercubeDimension = []hypercubeDimension{
	// 	newHypercubeDimension(-5.0, true, 1.0),
	// 	newHypercubeDimension(5.0, false, 1.0),
	// 	newHypercubeDimension(-5.0, true, 1.0),
	// }
	// var expectedHypercube []hypercubeDimension = []hypercubeDimension{
	// 	hypercubeDimension{first: -5.0, second: -5.0, base: -5.0, isMaximize: true, weight: 1.0},
	// 	hypercubeDimension{first: 5.0, second: 5.0, base: 5.0, isMaximize: false, weight: 1.0},
	// 	hypercubeDimension{first: -5.0, second: -5.0, base: -5.0, isMaximize: true, weight: 1.0},
	// }
	// c.Assert(hypercube, DeepEquals, expectedHypercube)
	//
	// // Stretch all dimensions.
	// expectedHypercube = []hypercubeDimension{
	// 	hypercubeDimension{first: -4.0, second: -5.0, base: -5.0, isMaximize: true, weight: 1.0},
	// 	hypercubeDimension{first: 4.0, second: 5.0, base: 5.0, isMaximize: false, weight: 1.0},
	// 	hypercubeDimension{first: -3.0, second: -5.0, base: -5.0, isMaximize: true, weight: 1.0},
	// }
	// hypercube = stretchDimensions(hypercube, []float64{-4.0, 4.0, -3.0})
	// c.Assert(hypercube, DeepEquals, expectedHypercube)
	//
	// // Stretch all dimensions.
	// expectedHypercube = []hypercubeDimension{
	// 	hypercubeDimension{first: 100.0, second: -4.0, base: -5.0, isMaximize: true, weight: 1.0},
	// 	hypercubeDimension{first: -100.0, second: 4.0, base: 5.0, isMaximize: false, weight: 1.0},
	// 	hypercubeDimension{first: 99.0, second: -3.0, base: -5.0, isMaximize: true, weight: 1.0},
	// }
	// hypercube = stretchDimensions(hypercube, []float64{100.0, -100.0, 99.0})
	// c.Assert(hypercube, DeepEquals, expectedHypercube)
	//
	// // Invalid parameters.
	// c.Assert(func() { stretchDimensions(hypercube, []float64{100.0, -100.0}) }, Panics, `stretchDimensions expects 3 dimensions, but outcomes have 2 dimensions`)

}

func (s *HypervolumeIndicatorSuite) Test_CalculateHypervolume(c *C) {

	// Calculate some hypervolumes. Base dimensions are always less than limit dimensions.
	c.Check(calculateHypervolume([]float64{3.0, 4.0, 5.0}, []float64{1.0, 2.0, 3.0}), Equals, 8.0)

	// Invalid parameters.
	c.Assert(func() { calculateHypervolume([]float64{3.0, 4.0, 5.0}, []float64{3.0, 2.0, 3.0}) }, Panics, `ASSERT: calculateHypervolume base can never be equal to or greater than limit in any dimension`)
	c.Assert(func() { calculateHypervolume([]float64{3.0, 4.0, 5.0}, []float64{1.0, 5.0, 3.0}) }, Panics, `ASSERT: calculateHypervolume base can never be equal to or greater than limit in any dimension`)
}
