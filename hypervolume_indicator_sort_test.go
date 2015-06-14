package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
	"sort"
)

// Create a suite.
type HypervolumeIndicatorSortSuite struct{}

var _ = Suite(&HypervolumeIndicatorSortSuite{})

// Add the tests.

func (s *HypervolumeIndicatorSortSuite) Test_DimensionMax_Stretch(c *C) {

	var dimension hypercubeDimension

	// Simple dimensionimize.
	dimension = newHypercubeDimension(-5.0, true, 1.0)
	c.Assert(dimension, Equals, hypercubeDimension{first: -5.0, second: -5.0, base: -5.0, isMaximize: true, weight: 1.0})
	dimension.stretch(-6.0) // Don't stretch, less than base.
	c.Assert(dimension, Equals, hypercubeDimension{first: -5.0, second: -5.0, base: -5.0, isMaximize: true, weight: 1.0})
	dimension.stretch(-5.0) // Don't stretch, equal to base.
	c.Assert(dimension, Equals, hypercubeDimension{first: -5.0, second: -5.0, base: -5.0, isMaximize: true, weight: 1.0})
	dimension.stretch(-4.0) // Stretch dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: -4.0, second: -5.0, base: -5.0, isMaximize: true, weight: 1.0})
	dimension.stretch(11.0) // Stretch dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: 11.0, second: -4.0, base: -5.0, isMaximize: true, weight: 1.0})
	dimension.stretch(12.0) // Stretch dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: 12.0, second: 11.0, base: -5.0, isMaximize: true, weight: 1.0})
	dimension.stretch(11.5) // Stretch second best dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: 12.0, second: 11.5, base: -5.0, isMaximize: true, weight: 1.0})
	dimension.stretch(11.5) // Doesn't stretch second best dimension. Equal.
	c.Assert(dimension, Equals, hypercubeDimension{first: 12.0, second: 11.5, base: -5.0, isMaximize: true, weight: 1.0})
	dimension.stretch(11.0) // Doesn't stretch dimension. Less than.
	c.Assert(dimension, Equals, hypercubeDimension{first: 12.0, second: 11.5, base: -5.0, isMaximize: true, weight: 1.0})
	dimension.stretch(12.0) // Stretch second best dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: 12.0, second: 12.0, base: -5.0, isMaximize: true, weight: 1.0})

	// Simple minimize.
	dimension = newHypercubeDimension(5.0, false, 1.0)
	c.Assert(dimension, Equals, hypercubeDimension{first: 5.0, second: 5.0, base: 5.0, isMaximize: false, weight: 1.0})
	dimension.stretch(6.0) // Don't stretch, greater than base.
	c.Assert(dimension, Equals, hypercubeDimension{first: 5.0, second: 5.0, base: 5.0, isMaximize: false, weight: 1.0})
	dimension.stretch(5.0) // Don't stretch, equal to base.
	c.Assert(dimension, Equals, hypercubeDimension{first: 5.0, second: 5.0, base: 5.0, isMaximize: false, weight: 1.0})
	dimension.stretch(4.0) // Stretch dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: 4.0, second: 5.0, base: 5.0, isMaximize: false, weight: 1.0})
	dimension.stretch(-11.0) // Stretch dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: -11.0, second: 4.0, base: 5.0, isMaximize: false, weight: 1.0})
	dimension.stretch(-12.0) // Stretch dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: -12.0, second: -11.0, base: 5.0, isMaximize: false, weight: 1.0})
	dimension.stretch(-11.5) // Stretch second best dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: -12.0, second: -11.5, base: 5.0, isMaximize: false, weight: 1.0})
	dimension.stretch(-11.5) // Doesn't stretch second best dimension. Equal.
	c.Assert(dimension, Equals, hypercubeDimension{first: -12.0, second: -11.5, base: 5.0, isMaximize: false, weight: 1.0})
	dimension.stretch(-11.0) // Doesn't stretch dimension. Less than.
	c.Assert(dimension, Equals, hypercubeDimension{first: -12.0, second: -11.5, base: 5.0, isMaximize: false, weight: 1.0})
	dimension.stretch(-12.0) // Stretch second best dimension.
	c.Assert(dimension, Equals, hypercubeDimension{first: -12.0, second: -12.0, base: 5.0, isMaximize: false, weight: 1.0})
}

func (s *HypervolumeIndicatorSortSuite) Test_StrechDimensions(c *C) {

	// Start with some hypercube.
	var hypercube []hypercubeDimension = []hypercubeDimension{
		newHypercubeDimension(-5.0, true, 1.0),
		newHypercubeDimension(5.0, false, 1.0),
		newHypercubeDimension(-5.0, true, 1.0),
	}
	var expectedHypercube []hypercubeDimension = []hypercubeDimension{
		hypercubeDimension{first: -5.0, second: -5.0, base: -5.0, isMaximize: true, weight: 1.0},
		hypercubeDimension{first: 5.0, second: 5.0, base: 5.0, isMaximize: false, weight: 1.0},
		hypercubeDimension{first: -5.0, second: -5.0, base: -5.0, isMaximize: true, weight: 1.0},
	}
	c.Assert(hypercube, DeepEquals, expectedHypercube)

	// Stretch all dimensions.
	expectedHypercube = []hypercubeDimension{
		hypercubeDimension{first: -4.0, second: -5.0, base: -5.0, isMaximize: true, weight: 1.0},
		hypercubeDimension{first: 4.0, second: 5.0, base: 5.0, isMaximize: false, weight: 1.0},
		hypercubeDimension{first: -3.0, second: -5.0, base: -5.0, isMaximize: true, weight: 1.0},
	}
	hypercube = stretchDimensions(hypercube, []float64{-4.0, 4.0, -3.0})
	c.Assert(hypercube, DeepEquals, expectedHypercube)

	// Stretch all dimensions.
	expectedHypercube = []hypercubeDimension{
		hypercubeDimension{first: 100.0, second: -4.0, base: -5.0, isMaximize: true, weight: 1.0},
		hypercubeDimension{first: -100.0, second: 4.0, base: 5.0, isMaximize: false, weight: 1.0},
		hypercubeDimension{first: 99.0, second: -3.0, base: -5.0, isMaximize: true, weight: 1.0},
	}
	hypercube = stretchDimensions(hypercube, []float64{100.0, -100.0, 99.0})
	c.Assert(hypercube, DeepEquals, expectedHypercube)

	// Invalid parameters.
	c.Assert(func() { stretchDimensions(hypercube, []float64{100.0, -100.0}) }, Panics, `stretchDimensions expects 3 dimensions, but outcomes have 2 dimensions`)

}

func (s *HypervolumeIndicatorSortSuite) Test_CalculateHypercubeContribution(c *C) {
	var indicator float64
	var volume float64

	// Start with some hypercube.
	var hypercube []hypercubeDimension = []hypercubeDimension{
		hypercubeDimension{first: 5.0, second: 1.0, base: -1.0, isMaximize: true, weight: 3.0}, // Examine different weights.
		hypercubeDimension{first: 5.0, second: 2.0, base: -2.0, isMaximize: true, weight: 2.0},
		hypercubeDimension{first: 5.0, second: 3.0, base: -3.0, isMaximize: true, weight: 1.0},
	}

	// A point that doesn't define any max of the hypercube dimension.
	indicator, volume = calculateHypercubeContribution(hypercube, []float64{4.0, 4.0, 4.0})
	c.Check(indicator, Equals, 0.0)  // Every dimension dominated by some other member of population.
	c.Assert(volume, Equals, 1260.0) // 15 * 12 * 7 == (4.0-(-1.0)) x 3.0 * (4.0-(-2.0)) x 2.0 * (4.0-(-3.0)) x 1.0

	// Now examine the point that defines the hypervolume. It dominates every other point.
	indicator, volume = calculateHypercubeContribution(hypercube, []float64{5.0, 5.0, 5.0})
	c.Check(indicator, Equals, 144.0) // 12 * 6 * 2 == (5.0-1.0) x 3.0 * (5.0-2.0) x 2.0 * (5.0-3.0) x 1.0
	c.Assert(volume, Equals, 2016.0)  // 18 * 14 * 8 == (5.0-(-1.0)) x 3.0 * (5.0-(-2.0)) x 2.0 * (5.0-(-3.0)) x 1.0

	// What if we have a point that only stretches one dimension?
	indicator, volume = calculateHypercubeContribution(hypercube, []float64{5.0, 1.0, 1.0})
	c.Check(indicator, Equals, 12.0) // 12 == (5.0-1.0) x 3.0
	c.Assert(volume, Equals, 432.0)  // 18 * 6 * 4 == (5.0-(-1.0)) x 3.0 * (1.0-(-2.0)) x 2.0 * (1.0-(-3.0)) x 1.0

	// What if we have a point that only stretches one dimension?
	indicator, volume = calculateHypercubeContribution(hypercube, []float64{1.0, 5.0, 1.0})
	c.Check(indicator, Equals, 6.0) // 6 == (5.0-2.0) x 2.0
	c.Assert(volume, Equals, 336.0) // 6 * 14 * 4 == (1.0-(-1.0)) x 3.0 * (5.0-(-2.0)) x 2.0 * (1.0-(-3.0)) x 1.0

	// What if we have a point that only stretches one dimension?
	indicator, volume = calculateHypercubeContribution(hypercube, []float64{1.0, 1.0, 5.0})
	c.Check(indicator, Equals, 2.0) // 2 == (5.0-3.0) x 1.0
	c.Assert(volume, Equals, 288.0) // 6 * 6 * 8 == (1.0-(-1.0)) x 3.0 * (1.0-(-2.0)) x 2.0 * (5.0-(-3.0)) x 1.0

	// What about a hypercube with "zero values", the base reference point
	indicator, volume = calculateHypercubeContribution(hypercube, []float64{-1.0, 5.0, -3.0})
	c.Check(indicator, Equals, 6.0) // (5.0-2.0) x 2.0
	c.Assert(volume, Equals, 0.0)   // Any dimension at or below the reference point will have no volume.

	// What about a hypercube below "zero values", the base reference point
	indicator, volume = calculateHypercubeContribution(hypercube, []float64{-10.0, 5.0, -30.0})
	c.Check(indicator, Equals, 6.0) // (5.0-2.0) x 2.0
	c.Assert(volume, Equals, 0.0)   // Any dimension at or below the reference point will have no volume.

	// What about a hypercube completely at "zero values", the base reference point
	indicator, volume = calculateHypercubeContribution(hypercube, []float64{-1.0, -2.0, -3.0})
	c.Check(indicator, Equals, 0.0) // Totally dominated.
	c.Assert(volume, Equals, 0.0)   // No volume, the ultimately bad specimens.

	// Throw in a minimized value.
	hypercube = []hypercubeDimension{
		hypercubeDimension{first: 5.0, second: 1.0, base: -1.0, isMaximize: true, weight: 3.0},
		hypercubeDimension{first: -2.0, second: 5.0, base: 6.0, isMaximize: false, weight: 2.0}, // Minimize value.
		hypercubeDimension{first: 5.0, second: 3.0, base: -3.0, isMaximize: true, weight: 1.0},
	}

	// The point that defines the hypercube.
	indicator, volume = calculateHypercubeContribution(hypercube, []float64{5.0, -2.0, 5.0})
	c.Check(indicator, Equals, 336.0) // 12 * 14 * 2 == (5.0-1.0) x 3.0 * (5.0-(-2.0)) x 2.0 * (5.0-3.0) x 1.0
	c.Assert(volume, Equals, 2304.0)  // 18 * 16 * 8 == (5.0-(-1.0)) x 3.0 * (6.0-(-2.0)) x 2.0 * (5.0-(-3.0)) x 1.0

	// The point touches the minimizing reference point value.
	indicator, volume = calculateHypercubeContribution(hypercube, []float64{5.0, 6.0, 5.0})
	c.Check(indicator, Equals, 24.0) // 12 * 2 == (5.0-1.0) x 3.0 * (5.0-3.0) x 1.0
	c.Assert(volume, Equals, 0.0)    // Any dimension at or above a minimizing reference point will have no volume.
}

func (s *HypervolumeIndicatorSortSuite) Test_NewHypercubeContribution(c *C) {

	// Start with some hypercube.
	var hypercube []hypercubeDimension = []hypercubeDimension{
		hypercubeDimension{first: 5.0, second: 1.0, base: -1.0, isMaximize: true, weight: 3.0}, // Examine different weights.
		hypercubeDimension{first: 5.0, second: 2.0, base: -2.0, isMaximize: true, weight: 2.0},
		hypercubeDimension{first: 5.0, second: 3.0, base: -3.0, isMaximize: true, weight: 1.0},
	}

	// A specimen that has two members in its species.
	var specimen Specimen = Specimen{
		SpeciesMemberCount: 2, // Two members in the species.
		Outcomes:           []float64{5.0, 5.0, 5.0},
	}

	// The value for making a new contribution.
	var expectedContribution hypercubeContribution = hypercubeContribution{
		indicator:         144.0,  // 12 * 6 * 2 == (5.0-1.0) x 3.0 * (5.0-2.0) x 2.0 * (5.0-3.0) x 1.0
		volume:            2016.0, // 18 * 14 * 8 == (5.0-(-1.0)) x 3.0 * (5.0-(-2.0)) x 2.0 * (5.0-(-3.0)) x 1.0
		weightedIndicator: 72.0,   // Two members in species.
		weightedVolume:    1008.0, // Two members in species.
		specimen:          specimen,
	}

	// Now examine the point that defines the hypervolume. It dominates every other point.
	c.Assert(newHypercubeContribution(hypercube, specimen), DeepEquals, expectedContribution)
}

func (s *HypervolumeIndicatorSortSuite) Test_HypercubeContribution_Sort(c *C) {

	// Unsorted.
	var contributions []hypercubeContribution = []hypercubeContribution{
		hypercubeContribution{weightedIndicator: 10.0, weightedVolume: 400.0},
		hypercubeContribution{weightedIndicator: 10.0, weightedVolume: 500.0},
		hypercubeContribution{weightedIndicator: 11.0, weightedVolume: 400.0},
	}

	// Sort.
	sort.Sort(byHypervolumeIndicator(contributions))

	// Expected sorted.
	var expectedContributions []hypercubeContribution = []hypercubeContribution{
		hypercubeContribution{weightedIndicator: 11.0, weightedVolume: 400.0},
		hypercubeContribution{weightedIndicator: 10.0, weightedVolume: 500.0},
		hypercubeContribution{weightedIndicator: 10.0, weightedVolume: 400.0},
	}
	c.Assert(contributions, DeepEquals, expectedContributions)
}
