package genetic

import (
	// "encoding/json"
	// "io/ioutil"
	"log"
	"sort"
)

// HypervolumeIndicator evaluates multi-outcome scores and scores each on how much it stretches
// the whole population's improvement by looking at all the specimens as overlapping cubes (hyper cubes).
type HypervolumeIndicatorContextScorer struct {
	ReferencePoint []float64 // A cube is defined by two multi-demensional points which represent opposite corners of the cube. To examine the mult-value outcomes as comparable cubes, they all share one arbitary point as the opposite corner.
	IsMaximize     []bool    // For each dimension, are we trying to maximize or minimize the value?
	Weights        []float64 // For each dimension, do some contribute more to the score than others?
}

// hypercubeDimension is the tuple of best and second best points in a dimension.
type hypercubeDimension struct {
	first      float64 // The value that most stretches the hypercube volume in this dimension.
	second     float64 // The value that second most stretches the hypercube volume in this dimension.
	base       float64 // The reference point's value in this dimension.
	isMaximize bool    // If true, values subtract the base to determine their contribution. If false, the base subtracts the value to determine the contribution.
	weight     float64 // How much does this dimension contribute to the final hypervolume indicator score.
}

// newHypercubeDimension creates a well-formed dimensionMax.
func newHypercubeDimension(base float64, isMaximize bool, weight float64) hypercubeDimension {
	return hypercubeDimension{
		first:      base,       // All values allowed must be on one side of the reference point value.
		second:     base,       // All values allowed must be on one side of the reference point value.
		base:       base,       // What is the base value for this dimension?
		isMaximize: isMaximize, // Determines which side of of the reference point value is valid side.
		weight:     weight,     // How much does this dimension contribute to the final hypervolume indicator score.
	}
}

// stretch pushes the first or second best dimension max farther if the value shoud do so, or leaves the max unchanged.
func (d *hypercubeDimension) stretch(value float64) {
	// If value is better than the fist value, knock the first value to the second place, ejecting that value.
	// If value is worse than the first and better than teh second, knock the second place value from the struct.
	// If the value is neither, leave the struct unchanged.
	if d.isMaximize {
		switch {
		case value > d.first:
			d.second = d.first
			d.first = value
		case value > d.second:
			d.second = value
		}
	} else {
		switch {
		case value < d.first:
			d.second = d.first
			d.first = value
		case value < d.second:
			d.second = value
		}
	}
}

// hypercubeContribution is the contribution of a given specimen to the population's hypercube
type hypercubeContribution struct {
	indicator         float64  // The amount by which this specimen stretched the population's hypercube.
	volume            float64  // The hypercube volume of the specimen itself, some sub-cube withint the population's cube.
	weightedIndicator float64  // Indicator inversely weighted by number of members of the specimen's species.
	weightedVolume    float64  // Volume inversely weighted by number of members of the specimen's species.
	specimen          Specimen // The specimen in question.
}

// ByHypervolumeIndicator implements sort.Interface to sort descending by hypervolume indicator.
// Example: sort.Sort(byHypervolumeIndicator(contributions))
type byHypervolumeIndicator []hypercubeContribution

func (a byHypervolumeIndicator) Len() int      { return len(a) }
func (a byHypervolumeIndicator) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byHypervolumeIndicator) Less(i, j int) bool {
	// If both contributions have the same indicator then use specimen volume instead.
	if a[i].weightedIndicator == a[j].weightedIndicator {
		return a[i].weightedVolume > a[j].weightedVolume
	}
	// If the indicators are different, use them.
	return a[i].weightedIndicator > a[j].weightedIndicator
}

// newHypercubeContribution calculates the contribution a specimen makes to the population's hypercube.
func newHypercubeContribution(hypercube []hypercubeDimension, specimen Specimen) hypercubeContribution {
	var indicator float64
	var volume float64
	indicator, volume = calculateHypercubeContribution(hypercube, specimen.Outcomes)
	return hypercubeContribution{
		indicator:         indicator,
		volume:            volume,
		weightedIndicator: indicator / float64(specimen.SpeciesMemberCount),
		weightedVolume:    volume / float64(specimen.SpeciesMemberCount),
		specimen:          specimen,
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
	var hypercube []hypercubeDimension
	for i := 0; i < len(s.ReferencePoint); i++ {
		hypercube = append(hypercube, newHypercubeDimension(s.ReferencePoint[i], s.IsMaximize[i], s.Weights[i]))
	}

	// For each member of the population, determine if it contributes to the whole population's hypercube
	// or the whole population's second-best hypercube.
	for _, specimen := range specimens {
		hypercube = stretchDimensions(hypercube, specimen.Outcomes)
	}

	// Update the specimen score to be the score calculated by each specimens contribution to the
	// who population's hypercube.
	var contributions []hypercubeContribution
	for _, specimen := range specimens {
		contributions = append(contributions, newHypercubeContribution(hypercube, specimen))
	}

	// Sort the specimens descending by their contribution.
	sort.Sort(byHypervolumeIndicator(contributions))

	// Dump the specimens into their sort.
	var sortedSpecimens []Specimen
	for _, contribution := range contributions {
		sortedSpecimens = append(sortedSpecimens, contribution.specimen)
	}

	// The specimens are sorted from fittest to least fit.
	return sortedSpecimens
}

// stretchDimensions runs all the outcomes from a single specimen through the maxes, stretching each in turn.
func stretchDimensions(hypercube []hypercubeDimension, outcomes []float64) []hypercubeDimension {

	// The dimensions must be equal.
	if len(hypercube) != len(outcomes) {
		log.Panicf("stretchDimensions expects %d dimensions, but outcomes have %d dimensions", len(hypercube), len(outcomes))
	}

	for i := 0; i < len(outcomes); i++ {
		hypercube[i].stretch(outcomes[i])
	}
	return hypercube
}

// calculateHypercubeContribution determines how a single specimen's multi-outcome compares to
// the population's hypercube.
func calculateHypercubeContribution(hypercube []hypercubeDimension, outcomes []float64) (indicator float64, volume float64) {
	// Build up the indicator and volume one dimension at a time.
	indicator = 0.0
	volume = 1.0
	for i, dimension := range hypercube {

		// What is the outcome for this dimension.
		var outcome float64 = outcomes[i]

		// What are the lengths that contribute to this dimension.
		var indicatorLength float64
		var volumeLength float64

		// Is this dimension counting up or down?
		if dimension.isMaximize {

			// Is this outcome the defining limit of the hypercube in this dimension?
			if outcome == dimension.first {
				// The contribution it makes is the difference to the second highest value.
				indicatorLength = dimension.first - dimension.second // First value is higher number.
			}

			// Regardless of the indicator, what is the contribution to the specimen's hypercube volume.
			if outcome > dimension.base {
				volumeLength = outcome - dimension.base // Outcome is higher.
			} else {
				// If the specimen doesn't pass the referecnce point in this dimension, it has no contribution in *any* dimension.
				// As in, if we let an invalid dimension have a value of 1.0, it may beat a valid dimension that calculates
				// to some value below 1.0 and we can't allow that.
				volumeLength = 0.0
				log.Printf("WARNING: In dimension %d, specimen has outcome %f which is not greater than reference point value %f. Too many of these and specimens will not be meaningfully sorted.", i, outcome, dimension.base)
			}

		} else {

			// We are minimizing on this dimension.

			// Is this outcome the defining limit of the hypercube in this dimension?
			if outcome == dimension.first {
				// The contribution it makes is the difference to the second highest value.
				indicatorLength = dimension.second - dimension.first // First value is lower number.
			}

			// Regardless of the indicator, what is the contribution to the specimen's hypercube volume.
			if outcome < dimension.base {
				volumeLength = dimension.base - outcome // base is higher.
			} else {
				// If the specimen doesn't pass the referecnce point in this dimension, it has no contribution in *any* dimension.
				// As in, if we let an invalid dimension have a value of 1.0, it may beat a valid dimension that calculates
				// to some value below 1.0 and we can't allow that.
				volumeLength = 0.0
				log.Printf("WARNING: In dimension %d, specimen has outcome %f which is not less than reference point value %f. Too many of these and specimens will not be meaningfully sorted.", i, outcome, dimension.base)
			}
		}

		// Did this dimension contribute to the population's hypercube?
		if indicatorLength > 0.0 {
			// If the indicator has no value yet, start it.
			if indicator == 0.0 {
				indicator = 1.0 // The indicator is now active.
			}
			// Add this contribution.
			indicator *= indicatorLength * dimension.weight // Some dimensions are more important than others.
		}

		// Regardless, build the volume information for this hypercube.
		volume *= volumeLength * dimension.weight // Some dimensions are more important than others.
	}

	return indicator, volume
}

// // LoadHypervolumeIndicatorSelectorConfig loads the json filename as a new configuration.
// // This selector scores and sorts specimens and uses a second selector to determine how
// // to pick a member of the population.
// func LoadHypervolumeIndicatorSelectorConfig(selector Selector, filename string) (HypervolumeIndicatorSelector, error) {
// 	var err error
// 	var bytes []byte
// 	var selector SelectorElitism
//
// 	log.Printf("Loading truncate selector Config: '%s'\n", filename)
//
// 	// Load and parse from json.
// 	if bytes, err = ioutil.ReadFile(filename); err != nil {
// 		return SelectorElitism{}, err
// 	}
// 	if err = json.Unmarshal(bytes, &selector); err != nil {
// 		return SelectorElitism{}, err
// 	}
//
// 	return selector, error(nil)
// }
