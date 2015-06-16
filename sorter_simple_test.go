package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type SorterSimpleSuite struct{}

var _ = Suite(&SorterSimpleSuite{})

// Add the tests.

func (s *SorterSimpleSuite) Test_Maximize(c *C) {

	// Some specimens.
	var specimenA Specimen = Specimen{Score: 5.0, Bonus: 5.0, SpeciesMemberCount: 1}  // selection score: 10.0
	var specimenB Specimen = Specimen{Score: 10.0, Bonus: 5.0, SpeciesMemberCount: 2} // selection score: 7.5, best score but not first on list
	var specimenC Specimen = Specimen{Score: 5.0, Bonus: 0.0, SpeciesMemberCount: 2}  // selection score: 2.5, pure socre more valuable than weighted score.
	var specimenD Specimen = Specimen{Score: 2.5, Bonus: 0.0, SpeciesMemberCount: 1}  // selection score: 2.5, score more valuable than bonus.
	var specimenE Specimen = Specimen{Score: 0.0, Bonus: 5.0, SpeciesMemberCount: 2}  // selection score: 2.5, pure bonus more valuable than weighted bonus.
	var specimenF Specimen = Specimen{Score: 0.0, Bonus: 2.5, SpeciesMemberCount: 1}  // selection score: 2.5

	// Put in a list that we will sort. The exact wrong order.
	var specimens []Specimen = []Specimen{specimenF, specimenE, specimenD, specimenC, specimenB, specimenA}

	// Make the sorter.
	var sorter Sorter = NewSorterSimpleMaximize()

	// Do the sort.
	var bestScore float64
	var best string
	bestScore, best = sorter.Sort(specimens)

	// Add the selection score each specimen should get.
	specimenA.SelectionScore = 10.0
	specimenB.SelectionScore = 7.5
	specimenC.SelectionScore = 2.5
	specimenD.SelectionScore = 2.5
	specimenE.SelectionScore = 2.5
	specimenF.SelectionScore = 2.5

	// What should they become?
	var expectedSpecimens []Specimen = []Specimen{specimenA, specimenB, specimenC, specimenD, specimenE, specimenF}

	// Did we get what we expected?
	c.Check(specimens, DeepEquals, expectedSpecimens)
	c.Check(bestScore, Equals, 10.0)
	c.Check(best, Equals, "score: 10.000000, bonus: 5.000000, speciesmembercount: 2")
}

func (s *SorterSimpleSuite) Test_Minimize(c *C) {

	// Some specimens.
	var specimenA Specimen = Specimen{Score: 4.0, Bonus: -1.0, SpeciesMemberCount: 1}   // selection score: 3.0
	var specimenB Specimen = Specimen{Score: 3.0, Bonus: -1.0, SpeciesMemberCount: 2}   // selection score: 4.0, best score but not first on list
	var specimenC Specimen = Specimen{Score: 5.0, Bonus: 0.0, SpeciesMemberCount: 2}    // selection score: 10.0, pure score more valuable than bonus weighted score
	var specimenD Specimen = Specimen{Score: 10.0, Bonus: 0.0, SpeciesMemberCount: 1}   // selection score: 10.0, score more valuable than bonus.
	var specimenE Specimen = Specimen{Score: 15.0, Bonus: -10.0, SpeciesMemberCount: 2} // selection score: 10.0, pure bonus more valuable than weighted bonus.
	var specimenF Specimen = Specimen{Score: 15.0, Bonus: -5.0, SpeciesMemberCount: 1}  // selection score: 10.0

	// Put in a list that we will sort. The exact wrong order.
	var specimens []Specimen = []Specimen{specimenF, specimenE, specimenD, specimenC, specimenB, specimenA}

	// Make the sorter.
	var sorter Sorter = NewSorterSimpleMinimize()

	// Do the sort.
	var bestScore float64
	var best string
	bestScore, best = sorter.Sort(specimens)

	// Add the selection score each specimen should get.
	specimenA.SelectionScore = 3.0
	specimenB.SelectionScore = 4.0
	specimenC.SelectionScore = 10.0
	specimenD.SelectionScore = 10.0
	specimenE.SelectionScore = 10.0
	specimenF.SelectionScore = 10.0

	// What should they become?
	var expectedSpecimens []Specimen = []Specimen{specimenA, specimenB, specimenC, specimenD, specimenE, specimenF}

	// Did we get what we expected?
	c.Check(specimens, DeepEquals, expectedSpecimens)
	c.Check(bestScore, Equals, 3.0)
	c.Check(best, Equals, "score: 3.000000, bonus: -1.000000, speciesmembercount: 2")
}
