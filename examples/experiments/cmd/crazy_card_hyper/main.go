/*
	crazy_card_hyper is a genetic experiment trying to breed/mutate a neural net that can correctly guess a secret winning
	hand selected at the beginning of the experiment. It is configurable how big the hand is and how many decks together make up
	the card pool the neural net can pick from.

	This version of crazy card hand uses multiple outcomes (instead of a single score) and hypervolume indicators.

	The experiment is not particularly relevant, except architecturally. It shows how all the parts of hte go-genetic package can
	be pulled together to make an elaborate test of some underlying model.
*/
package main

import (
	"flag"
	"github.com/glemzurg/go-genetic"
	"log"
	"os"
)

func main() {
	var err error

	// Pass the configuration files as parameters:
	//
	//   bin/crazy_card_hyper -scorer=path/to/scorer.json -genetic=path/to/genetic.json -sorter=path/to/sorter.json -selector=path/to/selector.json
	//   bin/crazy_card_hyper -h
	//
	var scorerFilename *string = flag.String("scorer", "", "json file configuration for scoring part of experiment")
	var geneticFilename *string = flag.String("genetic", "", "json file configuration for the genetic part of experiment")
	var sorterFilename *string = flag.String("sorter", "", "json file configuration for the sorter part of experiment")
	var selectorFilename *string = flag.String("selector", "", "json file configuration for the selector part of experiment")
	flag.Parse()

	// Report the files we are using.
	var experimentName string = os.Args[0]
	log.Printf("Experiment: '%s'\n", experimentName)

	var scorer Scorer
	if scorer, err = LoadConfig(*scorerFilename); err != nil {
		log.Panic(err)
	}

	var geneticConfig genetic.Config
	if geneticConfig, err = genetic.LoadConfig(*geneticFilename); err != nil {
		log.Panic(err)
	}

	var sorter genetic.SorterHypervolumeIndicator
	if sorter, err = genetic.LoadSorterHypervolumeIndicatorConfig(*sorterFilename); err != nil {
		log.Panic(err)
	}

	var selector genetic.SelectorElitism
	if selector, err = genetic.LoadSelectorElitismConfig(*selectorFilename); err != nil {
		log.Panic(err)
	}

	// Run the experiment.
	genetic.RunExperiment(experimentName, geneticConfig, &sorter, &selector, &scorer)
	log.Println("Experiment Complete.")
	os.Exit(0)
}
