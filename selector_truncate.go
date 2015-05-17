package genetic

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// TruncateSelector just keeps the highest scoring members of each generation.
// It is a naive solution for selecting fittest members of a generation.
type TruncateSelector struct {
	KeepCount int // How many specimens should be kept in each generation?
}

// Select the N fittest members of a population.
func (s *TruncateSelector) Select(specimens []Specimen) (fittest []Specimen) {
	// The specimens are sorted from fittest to least fit.
	return specimens[:s.KeepCount]
}

// LoadTruncateSelectorConfig loads the json filename as a new configuration.
func LoadTruncateSelectorConfig(filename string) (TruncateSelector, error) {
	var err error
	var bytes []byte
	var selector TruncateSelector

	log.Printf("Loading truncate selector Config: '%s'\n", filename)

	// Load and parse from json.
	if bytes, err = ioutil.ReadFile(filename); err != nil {
		return TruncateSelector{}, err
	}
	if err = json.Unmarshal(bytes, &selector); err != nil {
		return TruncateSelector{}, err
	}

	return selector, error(nil)
}
