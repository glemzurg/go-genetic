package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type SorterHypervolumeIndicatorSuite struct{}

var _ = Suite(&SorterHypervolumeIndicatorSuite{})

// Add the tests.

func (s *SorterHypervolumeIndicatorSuite) Test_ValidOrPanic(c *C) {

	// A well-formed sorter.
	var goodSorter SorterHypervolumeIndicator = SorterHypervolumeIndicator{
		ReferencePoint: []float64{0.0, 0.0},
		Maximize:       []bool{true, true},
		Weights:        []float64{1.0, 1.0},
	}
	var sorter SorterHypervolumeIndicator

	// The well-formed sorter doesn't panic.
	goodSorter.validOrPanic()

	// No reference point panics.
	sorter = goodSorter
	sorter.ReferencePoint = nil
	c.Check(func() { sorter.validOrPanic() }, Panics, "ReferencePoint must have values")

	// No reference point panics.
	sorter = goodSorter
	sorter.ReferencePoint = []float64{}
	c.Check(func() { sorter.validOrPanic() }, Panics, "ReferencePoint must have values")

	// Different length parts panic.
	sorter = goodSorter
	sorter.ReferencePoint = append(sorter.ReferencePoint, 0.0)
	c.Check(func() { sorter.validOrPanic() }, Panics, "ReferencePoint ([0 0 0]), Maximize ([true true]), and Weights ([1 1]) must all have the same number of values")

	// Different length parts panic.
	sorter = goodSorter
	sorter.Maximize = append(sorter.Maximize, false)
	c.Check(func() { sorter.validOrPanic() }, Panics, "ReferencePoint ([0 0]), Maximize ([true true false]), and Weights ([1 1]) must all have the same number of values")

	// Different length parts panic.
	sorter = goodSorter
	sorter.Weights = append(sorter.Weights, 1.0)
	c.Check(func() { sorter.validOrPanic() }, Panics, "ReferencePoint ([0 0]), Maximize ([true true]), and Weights ([1 1 1]) must all have the same number of values")

	// Any zero-weight panics.
	sorter = goodSorter
	sorter.Weights[1] = 0.0
	c.Check(func() { sorter.validOrPanic() }, Panics, "Weights with a zero value are not allowed: [1 0]")
}

func (s *SorterHypervolumeIndicatorSuite) Test_SorterHypervolumeIndicator(c *C) {

	// Some specimens.
	var specimenA Specimen = Specimen{Outcomes: []float64{1.5, 1.5}, SpeciesMemberCount: 1} // selection score: 2.5, indicator: 0.25, volume: 2.25
	var specimenB Specimen = Specimen{Outcomes: []float64{2.0, 1.0}, SpeciesMemberCount: 1} // selection score: 2.5, indicator: 0.5, volume: 2.0
	var specimenC Specimen = Specimen{Outcomes: []float64{1.0, 2.0}, SpeciesMemberCount: 2} // selection score: 1.25, indicator: 0.5, volume: 2.0
	var specimenD Specimen = Specimen{Outcomes: []float64{0.5, 1.3}, SpeciesMemberCount: 1} // selection score: 0.65, dominated: 0.0, volume: 0.65

	// Put in a list that we will sort. The exact wrong order.
	var specimens []Specimen = []Specimen{specimenD, specimenC, specimenB, specimenA}

	// Make the sorter.
	var sorter Sorter = &SorterHypervolumeIndicator{
		ReferencePoint: []float64{0.0, 0.0},
		Maximize:       []bool{true, true},
		Weights:        []float64{1.0, 1.0},
	}

	// Do the sort.
	var bestScore float64
	var best string
	var sorted []Specimen
	bestScore, best, sorted = sorter.Sort(specimens)

	// Add the selection score each specimen should get.
	specimenA.SelectionScore = 2.5
	specimenB.SelectionScore = 2.5
	specimenC.SelectionScore = 1.25
	specimenD.SelectionScore = 0.65

	// What should they become?
	var expectedSpecimens []Specimen = []Specimen{specimenA, specimenB, specimenC, specimenD}

	// Did we get what we expected?
	c.Check(sorted, DeepEquals, expectedSpecimens)
	c.Check(bestScore, Equals, 2.25)                                                                           // Volume is the ultimate best score.
	c.Check(best, Equals, "indicator: 0.250000, volume: 2.250000, speciesmembercount: 1, outcomes: [1.5 1.5]") // Volume is the ultimate best score.
}
