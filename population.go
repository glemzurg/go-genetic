package genetic

import (
	"log"
	"math/rand"
)

// generationPopulation is all the specimens of a single generation.
type generationPopulation struct {
	config  PopulationConfig
	species []Species
}

// newPopulation creates a well-formed Population, ready for specimens to be added.
func newPopulation(config PopulationConfig) generationPopulation {
	return generationPopulation{config: config}
}

// AddNeuralNet adds a single member of the population, putting it into the appropriate species.
func (p *generationPopulation) AddNeuralNet(neuralNet NeatNeuralNet, score float64, bonus float64, outcomes []float64) {
	// Wrap this neural net as a specimen.
	var specimen Specimen = newSpecimen(neuralNet, score, bonus, outcomes)
	p.AddSpecimen(specimen)
}

// AddSpecimen adds a single member of the population, putting it into the appropriate species.
func (p *generationPopulation) AddSpecimen(specimen Specimen) {

	// Add this specimen to the first species it is related to.
	// If it is related to none, then create a new species.
	// The order of the species is important since it may fit into two, but will fall into the first.
	var wasAdded bool = false
	for i := range p.species {
		if wasAdded = p.species[i].AddSpecimen(specimen, p.config.Speciation); wasAdded {
			break
		}
	}

	// If we didn't add this specimen, create a new species.
	if !wasAdded {
		p.species = append(p.species, newSpecies(specimen))
	}
}

// FillOut grows the population from the fittest of last generation to a full population by mutation and mating.
func (p *generationPopulation) FillOut() {

	// Prepare the species for pulling random specimens.
	var specimenCount int = p.prepareRandomSpecimenIndexes()

	// Gather all the new specimens.
	var newSpecimens []Specimen

	// How many more specimens do we need?
	var specimensNeeded int = (p.config.PopulationSize - specimenCount)
	for i := 0; i < specimensNeeded; i++ {

		// Pick a random specimen.
		var specimen Specimen
		var speciesSpecimens []Specimen
		var specimenIndex int
		specimen, speciesSpecimens, specimenIndex = p.randomSpecimen(specimenCount)

		// Create a new specimen from an random change of this one.
		var mutant Specimen = specimen.MateMutate(speciesSpecimens, specimenIndex, p.config.Mutate)
		newSpecimens = append(newSpecimens, mutant)
	}

	// Add all the new specimens to species, creating any needed to house them.
	for _, specimen := range newSpecimens {
		p.AddSpecimen(specimen)
	}
}

// prepareRandomSpecimenIndexes prepares species with indexes allowing random specimens to be picked.
func (p *generationPopulation) prepareRandomSpecimenIndexes() (specimenCount int) {

	// How many specimens are there? This population has been pruned of dead weight. We don't know how many
	// specimens there currently are or what species they are in. Calculate some indexes that would be true
	// if all specimens where in one slice, ordered by their species.
	var nextIndex int = 0
	for i := range p.species {
		// How many specimens are in this species?
		var specimenCount int = len(p.species[i].Specimens)

		// Note the index of the first specimen in the species.
		//
		// Assume that we have 5 members of the species and the first population index
		// of this species is 2. The members of this population would have these indexes:
		//
		//   [2, 3, 4, 5, 6]
		//
		//   2 = firstIndex = nextIndex
		//   6 = lastIndex = nextIndex + specimenCount - 1
		//
		var firstIndex int = nextIndex
		var lastIndex int = nextIndex + specimenCount - 1

		// Remember the indexes.
		p.species[i].firstPopulationIndex = firstIndex
		p.species[i].lastPopulationIndex = lastIndex

		// Next index?
		nextIndex = lastIndex + 1
	}

	// The specimen count is actually the next index, because it is higher than
	// the last valid index.
	specimenCount = nextIndex
	return specimenCount
}

// randomSpecimen picks a random specimen from a population with enough supporting data to mate if necessary.
func (p *generationPopulation) randomSpecimen(specimenCount int) (specimen Specimen, speciesSpecimens []Specimen, specimenIndex int) {

	// Pick a random specimen from the whole population.
	var populationSpecimenIndex int = rand.Intn(specimenCount)

	// Find the specimen we want.
	for i := range p.species {

		// Get the specimen from the species, return if found.
		var specimenFound bool
		if specimen, speciesSpecimens, specimenIndex, specimenFound = p.species[i].pickSpecimen(populationSpecimenIndex); specimenFound {
			return specimen, speciesSpecimens, specimenIndex
		}
	}

	log.Panic("Specimen not found.")
	return
}

// DumpSpecimens removes all the specimens from the population and returns them, fit for selection.
func (p *generationPopulation) DumpSpecimens() []Specimen {
	// Empty out the species of their specimens, but keep them there, they need to stick around for
	// recategorization later.
	var specimens []Specimen
	for i := range p.species {
		for _, specimen := range p.species[i].Specimens {
			specimens = append(specimens, specimen)
		}
		p.species[i].Specimens = nil
	}
	return specimens
}

// DumpSpecimensAsNeuralNets removes all the specimens from the population and returns them, fit for scoring.
func (p *generationPopulation) DumpSpecimensAsNeuralNets() []NeatNeuralNet {
	// Empty out the species of their specimens, but keep them there, they need to stick around for
	// recategorization later.
	var neuralNets []NeatNeuralNet
	for i := range p.species {
		for _, specimen := range p.species[i].Specimens {
			neuralNets = append(neuralNets, specimen.NeuralNet)
		}
		p.species[i].Specimens = nil
	}
	return neuralNets
}

// WeightSpecies weights all the specimen scores by the size of their species.
func (p *generationPopulation) WeightSpecies() {
	for i := range p.species {
		p.species[i].WeightSpecies()
	}
}

// AddAllSpecimens restocks the population with specimens.
func (p *generationPopulation) AddAllSpecimens(specimens []Specimen) {
	for _, specimen := range specimens {
		p.AddSpecimen(specimen)
	}
	p.PruneEmptySpecies()
}

// PruneEmptySpecies removes any species that has no more specimens.
func (p *generationPopulation) PruneEmptySpecies() {
	// Preserve order of species. Matters to keep specimens always categorizing into the same species every generation.
	var speciesToKeep []Species
	for _, species := range p.species {
		if len(species.Specimens) > 0 {
			speciesToKeep = append(speciesToKeep, species)
		}
	}
	p.species = speciesToKeep
}
