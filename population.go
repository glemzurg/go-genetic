package genetic

import (
	"log"
	"math/rand"
)

const (
	// Enumerations for picking the random kind of mate/mutation.
	_CHANGE_MATE = iota
	_CHANGE_MUTATE_ADD_NODE
	_CHANGE_MUTATE_ADD_CONNECTION
	_CHANGE_MUTATE_ALTER_CONNECTION
)

// Population is all the specimens of a single generation.
type Population struct {
	config  PopulationConfig
	species []Species
}

// NewPopulation creates a well-formed Population, ready for specimens to be added.
func NewPopulation(config PopulationConfig) Population {
	return Population{config: config}
}

// AddNeuralNet adds a single member of the population, putting it into the appropriate species.
func (p *Population) AddNeuralNet(neuralNet NeatNeuralNet, score float64, bonus float64, outcomes []float64) {
	// Wrap this neural net as a specimen.
	var specimen Specimen = newSpecimen(neuralNet, score, bonus, outcomes)
	p.AddSpecimen(specimen)
}

// AddSpecimen adds a single member of the population, putting it into the appropriate species.
func (p *Population) AddSpecimen(specimen Specimen) {

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
func (p *Population) FillOut() {

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
func (p *Population) prepareRandomSpecimenIndexes() (specimenCount int) {

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
func (p *Population) randomSpecimen(specimenCount int) (specimen Specimen, speciesSpecimens []Specimen, specimenIndex int) {

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
func (p *Population) DumpSpecimens() []Specimen {
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
func (p *Population) DumpSpecimensAsNeuralNets() []NeatNeuralNet {
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
func (p *Population) WeightSpecies() {
	for i := range p.species {
		p.species[i].WeightSpecies()
	}
}

// AddAllSpecimens restocks the population with specimens.
func (p *Population) AddAllSpecimens(specimens []Specimen) {
	for _, specimen := range specimens {
		p.AddSpecimen(specimen)
	}
	p.PruneEmptySpecies()
}

// PruneEmptySpecies removes any species that has no more specimens.
func (p *Population) PruneEmptySpecies() {
	// Preserve order of species. Matters to keep specimens always categorizing into the same species every generation.
	var speciesToKeep []Species
	for _, species := range p.species {
		if len(species.Specimens) > 0 {
			speciesToKeep = append(speciesToKeep, species)
		}
	}
	p.species = speciesToKeep
}

// Species is a collection of specimens deemed to be alike due to similarities in their genomes.
type Species struct {
	genome               NeatGenome
	Specimens            []Specimen
	firstPopulationIndex int // The index of the first specimen, as if all population specimens were in one slice.
	lastPopulationIndex  int // The index of the last specimen, as if all population specimens were in one slice.
}

// newSpecies creates a well-formed species of the population.
func newSpecies(specimen Specimen) Species {
	return Species{
		// The species will have the identiy genome of the specimen.
		// Make a copy of the genes so it is not tethered to the specimen itself.
		genome:    specimen.NeuralNet.Genome.Clone(),
		Specimens: []Specimen{specimen}, // Specimen speciation distance is 0.0 == it is the genome.
	}
}

// AddSpecimen adds a single member of the species, return true if it was added.
func (s *Species) AddSpecimen(specimen Specimen, config SpeciationConfig) (wasAdded bool) {
	var isSameSpecies bool
	var speciationDistance float64
	if isSameSpecies, speciationDistance = IsSameSpecies(s.genome, specimen.NeuralNet.Genome, config); isSameSpecies {
		// This specimen is a member of this species.
		// Stamp the speciation distance on them and add them.
		specimen.SpeciationDistance = speciationDistance
		s.Specimens = append(s.Specimens, specimen)
	}
	return isSameSpecies
}

// pickSpecimen picks the specimen from the species if it is in the species. Returns false if not found.
func (s *Species) pickSpecimen(populationSpecimenIndex int) (specimen Specimen, speciesSpecimens []Specimen, specimenIndex int, wasFound bool) {

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

// WeightSpecies weights adds to the specimens enough information to weigh by species in scoring.
func (s *Species) WeightSpecies() {

	// How many specimens are there?
	var specimenCount int = len(s.Specimens)
	for i := range s.Specimens {
		// The more specimens in the species, the lower the score. Keeps one species from
		// taking over the whole population.
		s.Specimens[i].SpeciesMemberCount = specimenCount
	}
}

// Specimen is a single member of a population, scored.
type Specimen struct {
	NeuralNet          NeatNeuralNet // The neural net that defines this specimen.
	Score              float64       // The score of the neural net for selectors that use it. 0.0 if unused.
	Bonus              float64       // The bonus for meta-qualiteis (e.g. novelty searches). 0.0 if unused.
	Outcomes           []float64     // Multi-outcomes for selectors that use it (e.g. hypervolume indicator). null if unused.
	SpeciationDistance float64       // The speciation distance from the species this specimen is in.
	SpeciesMemberCount int           //  How many specimens are in this specimen's species (including itself).
}

// newSpecimen creates a well-formed member of the population.
func newSpecimen(neuralNet NeatNeuralNet, score float64, bonus float64, outcomes []float64) Specimen {
	return Specimen{
		NeuralNet:          neuralNet,
		Score:              score,
		Bonus:              bonus,
		Outcomes:           outcomes,
		SpeciationDistance: 0.0, // Not calculated yet.
		SpeciesMemberCount: 0,   // Not calculated yet.
	}
}

// MateMutate produces another Specimen from modifying this specimen. It could be a mutated version or a child
// from mating. Mating can only be done with other members of the species. specimenIndex is this specimens index
// in the list (don't want to mate with self).
func (s *Specimen) MateMutate(speciesSpecimens []Specimen, specimenIndex int, config MutateConfig) Specimen {

	// Get the weights.
	var mateWeight uint = config.MateWeight
	var addNodeWeight uint = config.AddNodeWeight
	var addConnectionWeight uint = config.AddConnectionWeight
	var alterConnectionWeight uint = config.AlterConnectionWeight

	// If there is only one member of this species, we can't mate. It's this specimen.
	if len(speciesSpecimens) == 1 {
		mateWeight = 0
	}

	// Pick the type of change we're going to make. Then make it.
	var newNeuralNet NeatNeuralNet
	var changeType int = randomMateMutatePick(mateWeight, addNodeWeight, addConnectionWeight, alterConnectionWeight)
	switch changeType {

	case _CHANGE_MATE:
		// Pick another member of the species to mate with.
		var fitterParent Specimen = *s // Assume this is the fitter parent.
		var otherParent Specimen = randomSpecimenWithSkip(speciesSpecimens, specimenIndex)
		newNeuralNet = Mate(fitterParent.NeuralNet, otherParent.NeuralNet)

	case _CHANGE_MUTATE_ADD_NODE:
		newNeuralNet = s.NeuralNet.Clone()
		newNeuralNet.MutateAddNode(config.AvailableNodeFunctions)

	case _CHANGE_MUTATE_ADD_CONNECTION:
		newNeuralNet = s.NeuralNet.Clone()
		var added bool
		if added = newNeuralNet.MutateAddConnection(config.MaxAddConnectionAttempts); !added {
			// If we didn't succesfully add a connection, fall back to just altering a connection weight.
			newNeuralNet.MutateChangeConnectionWeight()
		}

	case _CHANGE_MUTATE_ALTER_CONNECTION:
		newNeuralNet = s.NeuralNet.Clone()
		newNeuralNet.MutateChangeConnectionWeight()

	default:
		log.Panicf("Unknown change type: %d", changeType)
	}

	// The new member of the population.
	return newSpecimen(newNeuralNet, 0.0, 0.0, nil) // No scores
}

// randomMateMutatePick randomly selects the kind of change we want to make to create a new member of the population.
func randomMateMutatePick(mateWeight uint, addNodeWeight uint, addConnectionWeight uint, alterConnectionWeight uint) int {
	// Randomly pick a kind of mutation based on the weighting factors.

	// If all values are zero they are of equal weight (make them all 1)
	if mateWeight == 0 && addNodeWeight == 0 && addConnectionWeight == 0 && alterConnectionWeight == 0 {
		mateWeight = 1
		addNodeWeight = 1
		addConnectionWeight = 1
		alterConnectionWeight = 1
	}

	// Hypothetically, imagine that we have these weights:
	//
	//   mate: 1
	//   add_node: 2
	//   add_connection: 3
	//   change_weight: 4
	//
	// We want to pick one of 10 values and choose the appropriate mutation.
	// We'll end up getting a number between 0 and 9.
	//
	//   [0]       -> mate     // random < mate
	//   [1,2]     -> add_node // random < mate + add_node
	//   [3,4,5]   -> add_connection // random < mate + add_node + add_connection
	//   [6,7,8,9] -> change_weight  // random < mate + add_node + add_connection + change_weight
	//
	var pickIndex uint = uint(rand.Intn(int(mateWeight + addNodeWeight + addConnectionWeight + alterConnectionWeight)))
	if pickIndex < mateWeight {
		return _CHANGE_MATE
	} else if pickIndex < (mateWeight + addNodeWeight) {
		return _CHANGE_MUTATE_ADD_NODE
	} else if pickIndex < (mateWeight + addNodeWeight + addConnectionWeight) {
		return _CHANGE_MUTATE_ADD_CONNECTION
	}
	// The else case is just returning the last option.
	return _CHANGE_MUTATE_ALTER_CONNECTION
}

// randomSpecimenWithSkip picks a random specimen from the list, skipping the specimen at the given index.
func randomSpecimenWithSkip(specimens []Specimen, skipIndex int) Specimen {
	// Since we are picking one less than the number of list items, use a length of one less.
	var pickIndex int = rand.Intn(len(specimens) - 1)
	// If the index is the index we want to skip or higher, we need to shift it up one to account for
	// the missing list item.
	if pickIndex >= skipIndex {
		pickIndex++
	}

	// Get that specimen and return it.
	return specimens[pickIndex]
}
