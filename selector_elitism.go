package genetic

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// SelectorElitism just keeps the highest scoring members of each generation.
// It is a naive solution for selecting fittest members of a generation.
type SelectorElitism struct {
	KeepCount int // How many specimens should be kept in each generation?
}

// Select the N fittest members of a population.
func (s *SelectorElitism) Select(specimens []Specimen) (fittest []Specimen) {
	// The specimens are sorted from fittest to least fit.
	return specimens[:s.KeepCount]
}

// LoadSelectorElitismConfig loads the json filename as a new configuration.
func LoadSelectorElitismConfig(filename string) (SelectorElitism, error) {
	var err error
	var bytes []byte
	var selector SelectorElitism

	log.Printf("Loading elitism selector Config: '%s'\n", filename)

	// Load and parse from json.
	if bytes, err = ioutil.ReadFile(filename); err != nil {
		return SelectorElitism{}, err
	}
	if err = json.Unmarshal(bytes, &selector); err != nil {
		return SelectorElitism{}, err
	}

	return selector, error(nil)
}
