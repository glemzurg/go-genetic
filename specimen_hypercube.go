package genetic

import (
	"log"
)

// specimenHypercube is the hypercube defined by the multi-outcome score for a specimen, normalized to be
// positive points relative to 0.0 in all dimensions. In an actual experiment, the experiment has a reference
// point as the base and the individual dimesnions may be either maximizing or minimizing. The normalization
// provides a clean way to sort and compare hypercubes for sorting by hypervolume indicator.
type specimenHypercube struct {

	// Some values are from the specimen alone.
	dimensions []float64 // The value on each dimension of the multi-outcome.
	volume     float64   // The hypercube volume for this specimen, 0.0 if any dimension incalculable based on reference point.
	specimen   Specimen  // The specimen in question.

	// Some values are from the specimen's relationship to the rest of the population.
	isDominated   bool      // True if this cube lies completly within another cube.
	indicator     float64   // The volume of the part of this hypercube that is not overlapped by any other hypercube in the population.
	indicatorBase []float64 // The corner of the indicator hypercube opposite the corner defined by this hypercube's dimensions.

	// Bookkeeping for doing comparisons.
	comparedWith map[*specimenHypercube]bool // True if we have compared against a cube already.
}

// newSpecimenHypercube creates a new normalized hypercube for the specimen.
func newSpecimenHypercube(specimen Specimen, referencePoint []float64, isMaximize []bool, weights []float64) specimenHypercube {

	// Normalize the specimens multi-outcome to just be lengths from the reference point.
	var dimensions []float64
	var volume float64 = 1.0
	dimensions, volume = specimenHypercubeDimensions(specimen.Outcomes, referencePoint, isMaximize, weights)

	// All cubes start without knowing what the base for the hypervolume indicator is.
	// Set it to (0.0, 0.0, ...)
	var indicatorBase []float64
	for i := 0; i < len(referencePoint); i++ {
		indicatorBase = append(indicatorBase, 0.0)
	}

	return specimenHypercube{
		dimensions:    dimensions,
		volume:        volume,
		specimen:      specimen,
		indicatorBase: indicatorBase,
	}
}

// isDominatedBy lets use know if this hypercube wholely exists inside another hypercube of the population.
func (h *specimenHypercube) isDominatedBy(other *specimenHypercube) bool {
	// If any of our dimensions are greater than the other's, we are *not* dominated.
	for i := range h.dimensions {
		if h.dimensions[i] > other.dimensions[i] {
			return false
		}
	}
	return true
}

// equals tells of two hypercubes have the exact same dimensions.
func (h *specimenHypercube) equals(other *specimenHypercube) bool {
	// If any of our dimensions don't match, we are *not* equal.
	for i := range h.dimensions {
		if h.dimensions[i] != other.dimensions[i] {
			return false
		}
	}
	return true
}

// setDominated updates this cube to the state that it's dominated. This cube lies wholely within another cube.
func (h *specimenHypercube) setDominated() {
	h.isDominated = true
	h.indicator = 0.0
	h.indicatorBase = nil
}

// calculateHypervolumeIndicator calculates the indicator if we are not dominated.
func (h *specimenHypercube) calculateHypervolumeIndicator() {
	if !h.isDominated {
		h.indicator = calculateHypervolume(h.dimensions, h.indicatorBase)
	}
}

// specimenHypercubeDimensions calculates the normalized hypercube.
func specimenHypercubeDimensions(outcomes []float64, referencePoint []float64, isMaximize []bool, weights []float64) (dimensions []float64, volume float64) {

	// The dimension count must be equal.
	if len(outcomes) != len(referencePoint) {
		log.Panicf("specimenHypercubeDimensions expects %d dimensions, but multi-outcome has %d dimensions", len(referencePoint), len(outcomes))
	}

	// Start the volume.
	volume = 1.0

	// Normalize the specimens multi-outcome to just be lengths from the reference point.
	for i, base := range referencePoint {

		// The specimen's value.
		var outcome float64 = outcomes[i]

		// How important this particular dimension is to the hypercube?
		// This is from the experiment configuration making some axis more important than others.
		var weight float64 = weights[i]

		// Compute the length, which will become the normalized point in this dimension relative to 0.0.
		var length float64

		// Is this particular dimension counting up or down?
		if isMaximize[i] {
			// We're counting up in this dimension.
			// Are we greater than the base?
			if outcome > base {
				length = (outcome - base) * weight // Outcome is greater than base.
			} else {
				// If the specimen doesn't pass the referecnce point in this dimension, it has no contribution in *any* dimension.
				// As in, if we let an invalid dimension have a value of 1.0, it may beat a valid dimension that calculates
				// to some value below 1.0 and we can't allow that.
				length = 0.0
				log.Printf("WARNING: In dimension %d, specimen has outcome %f which is not greater than reference point value %f. Too many of these and specimens will not be meaningfully sorted.", i, outcome, base)
			}
		} else {
			// We're counting down in this dimension.
			// Are we less than the base?
			if outcome < base {
				length = (base - outcome) * weight // Outcome is less than base.
			} else {
				// If the specimen doesn't pass the referecnce point in this dimension, it has no contribution in *any* dimension.
				// As in, if we let an invalid dimension have a value of 1.0, it may beat a valid dimension that calculates
				// to some value below 1.0 and we can't allow that.
				length = 0.0
				log.Printf("WARNING: In dimension %d, specimen has outcome %f which is not less than reference point value %f. Too many of these and specimens will not be meaningfully sorted.", i, outcome, base)
			}
		}

		// Add this point to the normalized points.
		dimensions = append(dimensions, length)

		// Compute the volume. If any of the lengths ever are 0.0, the volume is basically negated (which is correct).
		volume *= length
	}

	return dimensions, volume
}
