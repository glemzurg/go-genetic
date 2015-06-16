package genetic

import (
	. "gopkg.in/check.v1" // https://labix.org/gocheck
	"math/rand"
	"time"
)

// Create a suite.
type PopulationSuite struct{}

var _ = Suite(&PopulationSuite{})

// Add the tests.

func (s *PopulationSuite) Test_Population_AddSpecimen_NoSpeciesYet(c *C) {
	var population generationPopulation
	var expectedPopulation generationPopulation

	// We'll need a few genomes for this test.
	var genomeA neatGenome = gnm([]gn{gn{1, 0.2}})

	// Use a config that will be easy to study withs species.
	var config PopulationConfig = PopulationConfig{
		Speciation: SpeciationConfig{
			Threshold: 0.5, // One longer tail should be enough to indicate the genomes are not the same species.
			C1:        1.0, // Emphasies the longer tail of the younger genome.
			C2:        0.0,
			C3:        0.0,
		},
	}

	// Population is a new population with no species.
	population = newPopulation(config)
	expectedPopulation = generationPopulation{config: config}
	c.Assert(population, DeepEquals, expectedPopulation)

	// Add a specimen, which will create the first species.
	population.AddNeuralNet(NeatNeuralNet{Genome: genomeA}, 10.0, 100.0, []float64{1.0, 2.0})
	expectedPopulation = generationPopulation{
		config: config,
		species: []Species{
			Species{
				genome: genomeA, // Species has the identity genome of the first member.
				Specimens: []Specimen{
					Specimen{
						NeuralNet: NeatNeuralNet{Genome: genomeA},
						Score:     10.0,
						Bonus:     100.0,
						Outcomes:  []float64{1.0, 2.0},
					},
				},
			},
		},
	}
	c.Assert(population, DeepEquals, expectedPopulation)
}

func (s *PopulationSuite) Test_RandomMateMutatePick(c *C) {

	c.Skip("This test has been verified but is unpredictable so should be manually reviewed.")

	// Get the randomness rolling.
	rand.Seed(time.Now().UnixNano())

	// All equally chosen.
	c.Check(randomMateMutatePick(0, 0, 0, 0), DeepEquals, -1) // Can never equal this.
	c.Check(randomMateMutatePick(1, 1, 1, 1), DeepEquals, -1) // Can never equal this.

	// No mating.
	c.Check(randomMateMutatePick(0, 1, 1, 1), DeepEquals, -1) // Can never equal this.

	// Only one choice.
	c.Check(randomMateMutatePick(0, 1, 0, 0), DeepEquals, -1) // Can never equal this.
}

func (s *PopulationSuite) Test_Population_AddSpecimen_MatchingSpecies(c *C) {
	var population generationPopulation
	var expectedPopulation generationPopulation

	// We have two genomes that compute to a speciation distance of 1.0 (with the constants we're using)
	var genomeA neatGenome = gnm([]gn{gn{1, 0.2}})
	var genomeB neatGenome = gnm([]gn{gn{0, 0.0}, gn{2, 0.4}})

	// Use a config that will be easy to study withs species.
	var config PopulationConfig = PopulationConfig{
		Speciation: SpeciationConfig{
			Threshold: 0.5, // One longer tail should be enough to indicate the genomes are not the same species.
			C1:        1.0, // Emphasies the longer tail of the younger genome.
			C2:        0.0,
			C3:        0.0,
		},
	}

	// Population is a new population with no species.
	population = newPopulation(config)
	population.species = append(population.species, Species{genome: genomeB})
	population.species = append(population.species, Species{genome: genomeA}) // The matching genome is second.
	population.species = append(population.species, Species{genome: genomeA}) // Another matching genome is third.

	// Add a specimen, will with match the second species.
	population.AddNeuralNet(NeatNeuralNet{Genome: genomeA}, 0.0, 0.0, nil)
	expectedPopulation = generationPopulation{
		config: config,
		species: []Species{
			Species{
				genome:    genomeB, // The specimen doesn't match this species.
				Specimens: nil,
			},
			Species{
				genome: genomeA, // The specimen does match this one.
				Specimens: []Specimen{
					Specimen{NeuralNet: NeatNeuralNet{Genome: genomeA}},
				},
			},
			Species{
				genome:    genomeA, // The specimen would match this species, but was already placed in a species.
				Specimens: nil,
			},
		},
	}
	c.Assert(population, DeepEquals, expectedPopulation)
}

func (s *PopulationSuite) Test_Population_AddSpecimen_NoMatchingSpecies(c *C) {
	var population generationPopulation
	var expectedPopulation generationPopulation

	// We have two genomes that compute to a speciation distance of 1.0 (with the constants we're using)
	var genomeA neatGenome = gnm([]gn{gn{1, 0.2}})
	var genomeB neatGenome = gnm([]gn{gn{0, 0.0}, gn{2, 0.4}})

	// Use a config that will be easy to study withs species.
	var config PopulationConfig = PopulationConfig{
		Speciation: SpeciationConfig{
			Threshold: 0.5, // One longer tail should be enough to indicate the genomes are not the same species.
			C1:        1.0, // Emphasies the longer tail of the younger genome.
			C2:        0.0,
			C3:        0.0,
		},
	}

	// Population is a new population with no species.
	population = newPopulation(config)
	population.species = append(population.species, Species{genome: genomeB}) // Won't match this species.

	// Add a specimen, will with match the second species.
	population.AddNeuralNet(NeatNeuralNet{Genome: genomeA}, 0.0, 0.0, nil)
	expectedPopulation = generationPopulation{
		config: config,
		species: []Species{
			Species{
				genome:    genomeB, // The specimen doesn't match this species.
				Specimens: nil,
			},
			Species{
				genome: genomeA, // This new species was added.
				Specimens: []Specimen{
					Specimen{NeuralNet: NeatNeuralNet{Genome: genomeA}},
				},
			},
		},
	}
	c.Assert(population, DeepEquals, expectedPopulation)

}
