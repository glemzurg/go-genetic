package genetic

// Selector selects which members of a population to keep for the next generation.
type Selector interface {

	// Pick the population members to continue on to the next generation.
	Select(specimens []Specimen) (fittest []Specimen)
}
