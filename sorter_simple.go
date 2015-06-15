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
// by (score + bonus) / (# members in species). To keep a single species from taking over the populations, the larger
// the species, the worse the sort value.
func NewSorterSimpleMaximize() *sorterSimple {
	return &sorterSimple{
		Maximize: true,
	}
}

// NewSorterSimple create a new simple sorter that will sort ascending by
// (score + bonus) x (# members in species). To keep a single species from taking over the populations, the larger
// the species, the worse the sort value.
func NewSorterSimpleMinimize() *sorterSimple {
	return &sorterSimple{
		Maximize: false,
	}
}

// bySpeciesScoreAscending implements sort.Interface to sort ascending by (score + bonus) x (# members in species).
// Example: sort.Sort(bySpeciesScoreAscending(specimens))
type bySpeciesScoreAscending []Specimen

func (a bySpeciesScoreAscending) Len() int      { return len(a) }
func (a bySpeciesScoreAscending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a bySpeciesScoreAscending) Less(i, j int) bool {
	return (a[i].Score+a[i].Bonus)*float64(a[i].SpeciesMemberCount) < (a[j].Score+a[j].Bonus)*float64(a[j].SpeciesMemberCount)
}

// bySpeciesScoreDescending implements sort.Interface to sort descending by (score + bonus) / (# members in species).
// Example: sort.Sort(bySpeciesScoreDescending(specimens))
type bySpeciesScoreDescending []Specimen

func (a bySpeciesScoreDescending) Len() int      { return len(a) }
func (a bySpeciesScoreDescending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a bySpeciesScoreDescending) Less(i, j int) bool {
	return (a[i].Score+a[i].Bonus)/float64(a[i].SpeciesMemberCount) < (a[j].Score+a[j].Bonus)/float64(a[j].SpeciesMemberCount)
}

// Sort the specimens either ascending or descending.
func (s *sorterSimple) Sort(specimens []Specimen) (bestScore float64, best string) {

	// Sort the specimens.
	if s.Maximize {
		sort.Sort(bySpeciesScoreDescending(specimens))
	} else {
		sort.Sort(bySpeciesScoreAscending(specimens))
	}

	// What is the best score for this generation?
	bestScore = specimens[0].Score

	// The best value is the value at the specimen at the head of the list. It is the fittest.
	return bestScore, fmt.Sprintf("score: %f, bonus: %f, speciesmembercount: %d", bestScore, specimens[0].Bonus, specimens[0].SpeciesMemberCount)
}

// IsMaximize returns true if we higher scores are fitter and false if lower scores are fitter.
func (s *sorterSimple) IsMaximize() bool { return s.Maximize }
