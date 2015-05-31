package main

import (
	"flag"
	"github.com/glemzurg/go-genetic"
	"glemzurg/games/l5r"
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

	var selector genetic.TruncateSelector
	if selector, err = genetic.LoadTruncateSelectorConfig(*selectorFilename); err != nil {
		log.Panic(err)
	}

	// Run the experiment.
	genetic.RunExperiment(experimentName, geneticConfig, &selector, &scorer)
	log.Println("Experiment Complete.")
	os.Exit(0)
}
