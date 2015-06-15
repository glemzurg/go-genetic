package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type HypervolumeIndicatorSuite struct{}

var _ = Suite(&HypervolumeIndicatorSuite{})

// Add the tests.

func (s *HypervolumeIndicatorSuite) Test_CalculateHypervolumeIndicators(c *C) {

	// Test the two dimentional problem from this page: http://esa.github.io/pygmo/tutorials/getting_started_with_hyper_volumes.html
	// In our code we have been normalized to always be positive dimensions with zero as the base. In the examine the reference point
	// is 2.0, 2.0 and the cubes are stretching towards 0.0, 0.0.
	var cube2010 specimenHypercube = specimenHypercube{
		dimensions: []float64{2.0, 1.0},
	}
	var cube1515 specimenHypercube = specimenHypercube{
		dimensions: []float64{1.5, 1.5},
	}
	var cube1020 specimenHypercube = specimenHypercube{
		dimensions: []float64{1.0, 2.0},
	}
	var cube0513 specimenHypercube = specimenHypercube{
		dimensions: []float64{0.5, 1.3}, // 1.25 in the example, but dominated either way.
	}

	// What do we expect these cubes to compute to?
	var expectedCube2010 specimenHypercube = specimenHypercube{
		dimensions:    []float64{2.0, 1.0},
		isDominated:   false,
		indicator:     0.5,
		indicatorBase: []float64{1.5, 0.0},
	}
	var expectedCube1515 specimenHypercube = specimenHypercube{
		dimensions:    []float64{1.5, 1.5},
		isDominated:   false,
		indicator:     0.25,
		indicatorBase: []float64{1.0, 1.0},
	}
	var expectedCube1020 specimenHypercube = specimenHypercube{
		dimensions:    []float64{1.0, 2.0},
		isDominated:   false,
		indicator:     0.5,
		indicatorBase: []float64{0.0, 1.5},
	}
	var expectedCube0513 specimenHypercube = specimenHypercube{
		dimensions:    []float64{0.5, 1.3}, // 1.25 in the example, but dominated either way.
		isDominated:   true,                // Dominated.
		indicator:     0.0,                 // Dominated.
		indicatorBase: nil,                 // Dominated.
	}

	// Make the list of pointers.
	var hypercubes []*specimenHypercube = []*specimenHypercube{
		&cube2010,
		&cube1515,
		&cube1020,
		&cube0513,
	}

	// Calculat the hypervolume indicators.
	calculateHypervolumeIndicators(hypercubes)

	// Did we get what we expected?
	c.Check(cube2010, DeepEquals, expectedCube2010)
	c.Check(cube1515, DeepEquals, expectedCube1515)
	c.Check(cube1020, DeepEquals, expectedCube1020)
	c.Check(cube0513, DeepEquals, expectedCube0513)
}

func (s *HypervolumeIndicatorSuite) Test_CalculateHypervolume(c *C) {

	// Calculate some hypervolumes. Base dimensions are always less than limit dimensions.
	c.Check(calculateHypervolume([]float64{3.0, 4.0, 5.0}, []float64{1.0, 2.0, 3.0}), Equals, 8.0)

	// Invalid parameters.
	c.Assert(func() { calculateHypervolume([]float64{3.0, 4.0, 5.0}, []float64{3.0, 2.0, 3.0}) }, Panics, `ASSERT: calculateHypervolume base can never be equal to or greater than limit in any dimension`)
	c.Assert(func() { calculateHypervolume([]float64{3.0, 4.0, 5.0}, []float64{1.0, 5.0, 3.0}) }, Panics, `ASSERT: calculateHypervolume base can never be equal to or greater than limit in any dimension`)
}
