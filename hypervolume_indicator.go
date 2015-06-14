package genetic

import (
	// "encoding/json"
	// "io/ioutil"
	"log"
)

// HypervolumeIndicator evaluates multi-outcome scores and scores each on how much it stretches
// the whole population's improvement by looking at all the specimens as overlapping cubes (hyper cubes).
type HypervolumeIndicatorContextScorer struct {
	ReferencePoint    []float64 // What is the base reference point to use for calculating hypercubes?
	IsMaximize        []bool    // For each dimension, are we trying to maximize or minimize the value?
	DimensionalWeight []float64 // For each dimension, do some contribute more to the score than others?
}

// dimensionMax is the tuple of best and second best points in a dimension.
type dimensionMax struct {
	first      float64 // The value that most stretches the hypercube volume in this dimension.
	second     float64 // The value that second most stretche the hypercube volume in this dimension.
	isMaximize bool    // If true, values subtract the base to determine their contribution. If false, the base subtracts the value to determine the contribution.
}

// newDimensionMax creates a well-formed dimensionMax.
func newDimensionMax(base float64, isMaximize bool) dimensionMax {
	return dimensionMax{
		first:      base,       // All values allowed must be on one side of the reference point value.
		second:     base,       // All values allowed must be on one side of the reference point value.
		isMaximize: isMaximize, // Determines which side of of the reference point value is valid side.
	}
}

// stretch pushes the first or second best dimension max farther if the value shoud do so, or leaves the max unchanged.
func (m *dimensionMax) stretch(value float64) {
	// If value is better than the fist value, knock the first value to the second place, ejecting that value.
	// If value is worse than the first and better than teh second, knock the second place value from the struct.
	// If the value is neither, leave the struct unchanged.
	if m.isMaximize {
		switch {
		case value > m.first:
			m.second = m.first
			m.first = value
		case value > m.second:
			m.second = value
		}
	} else {
		switch {
		case value < m.first:
			m.second = m.first
			m.first = value
		case value < m.second:
			m.second = value
		}
	}
}

// MultiOutcomePopulationContextScore gives each specimen a single score based on how its multi-outcome
// relates to the entire population.
func (s *HypervolumeIndicatorContextScorer) MutliOutcomePopulationContextScore(specimens []Specimen) (scored []Specimen) {

	// First we want to get a feel for the scope of the volume of the whole population.
	// What are the dimensions of the hypercube defined by the best score in every dimenion?
	// We also want to know the second-best hypercube since we'll be considering the value of
	// each specimen if it was removed from the population and it may have the best score in
	// one of dimensions. We want to then know what contribution it made to the hypervolume.

	// For every dimension we want to remember two values: the first best and the second best.
	// Best in this case is always an increasing floating point value. If this is a minimization
	// problem, the positive value is determined by subtracting it from the reference point, rather
	// than adding it to the reference point.
	var dimensions int = len(s.ReferencePoint)
	var maxes []dimensionMax
	for i := 0; i < dimensions; i++ {
		maxes = append(maxes, newDimensionMax(s.ReferencePoint[i], s.IsMaximize[i]))
	}

	// For each member of the population, determine if it contributes to the whole population's hypercube
	// or the whole population's second-best hypercube.
	// for _, specimen := range specimens {
	//
	// }

	// The specimens are sorted from fittest to least fit.
	return nil
}

// stretchDimensions runs all the outcomes from a single specimen through the maxes, stretching each in turn.
func stretchDimensions(maxes []dimensionMax, outcomes []float64) []dimensionMax {

	// The dimensions must be equal.
	if len(maxes) != len(outcomes) {
		log.Panicf("stretchDimensions expects %d dimensions, but outcomes have %d dimensions", len(maxes), len(outcomes))
	}

	for i := 0; i < len(outcomes); i++ {
		maxes[i].stretch(outcomes[i])
	}
	return maxes
}

// // LoadHypervolumeIndicatorSelectorConfig loads the json filename as a new configuration.
// // This selector scores and sorts specimens and uses a second selector to determine how
// // to pick a member of the population.
// func LoadHypervolumeIndicatorSelectorConfig(selector Selector, filename string) (HypervolumeIndicatorSelector, error) {
// 	var err error
// 	var bytes []byte
// 	var selector TruncateSelector
//
// 	log.Printf("Loading truncate selector Config: '%s'\n", filename)
//
// 	// Load and parse from json.
// 	if bytes, err = ioutil.ReadFile(filename); err != nil {
// 		return TruncateSelector{}, err
// 	}
// 	if err = json.Unmarshal(bytes, &selector); err != nil {
// 		return TruncateSelector{}, err
// 	}
//
// 	return selector, error(nil)
// }
