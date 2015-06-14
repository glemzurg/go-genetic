package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type SpecimenHypercubeSuite struct{}

var _ = Suite(&SpecimenHypercubeSuite{})

// Add the tests.

func (s *SpecimenHypercubeSuite) Test_SpecimenHypercubeDimensions(c *C) {
	var dimensions []float64
	var volume float64

	// Some simple configurations.
	var referencePoint []float64 = []float64{-1.0, -2.0, 3.0}
	var isMaximize []bool = []bool{true, true, true}
	var weights []float64 = []float64{3.0, 2.0, 1.0}

	// A simple point.
	dimensions, volume = specimenHypercubeDimensions([]float64{4.0, 4.0, 4.0}, referencePoint, isMaximize, weights)
	c.Check(dimensions, DeepEquals, []float64{5.0 * 3.0, 6.0 * 2.0, 1.0 * 1.0})
	c.Assert(volume, Equals, (5.0*3.0)*(6.0*2.0)*(1.0*1.0))

	// A dimension at the reference point.
	dimensions, volume = specimenHypercubeDimensions([]float64{4.0, -2.0, 4.0}, referencePoint, isMaximize, weights)
	c.Check(dimensions, DeepEquals, []float64{5.0 * 3.0, 0.0, 1.0 * 1.0})
	c.Assert(volume, Equals, 0.0) // Any dimension at or below the reference point will have no volume.

	// A dimension below the reference point.
	dimensions, volume = specimenHypercubeDimensions([]float64{4.0, -3.0, 4.0}, referencePoint, isMaximize, weights)
	c.Check(dimensions, DeepEquals, []float64{5.0 * 3.0, 0.0, 1.0 * 1.0})
	c.Assert(volume, Equals, 0.0) // Any dimension at or below the reference point will have no volume.

	// Let's make one of the dimensions a minimizing dimension.
	isMaximize = []bool{true, false, true}

	// A simple point.
	dimensions, volume = specimenHypercubeDimensions([]float64{4.0, -4.0, 4.0}, referencePoint, isMaximize, weights)
	c.Check(dimensions, DeepEquals, []float64{5.0 * 3.0, 2.0 * 2.0, 1.0 * 1.0})
	c.Assert(volume, Equals, (5.0*3.0)*(2.0*2.0)*(1.0*1.0))

	// A dimension at the reference point.
	dimensions, volume = specimenHypercubeDimensions([]float64{4.0, -2.0, 4.0}, referencePoint, isMaximize, weights)
	c.Check(dimensions, DeepEquals, []float64{5.0 * 3.0, 0.0, 1.0 * 1.0})
	c.Assert(volume, Equals, 0.0) // Any dimension at or above the reference point will have no volume.

	// A dimension above the minimizing reference point.
	dimensions, volume = specimenHypercubeDimensions([]float64{4.0, -1.0, 4.0}, referencePoint, isMaximize, weights)
	c.Check(dimensions, DeepEquals, []float64{5.0 * 3.0, 0.0, 1.0 * 1.0})
	c.Assert(volume, Equals, 0.0) // Any dimension at or above the reference point will have no volume.

	// Invalid parameters.
	c.Assert(func() { specimenHypercubeDimensions([]float64{4.0, -1.0}, referencePoint, isMaximize, weights) }, Panics, `specimenHypercubeDimensions expects 3 dimensions, but multi-outcome has 2 dimensions`)
}

func (s *SpecimenHypercubeSuite) Test_NewSpecimenHypercube(c *C) {

	// Some simple configurations.
	var referencePoint []float64 = []float64{-1.0, -2.0, 3.0}
	var isMaximize []bool = []bool{true, true, true}
	var weights []float64 = []float64{3.0, 2.0, 1.0}

	// The specimen in question.
	var specimen Specimen = Specimen{
		Outcomes: []float64{4.0, 4.0, 4.0},
	}

	// The expected hypercube.
	var expectedHypercube specimenHypercube = specimenHypercube{
		dimensions: []float64{5.0 * 3.0, 6.0 * 2.0, 1.0 * 1.0},
		volume:     (5.0 * 3.0) * (6.0 * 2.0) * (1.0 * 1.0),
		specimen:   specimen,
	}

	// Create a new hypercube.
	var hypercube specimenHypercube = newSpecimenHypercube(specimen, referencePoint, isMaximize, weights)
	c.Assert(hypercube, DeepEquals, expectedHypercube)
}

func (s *SpecimenHypercubeSuite) Test_isDominatedBy(c *C) {

	// Some specimens.
	var hypercubeA specimenHypercube = specimenHypercube{
		dimensions: []float64{4.0, 4.0, 4.0},
	}
	var hypercubeB specimenHypercube = specimenHypercube{
		dimensions: []float64{4.0, 4.0, 4.0},
	}
	var hypercubeC specimenHypercube = specimenHypercube{
		dimensions: []float64{3.0, 3.0, 3.0},
	}
	var hypercubeD specimenHypercube = specimenHypercube{
		dimensions: []float64{5.0, 3.0, 3.0},
	}

	c.Check(hypercubeA.isDominatedBy(&hypercubeB), Equals, true)
	c.Check(hypercubeA.isDominatedBy(&hypercubeC), Equals, false)
	c.Check(hypercubeA.isDominatedBy(&hypercubeD), Equals, false)

	c.Check(hypercubeB.isDominatedBy(&hypercubeA), Equals, true)
	c.Check(hypercubeB.isDominatedBy(&hypercubeC), Equals, false)
	c.Check(hypercubeB.isDominatedBy(&hypercubeD), Equals, false)

	c.Check(hypercubeC.isDominatedBy(&hypercubeA), Equals, true)
	c.Check(hypercubeC.isDominatedBy(&hypercubeB), Equals, true)
	c.Check(hypercubeC.isDominatedBy(&hypercubeD), Equals, true)

	c.Check(hypercubeD.isDominatedBy(&hypercubeA), Equals, false)
	c.Check(hypercubeD.isDominatedBy(&hypercubeB), Equals, false)
	c.Check(hypercubeD.isDominatedBy(&hypercubeC), Equals, false)
}
