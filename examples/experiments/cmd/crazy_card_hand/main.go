/*
	crazy_card_hand is a genetic experiment trying to breed/mutate a neural net that can correctly reconstruct a secret winning
	hand selected at the beginning of the experiment. It is configurable how big the hand is and how many decks together make up
	the card pool the neural net can pick from.

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
	//   bin/crazy_card_hand -scorer=path/to/scorer.json -genetic=path/to/genetic.json -selector=path/to/selector.json
	//   bin/crazy_card_hand -h
	//
	var scorerFilename *string = flag.String("scorer", "", "json file configuration for scoring part of experiment")
	var geneticFilename *string = flag.String("genetic", "", "json file configuration for the genetic part of experiment")
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

	var selector genetic.SelectorElitism
	if selector, err = genetic.LoadSelectorElitismConfig(*selectorFilename); err != nil {
		log.Panic(err)
	}

	// We are maximizing the score for this experiment.
	var sorter genetic.Sorter = genetic.NewSorterSimpleMaximize() // Higher scores are fitter.

	// Run the experiment.
	genetic.RunExperiment(experimentName, geneticConfig, sorter, &selector, &scorer)
	log.Println("Experiment Complete.")
	os.Exit(0)
}
