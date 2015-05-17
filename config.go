package genetic

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config is the genetic-specific experiment configuration details (loaded from a json file).
type Config struct {
	Population   PopulationConfig   // How should each generation's population be managed.
	CppnInOut    CppnInOut          // What is the interface to the CPPNs in this experiment.
	EndCondition EndConditionConfig // What determines when the experiment should end. If nothing, must manually stop.
	Database     DatabaseConfig     // Database settings.
}

// EndConditionConfig describes how a genetic experiment should end. If blank, then the experiment must be manually stopped.
type EndConditionConfig struct {
	GenerationNum           uint64  // If not 0, the experiment ends after the specified generation no matter what.
	TargetScore             float64 // If not 0.0, the experiment ends when an experiment reaches this score.
	StagnantGenerationCount uint64  // If not 0, the experiment ends if the specified count of generations go by without fitness improving.
}

// PopulationConfig describes how the population should be managed.
type PopulationConfig struct {
	PopulationSize int              // How many specimens are in each generation of the experiment?
	Speciation     SpeciationConfig // What rules are used for identifying whether two specimens are the same species?
	Mutate         MutateConfig     // Rules for mating and mutating new members of the population.
}

// DatabaseConfig describes how the database should be interacted with.
type DatabaseConfig struct {
	RecordEveryNthGeneration uint64 // If 0, only record final generation. Otherwise, record every nth generation.
}

// SpeciationConfig describes how species are discovered to group specimens together by similarity.
type SpeciationConfig struct {
	Threshold float64 // Two genomes with a speciation distance below this number will be members of the same species.
	C1        float64 // A high configuration C1 gives more importance to excess genes (the tail of the longer genome).
	C2        float64 // A high configuration C2 gives more importance to disjoint genes (non-shared genes in either genome before the excess genes).
	C3        float64 // A high configuration C3 gives more importance to differences in shared genes.
}

// MutateConfig describes how new members of a population are created.
type MutateConfig struct {
	AvailableNodeFunctions   []string // A list of the activation functions we should use for creating new nodes. (e.g. ["bipolar_sigmoid", "inverse", "sine"])
	MaxAddConnectionAttempts int      // When adding a connection, we may create invalid ones. How many attempts until we just decide alter the weight of an existing connection.
	MateWeight               uint     // How likely is it that we'll mate during a mutation change. 6 is twice as likely to occur as 3.
	AddNodeWeight            uint     // How likely is it that we'll split an existing connection with a new node during a mutation change. 6 is twice as likely to occur as 3.
	AddConnectionWeight      uint     // How likely is it that we'll add a new connection during a mutation change. 6 is twice as likely to occur as 3.
	AlterConnectionWeight    uint     // How likely is it that we'll change the weight of an existing connection during a mutation change. 6 is twice as likely to occur as 3.
}

// LoadConfig loads the json filename as a new configuration.
func LoadConfig(filename string) (Config, error) {
	var err error
	var bytes []byte
	var config Config

	log.Printf("Loading genetic Config: '%s'\n", filename)

	// Load and parse from json.
	if bytes, err = ioutil.ReadFile(filename); err != nil {
		return Config{}, err
	}
	if err = json.Unmarshal(bytes, &config); err != nil {
		return Config{}, err
	}

	return config, error(nil)
}
