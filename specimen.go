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

// Specimen is a single member of a population, scored.
type Specimen struct {
	NeuralNet          NeatNeuralNet // The neural net that defines this specimen.
	Score              float64       // The score of the neural net for selectors that use it. 0.0 if unused.
	Bonus              float64       // The bonus for meta-qualiteis (e.g. novelty searches). 0.0 if unused.
	Outcomes           []float64     // Multi-outcomes for selectors that use it (e.g. hypervolume indicator). null if unused.
	SelectionScore     float64       // This is the score the specimen will ultimately be sorted on before being passed to the Selector.
	SpeciationDistance float64       // How different this specimen's genome is from the species identity genome.
	SpeciesMemberCount int           // How many specimens are in this specimen's species (including itself).
}

// newSpecimen creates a well-formed member of the population.
func newSpecimen(neuralNet NeatNeuralNet, score float64, bonus float64, outcomes []float64) Specimen {
	return Specimen{
		NeuralNet:          neuralNet,
		Score:              score,
		Bonus:              bonus,
		Outcomes:           outcomes,
		SelectionScore:     0.0, // Not calculated yet.
		SpeciationDistance: 0.0, // Not calculated yet.
		SpeciesMemberCount: 0,   // Not calculated yet.
	}
}

// setSelectionScore updates the selection score for the specimen.
func (s *Specimen) setSelectionScore(selectionScore float64) { s.SelectionScore = selectionScore }

// MateMutate produces another Specimen by modifying this specimen. It could be a mutated version or a child
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
