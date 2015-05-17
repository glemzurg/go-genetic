package genetic

import (
	"fmt"
	"math"
	"sort"
)

const (
	// The types of genes that can exist.
	_GENE_TYPE_CONNECTION = "connection"
	_GENE_TYPE_NODE       = "node"
)

// NeatGenome is the genome of a NEAT CPPN.
type NeatGenome struct {
	Genes []NeatGene
}

// Clone makes a copy of this genome ensuring no data structure are shared.
func (g *NeatGenome) Clone() (clone NeatGenome) {
	clone = NeatGenome{}
	clone.Genes = make([]NeatGene, len(g.Genes))
	copy(clone.Genes, g.Genes)
	return clone
}

// NeatGene is a single gene in a NeatGenome
type NeatGene struct {
	GeneId    uint64  // The unique (in an experiment) identity of this gene, shared by eventually many CPPNs.
	IsEnabled bool    // Genes can be disabled, but need to remain in order to compare ancestry of specimens.
	Type      string  // The type of gene this is.
	From      string  // Describing the source of a connection.
	To        string  // Describing the sink of a connection.
	Weight    float64 // A value between 0.0 and 1.0.
	Function  string  // The activation function for node genes.
}

// ByGeneId implements sort.Interface to sort ascending by GeneId.
// Example: sort.Sort(ByGeneId(genes))
type ByGeneId []NeatGene

func (a ByGeneId) Len() int           { return len(a) }
func (a ByGeneId) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByGeneId) Less(i, j int) bool { return a[i].GeneId < a[j].GeneId }

// SpeciationDistance computes how related to NEAT genomes are (i.e. are they the same species?).
// The lower the distance, the more alike the genomes are and the closer they are to being the same
// genome. Three constants C1, C2, C3 are used to configure what a particular experiment identifies
// as important for determining species.
//
// A high configuration C1 gives more importance to excess genes (the tail of the longer genome).
// A high configuration C2 gives more importance to disjoint genes (non-shared genes in either genome before the excess genes).
// A high configuration C3 gives more importance to differences in shared genes.
//
// SpeciationDistances = C1 * (ExcessGeneCount / LargestGeneCount) + C2 * (DisjointGeneCount / LargestGeneCount) + C3 * AverageWeightDiffOfSharedGenes
func SpeciationDistance(genomeA NeatGenome, genomeB NeatGenome, c1 float64, c2 float64, c3 float64) float64 {
	// Run a few sanity checks to ensure the code is working correctly.
	var geneCountBefore int
	var geneCountAfter int

	// The speciation distance is calculated by:
	//
	//   = C1 * (ExcessGeneCount   / LargestGeneCount)
	//   + C2 * (DisjointGeneCount / LargestGeneCount)
	//   + C3 * AverageWeightDiffOfSharedGenes
	//
	// C1, C2, and C3 are arbitrary values passed in to weight the results in different ways.

	// What is the longest gene count between the genomes?
	var longestGeneCount int = len(genomeA.Genes)
	if len(genomeB.Genes) > longestGeneCount {
		longestGeneCount = len(genomeB.Genes)
	}

	// Figure out which is the "longest" genome, the one with the most recent gene.
	var lastA NeatGene = genomeA.Genes[len(genomeA.Genes)-1]
	var lastB NeatGene = genomeB.Genes[len(genomeB.Genes)-1]

	// Make one the older and one the younger genome.
	// The younger genome is the one with more recent changes (higher gene ids).
	var youngerGenome NeatGenome = genomeA
	var olderGenome NeatGenome = genomeB
	if lastB.GeneId > lastA.GeneId {
		youngerGenome = genomeB
		olderGenome = genomeA
	}

	// From here on, we'll just be working with the gene slices.
	// Make a copy of the genes so we don't alter the original genomes.
	var youngerGenes []NeatGene = make([]NeatGene, len(youngerGenome.Genes))
	copy(youngerGenes, youngerGenome.Genes)
	var olderGenes []NeatGene = make([]NeatGene, len(olderGenome.Genes))
	copy(olderGenes, olderGenome.Genes)

	// Sanity check.
	geneCountBefore = len(youngerGenes)

	// Prune any excess genes from the younger genome, all genes more recent than the older genome.
	var prunedYoungerGenes []NeatGene = youngerGenes // Assume no excess initially.
	var excessGenes []NeatGene
	// It's possible there are none. The two genomes may end on the same gene.
	var highestOlderGene NeatGene = olderGenes[len(olderGenes)-1]
	var highestYoungerGene NeatGene = youngerGenes[len(youngerGenes)-1]
	if highestYoungerGene.GeneId > highestOlderGene.GeneId {
		// Find the first index of the younger genes higher than this gene id.
		var higherIndex int = sort.Search(len(youngerGenes), func(i int) bool { return youngerGenes[i].GeneId > highestOlderGene.GeneId })
		if higherIndex < len(youngerGenes) && youngerGenes[higherIndex].GeneId > highestOlderGene.GeneId {
			// We found the higher index. The excess genes (the tail of the younger genes) is this
			// index to the end of the younger gene.
			prunedYoungerGenes = youngerGenes[:higherIndex]
			excessGenes = youngerGenes[higherIndex:]
		}
	}
	var excessGeneCount int = len(excessGenes)

	// Sanity check.
	geneCountAfter = len(prunedYoungerGenes) + excessGeneCount
	if geneCountBefore != geneCountAfter {
		panic(fmt.Errorf("sanity check failed: %d != %d", geneCountBefore, geneCountAfter))
	}

	// Sanity check.
	geneCountBefore = len(prunedYoungerGenes) + len(olderGenes)

	// Analyze the genes that could be in both older and younger for disjoint and average weigth.
	var olderAndYoungerGenes []NeatGene
	olderAndYoungerGenes = append(olderAndYoungerGenes, prunedYoungerGenes...)
	olderAndYoungerGenes = append(olderAndYoungerGenes, olderGenes...)
	// Build the counts we need to calculate the speciation distance.
	var disjointGeneCount int
	var weightSum float64
	var weightContributors int
	// Sort and loop through all the genes.
	sort.Sort(ByGeneId(olderAndYoungerGenes))
	var lastGeneIndex int = len(olderAndYoungerGenes) - 1
	for i := 0; i < len(olderAndYoungerGenes); i++ {
		// Is this gene solo (a dijoint gene), or back-to-back with a gene of the same id (a shared gene)?
		// If this is the last gene index, it is a disjoint gene. No more genes can follow it.
		if i == lastGeneIndex {
			disjointGeneCount++
		} else {
			// There are genes following this gene. Take a peek.
			var thisGene NeatGene = olderAndYoungerGenes[i]
			var nextGene NeatGene = olderAndYoungerGenes[i+1]
			// Are they the same gene?
			if thisGene.GeneId == nextGene.GeneId {

				// These two genes are shared. Compute the difference in weigth between them.
				weightSum += math.Abs(thisGene.Weight - nextGene.Weight)
				weightContributors++

				// We've just "consumed" two genes instead of one so indicate one more
				// gene has been handled.
				i++

			} else {
				// Not the same gene. This is a disjoint gene.
				disjointGeneCount++
			}
		}
	}

	// Sanity check.
	geneCountAfter = disjointGeneCount + 2*weightContributors
	if geneCountBefore != geneCountAfter {
		panic(fmt.Errorf("sanity check failed: %d != %d", geneCountBefore, geneCountAfter))
	}

	// Calculate the numbers-of-interest for the equation.
	var averageWeightDiffOfSharedGenes float64
	if weightContributors > 0 {
		averageWeightDiffOfSharedGenes = weightSum / float64(weightContributors)
	}

	// The speciation distance itself.
	var speciationDistance float64 = c1*(float64(excessGeneCount)/float64(longestGeneCount)) + c2*(float64(disjointGeneCount)/float64(longestGeneCount)) + c3*averageWeightDiffOfSharedGenes
	return speciationDistance
}

// IsSameSpecies wraps the speciation distance calculation with a threshold value that is the breaking point
// for when the distance is too large to be in the same species. If the threshold is 0.0, then all genomes are
// expected to be part of one big species in the population (the feature is "turned off").
func IsSameSpecies(genomeA NeatGenome, genomeB NeatGenome, config SpeciationConfig) (isSameSpecies bool, speciationDistance float64) {
	speciationDistance = SpeciationDistance(genomeA, genomeB, config.C1, config.C2, config.C3)
	isSameSpecies = (config.Threshold == 0.0 || speciationDistance <= config.Threshold)
	return isSameSpecies, speciationDistance
}
