package genetic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
)

// SorterHypervolumeIndicator sorts multi-outcome specimens first by how must they improve on the rest of the population and
// secondly on how good their multi-outcomes are individually. Hypervolume indicators require essentially comparing each member
// of the population against all other members of the population (although there is some small optimication).
//
// The hypervolume indicator sort ignores the score and bonus and only operatons on the outcomes from the scorer.
type SorterHypervolumeIndicator struct {
	ReferencePoint []float64 // For each outcome, what is the base value the outcome is compared to.
	Maximize       []bool    // For each outcome, true means a higher value is more fit, false means a lower value is more fit.
	Weights        []float64 // For each outcome, what relative value does it have? A higher value means it will have more influence on the sort.
}

// LoadSorterHypervolumeIndicatorConfig loads the json filename as a new configuration.
func LoadSorterHypervolumeIndicatorConfig(filename string) (SorterHypervolumeIndicator, error) {
	var err error
	var bytes []byte
	var sorter SorterHypervolumeIndicator

	log.Printf("Loading hypervolume indicator Config: '%s'\n", filename)

	// Load and parse from json.
	if bytes, err = ioutil.ReadFile(filename); err != nil {
		return SorterHypervolumeIndicator{}, err
	}
	if err = json.Unmarshal(bytes, &sorter); err != nil {
		return SorterHypervolumeIndicator{}, err
	}
	sorter.validOrPanic()
	return sorter, error(nil)
}

// validOrPanic panics if we're not ready for use.
func (s *SorterHypervolumeIndicator) validOrPanic() {
	if len(s.ReferencePoint) == 0 {
		log.Panic("ReferencePoint must have values")
	}
	if len(s.ReferencePoint) != len(s.Maximize) || len(s.Maximize) != len(s.Weights) {
		log.Panicf("ReferencePoint (%v), Maximize (%v), and Weights (%v) must all have the same number of values", s.ReferencePoint, s.Maximize, s.Weights)
	}
	for _, weight := range s.Weights {
		if weight == 0.0 {
			log.Panicf("Weights with a zero value are not allowed: %v", s.Weights)
		}
	}
}

// Sort the specimens descending, first by how much they uniquely improve on the population, next on how good the specimen's multi-outcomes are.
func (s *SorterHypervolumeIndicator) Sort(specimens []Specimen) (bestScore float64, best string, sorted []Specimen) {

	// Create hypercubes of the specimens.
	var hypercubes []*specimenHypercube
	for _, specimen := range specimens {
		var hypercube specimenHypercube = newSpecimenHypercube(specimen, s.ReferencePoint, s.Maximize, s.Weights)
		hypercubes = append(hypercubes, &hypercube)
	}

	// Calculate the hypervolume indicators.
	calculateHypervolumeIndicators(hypercubes)

	// Give each specimen a selection score.
	var bestSpecimen *specimenHypercube
	for i := 0; i < len(hypercubes); i++ {

		// Score based on hypervolume indicator (will always be positive).
		var selectionScore float64 = hypercubes[i].indicator / float64(hypercubes[i].specimen.SpeciesMemberCount)
		hypercubes[i].specimen.setSelectionScore(selectionScore)

		// Is this the best specimen we've found?
		if bestSpecimen == nil {
			bestSpecimen = hypercubes[i]
		} else {
			if hypercubes[i].indicator > bestScore {
				bestSpecimen = hypercubes[i]
			}
		}
		bestScore = bestSpecimen.indicator
	}

	// What is the text summary of the best specimen.
	best = fmt.Sprintf("indicator: %f, volume: %f, speciesmembercount: %d, outcomes: %v", bestSpecimen.indicator, bestSpecimen.volume, bestSpecimen.specimen.SpeciesMemberCount, bestSpecimen.specimen.Outcomes)

	// Sort the hypercubes descending by selection score, then indicator, then volume.
	sort.Sort(byHypervolumeIndicatorDescending(hypercubes))

	// Dump the specimens into the original pointer list, preserving their order.
	sorted = nil
	for _, hypercube := range hypercubes {
		sorted = append(sorted, hypercube.specimen)
	}

	// The best information of the population (may not be the specimen at the head of the list).
	return bestScore, best, sorted
}

// IsMaximize returns true. Hypervolume indicator sort makes normalized hypercubes that increase in volume when fitter.
func (s *SorterHypervolumeIndicator) IsMaximize() bool { return true }

// byHypervolumeIndicatorDescending implements sort.Interface to sort descending by selection score, then indicator, then volume.
// Example: sort.Sort(byHypervolumeIndicatorDescending(hypercubes))
type byHypervolumeIndicatorDescending []*specimenHypercube

func (a byHypervolumeIndicatorDescending) Len() int      { return len(a) }
func (a byHypervolumeIndicatorDescending) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byHypervolumeIndicatorDescending) Less(i, j int) bool {
	if a[i].specimen.SelectionScore == a[j].specimen.SelectionScore {
		if a[i].indicator == a[j].indicator {
			return a[i].volume > a[j].volume // Third by volume.
		} else {
			return a[i].indicator > a[j].indicator // Second by indicator.
		}
	}
	return a[i].specimen.SelectionScore > a[j].specimen.SelectionScore // First by selection score.
}
