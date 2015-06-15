package genetic

// Sorter orders the specimens for the Selector.
type Sorter interface {

	// Order the specimens and report how well the population did.
	Sort(specimens []Specimen) (bestScore float64, best string)

	// IsMaximize returns true if this experiment is seeking higher values, false if seeking lower values.
	IsMaximize() bool
}
