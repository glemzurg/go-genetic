package genetic

import ()

// calcualteHypervolumeIndicators computes all the indicators for each specimen in a population.
//
// A hypervolume indicator is the volume of a hypercube that is unique to the hypercube in the population. No other hypercubes are inside that volume.
// If you imagine the overlaying of all the hypercubes in a population as a jagged terrain, the hypervolume indicators are all the jaggs, each belonging
// to a specific hypercube (or in this case specimen's hypercube).
//
// Each specimen's hypercube represents how far that specimen has pushed the multiple outcomes being testing, each outcome value being a dimension of the hypercube.
// The hypervolume indicator, or the jagg on the jagged terrain, represents the contribution this specimen has made to the population as a whole, what it has discovered
// in solving the problem that no other specimen has.
//
// If this specimen's hypercube is wholely consumed inside another specimen's hypercube, the specimen is "dominated", represents no jagg on the terrain, and is not meaningfully
// stretching the population's solution to the problem at hand.
//
// The specimen hypercubes are defined by two points: (0.0, 0.0, ...) and the defining point for the opposite corner (all positive values). Every hypercube shares the base
// point of (0.0, 0.0, ...) so the hypervolume indicator is found by keeping the defining point and then moving the base point closer to it. We move the base point closer to
// the defining point whenever we find another hypercube that has "consumed" part of this shrinking hypercube's volume. After going through all the other members of the population,
// our remaining volume is the hypervolume indicator. If the remaining volume has become zero, this hypercube is dominated.
func calcualteHypervolumeIndicators(hypercubes []*specimenHypercube) {

	// // Sort the hypercubes into a k-d tree for quicker searching.
	// var kdTree hypercubeKdTree = newHypecubeKdTree(hypercubes)
	//
	// // Calculate the hypervolume indicator for each hypercube.
	// for _, hypercube := range hypercubes {
	//
	// 	// Is ths hypercube not dominated? We're working with pointers so working with another hypercube may have discovered this one was dominated.
	// 	if !hypercube.isDominated {
	//
	// 		// This hypercube is not dominated, let's calculuate its hypervolume indicator.
	// 		kdTree.calculateHypervolumeIndicator(hypercube)
	// 	}
	// }
}
