package genetic

// Sorter orders the specimens for the Selector and gives each a selector score, a score used by some selectors
// to pick the fittest members of a population.
type Sorter interface {

	// Order the specimens and report how well the population did.
	Sort(specimens []Specimen) (bestScore float64, best string, sorted []Specimen)

	// IsMaximize returns true if this experiment is seeking higher values, false if seeking lower values.
	IsMaximize() bool
}
