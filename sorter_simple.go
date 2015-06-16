package genetic

import (
	"fmt"
	"sort"
)

// sorterSimple sorts the specimens either up or down, depending on the experiment.
type sorterSimple struct {
	Maximize bool // True if we are looking for higher values.
}

// NewSorterSimpleMaximize create a new simple sorter that will will sort descending
// by (score + bonus) / (# members in species). Better bonuses should be larger numbers.
// To keep a single species from taking over the populations, the larger the species, the worse the sort value.
func NewSorterSimpleMaximize() *sorterSimple {
	return &sorterSimple{
		Maximize: true,
	}
}

// NewSorterSimple create a new simple sorter that will sort ascending by
// (score + bonus) x (# members in species). Better bonuses should be smaller numbers, going negative if necessary.
// To keep a single species from taking over the populations, the larger the species, the worse the sort value.
func NewSorterSimpleMinimize() *sorterSimple {
	return &sorterSimple{
		Maximize: false,
	}
}

// bySimpleSortAscending implements sort.Interface to sort ascending by (score + bonus) x (# members in species).
// Example: sort.Sort(bySimpleSortAscending(specimens))
type bySimpleSortAscending []Specimen

func (a bySimpleSortAscending) Len() int      { return len(a) }
func (a bySimpleSortAscending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a bySimpleSortAscending) Less(i, j int) bool {
	if a[i].SelectionScore == a[j].SelectionScore {
		if a[i].Score == a[j].Score {
			return a[i].Bonus < a[j].Bonus // Third by bonus.
		} else {
			return a[i].Score < a[j].Score // Second by score.
		}
	}
	return a[i].SelectionScore < a[j].SelectionScore // First by selection score.
}

// bySimpleSortDescending implements sort.Interface to sort descending by (score + bonus) / (# members in species).
// Example: sort.Sort(bySimpleSortDescending(specimens))
type bySimpleSortDescending []Specimen

func (a bySimpleSortDescending) Len() int      { return len(a) }
func (a bySimpleSortDescending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a bySimpleSortDescending) Less(i, j int) bool {
	if a[i].SelectionScore == a[j].SelectionScore {
		if a[i].Score == a[j].Score {
			return a[i].Bonus > a[j].Bonus // Third by bonus.
		} else {
			return a[i].Score > a[j].Score // Second by score.
		}
	}
	return a[i].SelectionScore > a[j].SelectionScore // First by selection score.
}

// Sort the specimens either ascending or descending.
func (s *sorterSimple) Sort(specimens []Specimen) (bestScore float64, best string) {

	// Give each specimen a selection score.
	var bestSpecimen *Specimen
	for i := 0; i < len(specimens); i++ {

		// Score based on whether we are ascending or descending.
		var selectionScore float64
		if s.Maximize {
			// Value gets smaller if specimen is in a larger species, making it a worse score.
			selectionScore = (specimens[i].Score + specimens[i].Bonus) / float64(specimens[i].SpeciesMemberCount)
		} else {
			// Value gets larger if specimen is in a larger species, making it a worse score.
			selectionScore = (specimens[i].Score + specimens[i].Bonus) * float64(specimens[i].SpeciesMemberCount)
		}
		specimens[i].setSelectionScore(selectionScore)

		// Is this the best specimen we've found?
		if bestSpecimen == nil {
			bestSpecimen = &specimens[i]
		} else {
			if s.Maximize {
				if specimens[i].Score > bestScore {
					bestSpecimen = &specimens[i]
				}
			} else {
				if specimens[i].Score < bestScore {
					bestSpecimen = &specimens[i]
				}
			}
		}
		bestScore = bestSpecimen.Score
	}

	// What is the text summary of the best specimen.
	best = fmt.Sprintf("score: %f, bonus: %f, speciesmembercount: %d", bestSpecimen.Score, bestSpecimen.Bonus, bestSpecimen.SpeciesMemberCount)

	// Sort the specimens.
	if s.Maximize {
		// Count down from left to right.
		sort.Sort(bySimpleSortDescending(specimens))
	} else {
		// Count up form left to right.
		sort.Sort(bySimpleSortAscending(specimens))
	}

	// The best value is the value at the specimen at the head of the list. It is the fittest.
	return bestScore, best
}

// IsMaximize returns true if we higher scores are fitter and false if lower scores are fitter.
func (s *sorterSimple) IsMaximize() bool { return s.Maximize }
