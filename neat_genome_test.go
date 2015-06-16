package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
)

// Create a suite.
type neatGenomeSuite struct{}

var _ = Suite(&neatGenomeSuite{})

// Add the tests.

func (s *neatGenomeSuite) Test_CalculateSpeciationDistance(c *C) {
	var genomeA, genomeB neatGenome
	var expectedDistance float64

	// SpeciationDistance = C1  * (Excess/Longest) + C2  * (Disjoint/Longest) + C3  * AverageWeightDiff

	// Excess genes.

	// no excess gene
	expectedDistance = 1.0*(0.0/1.0) + 0.0 + 0.0
	genomeA = gnm([]gn{gn{1, 0.7}})
	genomeB = gnm([]gn{gn{1, 0.2}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 1.0, 0.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 1.0, 0.0, 0.0), Equals, expectedDistance)

	// no excess gene, but disjoint genes
	expectedDistance = 1.0*(0.0/2.0) + 0.0 + 0.0
	genomeA = gnm([]gn{gn{1, 0.2}, gn{2, 0.2}})
	genomeB = gnm([]gn{gn{0, 0.0}, gn{2, 0.7}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 1.0, 0.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 1.0, 0.0, 0.0), Equals, expectedDistance)

	// one excess gene
	expectedDistance = 1.0*(1.0/2.0) + 0.0 + 0.0
	genomeA = gnm([]gn{gn{1, 0.2}})
	genomeB = gnm([]gn{gn{1, 0.7}, gn{2, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 1.0, 0.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 1.0, 0.0, 0.0), Equals, expectedDistance)

	// one excess gene, altering C1
	expectedDistance = 0.5*(1.0/2.0) + 0.0 + 0.0
	genomeA = gnm([]gn{gn{1, 0.2}})
	genomeB = gnm([]gn{gn{1, 0.7}, gn{2, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.5, 0.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.5, 0.0, 0.0), Equals, expectedDistance)

	// one excess gene, following single disjoint in other genome
	expectedDistance = 1.0*(1.0/1.0) + 0.0 + 0.0
	genomeA = gnm([]gn{gn{1, 0.2}})
	genomeB = gnm([]gn{gn{0, 0.0}, gn{2, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 1.0, 0.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 1.0, 0.0, 0.0), Equals, expectedDistance)

	// three excess genes
	expectedDistance = 1.0*(3.0/4.0) + 0.0 + 0.0
	genomeA = gnm([]gn{gn{1, 0.2}})
	genomeB = gnm([]gn{gn{1, 0.7}, gn{2, 0.4}, gn{7, 0.4}, gn{9, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 1.0, 0.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 1.0, 0.0, 0.0), Equals, expectedDistance)

	// Disjoint genes.

	// no disjoint gene
	expectedDistance = 0.0 + 1.0*(0.0/1.0) + 0.0
	genomeA = gnm([]gn{gn{1, 0.2}})
	genomeB = gnm([]gn{gn{1, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 1.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 1.0, 0.0), Equals, expectedDistance)

	// no disjoint gene, but excess gene
	expectedDistance = 0.0 + 1.0*(0.0/2.0) + 0.0
	genomeA = gnm([]gn{gn{1, 0.2}})
	genomeB = gnm([]gn{gn{1, 0.4}, gn{9, 0.1}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 1.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 1.0, 0.0), Equals, expectedDistance)

	// one disjoint gene
	expectedDistance = 0.0 + 1.0*(1.0/2.0) + 0.0
	genomeA = gnm([]gn{gn{0, 0.0}, gn{2, 0.2}})
	genomeB = gnm([]gn{gn{1, 0.7}, gn{2, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 1.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 1.0, 0.0), Equals, expectedDistance)

	// one disjoint gene, altering C2
	expectedDistance = 0.0 + 0.5*(1.0/2.0) + 0.0
	genomeA = gnm([]gn{gn{0, 0.0}, gn{2, 0.2}})
	genomeB = gnm([]gn{gn{1, 0.7}, gn{2, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 0.5, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 0.5, 0.0), Equals, expectedDistance)

	// one disjoint gene, followed by excess on other gene
	expectedDistance = 0.0 + 1.0*(1.0/1.0) + 0.0
	genomeA = gnm([]gn{gn{0, 0.0}, gn{2, 0.2}})
	genomeB = gnm([]gn{gn{1, 0.7}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 1.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 1.0, 0.0), Equals, expectedDistance)

	// one disjoint gene, followed by shared then excess in same genome
	expectedDistance = 0.0 + 1.0*(1.0/3.0) + 0.0
	genomeA = gnm([]gn{gn{0, 0.0}, gn{2, 0.2}})
	genomeB = gnm([]gn{gn{1, 0.7}, gn{2, 0.4}, gn{3, 0.1}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 1.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 1.0, 0.0), Equals, expectedDistance)

	// one disjoint gene, followed by shared then excess in other genome
	expectedDistance = 0.0 + 1.0*(1.0/2.0) + 0.0
	genomeA = gnm([]gn{gn{0, 0.0}, gn{2, 0.2}, gn{3, 0.1}})
	genomeB = gnm([]gn{gn{1, 0.7}, gn{2, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 1.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 1.0, 0.0), Equals, expectedDistance)

	// one disjoint gene, following shared in genome then followed by excess in other genome
	expectedDistance = 0.0 + 1.0*(1.0/2.0) + 0.0
	genomeA = gnm([]gn{gn{1, 0.2}, gn{0, 0.0}, gn{3, 0.1}})
	genomeB = gnm([]gn{gn{1, 0.7}, gn{2, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 1.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 1.0, 0.0), Equals, expectedDistance)

	// two disjoint genes, before shared
	expectedDistance = 0.0 + 1.0*(2.0/3.0) + 0.0
	genomeA = gnm([]gn{gn{0, 0.0}, gn{0, 0.0}, gn{3, 0.1}})
	genomeB = gnm([]gn{gn{1, 0.7}, gn{2, 0.4}, gn{3, 0.2}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 1.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 1.0, 0.0), Equals, expectedDistance)

	// two disjoint genes, before shared and alternating
	expectedDistance = 0.0 + 1.0*(2.0/2.0) + 0.0
	genomeA = gnm([]gn{gn{0, 0.0}, gn{2, 0.4}, gn{3, 0.1}})
	genomeB = gnm([]gn{gn{1, 0.7}, gn{0, 0.0}, gn{3, 0.2}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 1.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 1.0, 0.0), Equals, expectedDistance)

	// two disjoint genes, before excess and alternating
	expectedDistance = 0.0 + 1.0*(2.0/2.0) + 0.0
	genomeA = gnm([]gn{gn{0, 0.0}, gn{2, 0.4}})
	genomeB = gnm([]gn{gn{1, 0.7}, gn{0, 0.0}, gn{3, 0.2}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 1.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 1.0, 0.0), Equals, expectedDistance)

	// three disjoint genes
	expectedDistance = 0.0 + 1.0*(3.0/4.0) + 0.0
	genomeA = gnm([]gn{gn{0, 0.0}, gn{0, 0.0}, gn{0, 0.0}, gn{9, 0.2}})
	genomeB = gnm([]gn{gn{1, 0.7}, gn{2, 0.4}, gn{7, 0.4}, gn{9, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 1.0, 0.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 1.0, 0.0), Equals, expectedDistance)

	// Average weight differences between shared genes.

	// average weight diff, no shared gene
	expectedDistance = 0.0 + 0.0 + 1.0*0.0
	genomeA = gnm([]gn{gn{1, 0.7}})
	genomeB = gnm([]gn{gn{0, 0.0}, gn{2, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 0.0, 1.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 0.0, 1.0), Equals, expectedDistance)

	// average weight diff, one gene
	expectedDistance = 0.0 + 0.0 + 1.0*0.5 // Ave of 0.5
	genomeA = gnm([]gn{gn{1, 0.9}})
	genomeB = gnm([]gn{gn{1, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 0.0, 1.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 0.0, 1.0), Equals, expectedDistance)

	// average weight diff, one gene, C3 modified
	expectedDistance = 0.0 + 0.0 + 0.5*0.5 // Ave of 0.5
	genomeA = gnm([]gn{gn{1, 0.9}})
	genomeB = gnm([]gn{gn{1, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 0.0, 0.5), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 0.0, 0.5), Equals, expectedDistance)

	// average weight diff, two genes
	expectedDistance = 0.0 + 0.0 + 1.0*(0.5+0.1)/2.0 // Ave of 0.5, 0.1
	genomeA = gnm([]gn{gn{1, 0.9}, gn{2, 0.4}})
	genomeB = gnm([]gn{gn{1, 0.4}, gn{2, 0.5}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 0.0, 1.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 0.0, 1.0), Equals, expectedDistance)

	// average weight diff, many genes with disjoint genes mixed in
	expectedDistance = 0.0 + 0.0 + 1.0*(0.5+0.1)/2.0 // Ave of 0.5, 0.1
	genomeA = gnm([]gn{gn{1, 0.4}, gn{2, 0.4}, gn{0, 0.0}, gn{4, 0.4}, gn{5, 0.5}, gn{6, 0.5}})
	genomeB = gnm([]gn{gn{0, 0.0}, gn{2, 0.9}, gn{3, 0.4}, gn{0, 0.0}, gn{5, 0.4}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 0.0, 0.0, 1.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 0.0, 0.0, 1.0), Equals, expectedDistance)

	// Bring it all together.

	// shared, disjoint, and excess genes. All constants set.
	expectedDistance = 1.0*(1.0/3.0) + 2.0*(1.0/3.0) + 3.0*0.3 // Ave of 0.5, -0.1
	genomeA = gnm([]gn{gn{1, 0.75}, gn{2, 0.4}, gn{3, 0.3}})
	genomeB = gnm([]gn{gn{1, 0.25}, gn{2, 0.5}, gn{0, 0.0}, gn{4, 0.1}})
	c.Assert(calculateSpeciationDistance(genomeA, genomeB, 1.0, 2.0, 3.0), Equals, expectedDistance)
	c.Assert(calculateSpeciationDistance(genomeB, genomeA, 1.0, 2.0, 3.0), Equals, expectedDistance)
}

func (s *neatGenomeSuite) Test_IsSameSpecies(c *C) {

	// Create a sample configuration.
	var config SpeciationConfig = SpeciationConfig{
		Threshold: 0.0,
		C1:        1.0, // Emphasies the longer tail of the younger genome.
		C2:        0.0,
		C3:        0.0,
	}

	// Verify that we see .
	var isSpecies bool
	var speciationDistance float64

	// We have two genomes that compute to a speciation distance of 1.0 (with the constants we're using)
	var genomeA neatGenome = gnm([]gn{gn{1, 0.2}})
	var genomeB neatGenome = gnm([]gn{gn{0, 0.0}, gn{2, 0.4}})

	// Set the threshold so they are in the same species.
	config.Threshold = 1.0
	isSpecies, speciationDistance = isSameSpecies(genomeA, genomeB, config)
	c.Check(isSpecies, Equals, true)
	c.Assert(speciationDistance, Equals, 1.0)
	// Reverse parameters.
	isSpecies, speciationDistance = isSameSpecies(genomeB, genomeA, config)
	c.Check(isSpecies, Equals, true)
	c.Assert(speciationDistance, Equals, 1.0)

	// Set the threshold so they are NOT in the same species.
	config.Threshold = 0.9
	isSpecies, speciationDistance = isSameSpecies(genomeA, genomeB, config)
	c.Check(isSpecies, Equals, false)
	c.Assert(speciationDistance, Equals, 1.0)
	// Reverse parameters.
	isSpecies, speciationDistance = isSameSpecies(genomeB, genomeA, config)
	c.Check(isSpecies, Equals, false)
	c.Assert(speciationDistance, Equals, 1.0)

	// But if the threshold is 0.0, then the feature is "turned off" and everything is part
	// of the same big species that comprises the whole population.
	config.Threshold = 0.0
	isSpecies, speciationDistance = isSameSpecies(genomeA, genomeB, config)
	c.Check(isSpecies, Equals, true)
	c.Assert(speciationDistance, Equals, 1.0)
	// Reverse parameters.
	isSpecies, speciationDistance = isSameSpecies(genomeB, genomeA, config)
	c.Check(isSpecies, Equals, true)
	c.Assert(speciationDistance, Equals, 1.0)

}

// To keep the tests easy to read, have a minimal way to show the genome.
// Only include the data needed for this test.
type gn struct {
	GeneId uint64 // Set to 0 if it is supposed to be ignored. Allows consistent visual formatting for tests.
	Weight float64
}

// Create a genome for testing from minimal genes.
func gnm(minimalGenes []gn) (genome neatGenome) {
	// Only create the data meaningful for calculating speciation distance.
	for _, minimalGene := range minimalGenes {
		if minimalGene.GeneId != 0 {
			genome.Genes = append(genome.Genes, neatGene{GeneId: minimalGene.GeneId, Weight: minimalGene.Weight})
		}
	}
	return genome
}
