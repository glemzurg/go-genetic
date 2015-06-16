package genetic

// genSpecies is a collection of specimens deemed to be alike due to similarities in their genomes.
type genSpecies struct {
	genome               neatGenome
	Specimens            []Specimen
	firstPopulationIndex int // The index of the first specimen, as if all population specimens were in one slice.
	lastPopulationIndex  int // The index of the last specimen, as if all population specimens were in one slice.
}

// newSpecies creates a well-formed species of the population.
func newSpecies(specimen Specimen) genSpecies {
	return genSpecies{
		// The species will have the identiy genome of the specimen.
		// Make a copy of the genes so it is not tethered to the specimen itself.
		genome:    specimen.NeuralNet.Genome.Clone(),
		Specimens: []Specimen{specimen}, // Specimen speciation distance is 0.0 == it is the genome.
	}
}

// AddSpecimen adds a single member of the species, return true if it was added.
func (s *genSpecies) AddSpecimen(specimen Specimen, config ConfigSpeciation) (wasAdded bool) {
	var isSpecies bool
	var speciationDistance float64
	if isSpecies, speciationDistance = isSameSpecies(s.genome, specimen.NeuralNet.Genome, config); isSpecies {
		// This specimen is a member of this species.
		// Stamp the speciation distance on them and add them.
		specimen.SpeciationDistance = speciationDistance
		s.Specimens = append(s.Specimens, specimen)
	}
	return isSpecies
}

// pickSpecimen picks the specimen from the species if it is in the species. Returns false if not found.
func (s *genSpecies) pickSpecimen(populationSpecimenIndex int) (specimen Specimen, speciesSpecimens []Specimen, specimenIndex int, wasFound bool) {

	// If this index is not in species, return that we didn't find anything.
	if populationSpecimenIndex < s.firstPopulationIndex || populationSpecimenIndex > s.lastPopulationIndex {
		return Specimen{}, nil, 0, false
	}

	// The specimen is in this species. What index is it?
	//
	// Assume that we have 5 members of the species and the first population index
	// of this species is 2. The members of this population would have these indexes:
	//
	//   [2, 3, 4, 5, 6]
	//
	//   pop 2 = firstIndex
	//   pop 6 = lastIndex
	//
	// Say we want the specimen with population index 4 (index 2 in the local specimens)
	//
	//   localIndex 2 = populationIndex 4 - firstIndex 2
	//
	var localIndex int = populationSpecimenIndex - s.firstPopulationIndex

	// Return the find.
	return s.Specimens[localIndex], s.Specimens, localIndex, true
}

// WeightSpecies adds to the specimens enough information to weigh by species in scoring.
func (s *genSpecies) WeightSpecies() {

	// How many specimens are there?
	var specimenCount int = len(s.Specimens)
	for i := range s.Specimens {
		// The more specimens in the species, the worse the score. Keeps one species from
		// taking over the whole population.
		s.Specimens[i].SpeciesMemberCount = specimenCount
	}
}
