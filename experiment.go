package genetic

import (
	"bufio"
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"reflect"
	"time"
)

const (
	_DEFAULT_END_GENERATION_NUM = math.MaxUint64 // Run experiment "forever."
)

// geneticExperiment holds the state of the running experiment between generations as well as the
// beginning and ending results.
type geneticExperiment struct {
	experimentName string   // The name of this experiment.
	experimentId   int64    // The SQL id for this experiment after it is added to the database.
	config         Config   // The configuration for this experiment.
	scorer         Scorer   // The scorer of each specimen of a generation.
	sorter         Sorter   // The sorter that orders the specimens for selection.
	selector       Selector // The selector of which specimens should continue on to the next generation.
	db             *sql.DB  // The database connection.
}

// RunExperiment runs a genetic experiment until stopped manually or an end condition is met.
func RunExperiment(experimentName string, config Config, sorter Sorter, selector Selector, scorer Scorer) {

	// Create the experiment.
	var experiment geneticExperiment = geneticExperiment{
		experimentName: experimentName,
		config:         config,
		scorer:         scorer,
		sorter:         sorter,
		selector:       selector,
		db:             newDatabaseConnection(),
	}

	// Ensure the interface defined in the config is valid.
	experiment.config.NeuralNetInOut.validate()

	// Get the randomness rolling.
	rand.Seed(time.Now().UnixNano())

	// Record the start of the experiment.
	experiment.recordStart()

	// Start the current max geneid.
	setMaxGeneId(0)

	// Create an initial neural net that will seed the population, creating
	// a single specimen in a single species. In the first generation, this neural net will
	// be mutated into a full population through the normal mechanism to fill out a generation.
	var population generationPopulation = newPopulation(experiment.config.Population)
	var neuralNet NeatNeuralNet = newNeatNeuralNet(experiment.config.NeuralNetInOut)
	population.AddNeuralNet(neuralNet, 0.0, 0.0, nil) // The specimen has no scores.

	// Keep track if this experiment becomes stagnant (fitness never improving).
	// In this case, we have found some maxima and cannot move beyond it.
	var bestExperimentScore float64 = 0.0
	var stagnantGenerationCount uint64 = 0

	// What generation is the last of the experiment?
	var endConditionGenerationNum uint64 = experiment.config.EndCondition.GenerationNum
	if endConditionGenerationNum == 0 {
		endConditionGenerationNum = _DEFAULT_END_GENERATION_NUM
	}

	// Every generation capture the details of the best member of that generation.
	var best string

	// Keep track of why the experiment ends.
	var endReason string

	// Create a channel for listening for a manual stop command from the user.
	var manualStopChannel chan bool = make(chan bool)

	// Start listening for a manual stop.
	go listenForManualStop(manualStopChannel)

	// Run a generation of the experiment.
	var generationNum uint64
	for generationNum = 1; generationNum <= endConditionGenerationNum; generationNum++ {

		// Tell the scorer that a new generaiton has started.
		// It may want to prepare internal data structures.
		experiment.scorer.GenerationStart(generationNum)

		// Fill out the population to the correct size.
		// We either have the first generation's initial specimen or we have
		// the fittest specimens from the prior generation.
		population.FillOut()

		// Dump the neural nets from the population for examining, ready for scoring.
		var neuralNets []NeatNeuralNet = population.DumpSpecimensAsNeuralNets()

		// Score each neural net, one at a time. Bundle with its scores to make a specimen.
		for i, neuralNet := range neuralNets {

			// For scoring get these results:
			//
			//  score is the score for the neural net
			//  bonus is decided by meta-decisions (e.g. novelty search), 0.0 if nothing
			//  outcomes are for use with multi-outcome selectors (e.g. hyper-volume indicator), null otherwise
			var score float64
			var bonus float64
			var outcomes []float64

			// Score this neural net.
			score, bonus, outcomes = scorer.Score(neuralNet, neuralNets, i)

			// Re-add the specimen into the population.
			population.AddNeuralNet(neuralNet, score, bonus, outcomes)
		}

		// Modify the scores of the specimens by the size of their species.
		population.WeightSpecies()

		// Dump the specimens from the population, ready for selection.
		var specimens []Specimen = population.DumpSpecimens()

		// Sort the specimens. The specimens earlier in the slice are considered more fit.
		var bestScore float64
		var sorted []Specimen
		bestScore, best, sorted = experiment.sorter.Sort(specimens)

		// Select the fittest specimens.
		var fittestSpecimens []Specimen = experiment.selector.Select(sorted)

		// Add the the specimens back into the population.
		population.AddAllSpecimens(fittestSpecimens)

		// Did we improve over prior generations?
		if experiment.sorter.IsMaximize() {
			// Maximizing score.
			if bestScore > bestExperimentScore {
				bestExperimentScore = bestScore
				stagnantGenerationCount = 0
			} else {
				stagnantGenerationCount++
			}
		} else {
			// Minimizing score.
			if bestScore < bestExperimentScore {
				bestExperimentScore = bestScore
				stagnantGenerationCount = 0
			} else {
				stagnantGenerationCount++
			}
		}

		// Is this experiment over?

		// Have we reached the final generation?
		if generationNum >= endConditionGenerationNum {
			// End the experiment.
			endReason = fmt.Sprintf("reached generation: %d", generationNum)
			break
		}

		// Have we reached a target score?
		if experiment.sorter.IsMaximize() {
			// Maximizing score.
			if experiment.config.EndCondition.TargetScore > 0.0 && bestExperimentScore >= experiment.config.EndCondition.TargetScore {
				// End the experiment.
				endReason = fmt.Sprintf("target score %f reached: %f", experiment.config.EndCondition.TargetScore, bestExperimentScore)
				break
			}
		} else {
			// Minimizing score.
			// Since an uninitialized config will give a target score of 0.0, assume that is our target or anything else specified.
			if bestExperimentScore <= experiment.config.EndCondition.TargetScore {
				// End the experiment.
				endReason = fmt.Sprintf("target score %f reached: %f", experiment.config.EndCondition.TargetScore, bestExperimentScore)
				break
			}
		}

		// Is the experiment stuck and not improving?
		if experiment.config.EndCondition.StagnantGenerationCount > 0 && stagnantGenerationCount >= experiment.config.EndCondition.StagnantGenerationCount {
			// End the experiment.
			endReason = fmt.Sprintf("stagnant generation reached: %d", stagnantGenerationCount)
			break
		}

		// Is there a manual stop from the user?
		// Anything sent of the manual stop channel means it's time to stop.
		var isStop bool
		// Use select to create non-blocking channel receive.
		select {
		case isStop = <-manualStopChannel: // Only true will ever be sent over this channel.
		default: // Nothing to do, but creates a non-blocking receive.
		}
		// Manual stop?
		if isStop {
			// End the experiment.
			endReason = fmt.Sprintf("manual stop triggered")
			break
		}

		// Is this a generation we want to record?
		// We don't want to record every generation because of the time it takes to write.
		var isRecordGeneration bool = false
		if experiment.config.Database.RecordEveryNthGeneration > 0 {
			isRecordGeneration = ((generationNum % experiment.config.Database.RecordEveryNthGeneration) == 0)
		}

		// Record the generation of the experiment.
		if isRecordGeneration {
			experiment.recordGeneration(generationNum, bestExperimentScore, stagnantGenerationCount, best, population)
		}
	}

	// If we just ended the experiment we have yet to record this last generation.
	experiment.recordGeneration(generationNum, bestExperimentScore, stagnantGenerationCount, best, population)

	// Record the end of the experiment.
	experiment.recordEnd(generationNum, endReason, population)
}

// recordStart records details about the experiment before it runs.
func (e *geneticExperiment) recordStart() {
	var err error
	var bytes []byte

	// Get the config as json.
	if bytes, err = json.Marshal(e.config); err != nil {
		log.Panic(err)
	}
	var configJson string = string(bytes)

	// Get the scorer as json.
	// The experiment name indicates the "type" of the scorer.
	if bytes, err = json.Marshal(e.scorer); err != nil {
		log.Panic(err)
	}
	var scorerJson string = string(bytes)

	// Get the sorter as json.
	if bytes, err = json.Marshal(e.sorter); err != nil {
		log.Panic(err)
	}
	var sorterJson string = string(bytes)

	// Also save the type of scorer we're using.
	var sorterType string = reflect.ValueOf(e.sorter).Elem().Type().String()

	// Get the selector as json.
	if bytes, err = json.Marshal(e.selector); err != nil {
		log.Panic(err)
	}
	var selectorJson string = string(bytes)

	// Also save the type of selector we're using.
	var selectorType string = reflect.ValueOf(e.selector).Elem().Type().String()

	// Write the experiment to the database.
	var result sql.Result
	if result, err = e.db.Exec(
		`INSERT INTO genetic.experiment
         SET experiment=?,
             datetime=NOW(),
             config=?,
             scorer=?,
             sorter_type=?,
             sorter=?,
             selector_type=?,
             selector=?`,
		e.experimentName,
		configJson,
		scorerJson,
		sorterType,
		sorterJson,
		selectorType,
		selectorJson); err != nil {

		log.Panic(err)
	}

	// Depending on data in the database zero to two rows may be effected.
	var rowsAffected int64
	if rowsAffected, err = result.RowsAffected(); err != nil {
		log.Panic(err)
	}
	if rowsAffected != 1 {
		log.Panicf("Inserting experiment expected 1 row affected but was: %d", rowsAffected)
	}

	// What is the insert id we just created?
	var experimentId int64
	if experimentId, err = result.LastInsertId(); err != nil {
		log.Panic(err)
	}
	log.Printf("Experiment %d starting.\n", experimentId)

	// Remember the experiment for future SQL.
	e.experimentId = experimentId
}

// recordGeneration records details of a single generation of the experiment.
func (e *geneticExperiment) recordGeneration(generationNum uint64, bestExperimentScore float64, stagnantGenerationCount uint64, best string, population generationPopulation) {
	var err error

	// Get generation details from the scorer.
	var scorerBytes []byte = e.scorer.GenerationDetails()

	// Write the core experiment record to the database.
	var result sql.Result
	if result, err = e.db.Exec(
		`INSERT INTO genetic.experiment_generation
         SET experimentid=?,
             generation_num=?,
             datetime=NOW(),
             best_experiment_score=?,
             stagnant_generations=?,
             best=?,
             details=?`,
		e.experimentId,
		generationNum,
		bestExperimentScore,
		stagnantGenerationCount,
		best,
		scorerBytes); err != nil {

		log.Panic(err)
	}

	// Depending on data in the database zero to two rows may be effected.
	var rowsAffected int64
	if rowsAffected, err = result.RowsAffected(); err != nil {
		log.Panic(err)
	}
	if rowsAffected != 1 {
		log.Panicf("Inserting experiment generation expected 1 row affected but was: %d", rowsAffected)
	}

	// Add each species to the database with an overview of it.
	for _, species := range population.species {

		// We need to create a unique hash for the species.
		// Make an md5 of the genome.

		// Get the gnome as json.
		var bytes []byte
		if bytes, err = json.Marshal(species.genome); err != nil {
			log.Panic(err)
		}
		var genomeJson string = string(bytes)
		var genomeMd5 string = md5Of(genomeJson)

		// Species overview.
		var specimenCount int = len(species.Specimens)
		var specimenBest string
		var specimenBestScore float64
		specimenBestScore, specimenBest, _ = e.sorter.Sort(species.Specimens)

		// Write the core experiment record to the database.
		var result sql.Result
		if result, err = e.db.Exec(
			`INSERT INTO genetic.experiment_generation_species
         SET experimentid=?,
             generation_num=?,
             species_fingerprint=?,
             specimens=?,
             best_score=?,
             best=?`,
			e.experimentId,
			generationNum,
			genomeMd5,
			specimenCount,
			specimenBestScore,
			specimenBest); err != nil {

			log.Panic(err)
		}

		// Depending on data in the database zero to two rows may be effected.
		var rowsAffected int64
		if rowsAffected, err = result.RowsAffected(); err != nil {
			log.Panic(err)
		}
		if rowsAffected != 1 {
			log.Panicf("Inserting experiment generation species expected 1 row affected but was: %d", rowsAffected)
		}
	}
}

// recordEnd records details about the experiment after it ends.
func (e *geneticExperiment) recordEnd(generationNum uint64, endReason string, population generationPopulation) {
	var err error

	// Get the population as json.
	var bytes []byte
	if bytes, err = json.Marshal(population.species); err != nil {
		log.Panic(err)
	}
	var populationJson string = string(bytes)

	// Write the core experiment record to the database.
	var result sql.Result
	if result, err = e.db.Exec(
		`INSERT INTO genetic.experiment_end
         SET experimentid=?,
             end_reason=?,
             datetime=NOW(),
             generation_num=?,
             results=?`,
		e.experimentId,
		endReason,
		generationNum,
		populationJson); err != nil {

		log.Panic(err)
	}

	// Depending on data in the database zero to two rows may be effected.
	var rowsAffected int64
	if rowsAffected, err = result.RowsAffected(); err != nil {
		log.Panic(err)
	}
	if rowsAffected != 1 {
		log.Panicf("Inserting experiment end expected 1 row affected but was: %d", rowsAffected)
	}

	log.Printf("Experiment %d ended: %s\n", e.experimentId, endReason)
}

// listenForManualStop informs the experiment when a manual stop is issued by the user. Meant to be started
// as a goroutine.
func listenForManualStop(manualStopChannel chan bool) {
	var err error

	log.Println("To manually stop the experiment, press 'Q' then RETURN.")

	// Create a reader watching the standard input.
	var reader *bufio.Reader = bufio.NewReader(os.Stdin)

	// Loop until we hear something on stdin.
	for {

		var char rune

		// Anything to read?
		if char, _, err = reader.ReadRune(); err != nil {
			log.Panic(err)
		}

		// Is this the key to stop?
		if char == 'q' || char == 'Q' {
			manualStopChannel <- true
			close(manualStopChannel)
			break
		}
	}
}

// md5Of creates an md5 (as string) for an input string.
func md5Of(value string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(value)))
}
