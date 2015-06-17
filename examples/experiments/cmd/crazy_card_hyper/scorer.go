package main

import (
	"encoding/json"
	"github.com/glemzurg/go-genetic"
	"github.com/glemzurg/go-genetic/examples/games/cards"
	"hash/crc64"
	"io/ioutil"
	"log"
	"math"
	"sort"
)

// Scorer is the game-specific experiment scoring for a running experiment.
type Scorer struct {

	// Public members will be recorded in the configuration for the experiment.
	DeckCount     int      // How many decks together create the card pool.
	JokersPerDeck int      // How many jokers should be included in each deck.
	WinningHand   []string // The faces of the winning hand (suit ignored), order is important.

	// Private members will not be recorded in the configuration for the experiment.
	handSize                int            // How big is the secret hand?
	winningHandValues       []int          // The card value for each face of the winning hand.
	availableCards          []cards.Card   // Cards that can be used to build our hands.
	generationCardDensities map[string]int // The number of times we've seen a card in the current generation, keyed by card suite + face.

	// We want to hash strings as integers quickly.
	// Neural nets only take numeric input.
	// We'll use the CRC checksums as the integer values for the strings.
	crcTable *crc64.Table
}

// LoadConfig loads the json filename as a new configuration.
func LoadConfig(filename string) (Scorer, error) {
	var err error
	var bytes []byte
	var scorer Scorer

	log.Printf("Loading scorer Config: '%s'\n", filename)

	// Load and parse from json.
	if bytes, err = ioutil.ReadFile(filename); err != nil {
		return Scorer{}, err
	}

	if err = json.Unmarshal(bytes, &scorer); err != nil {
		return Scorer{}, err
	}

	// We need to know the size of the hand.
	scorer.handSize = len(scorer.WinningHand)

	// Get the values for each face in the winning hand.
	for _, face := range scorer.WinningHand {
		scorer.winningHandValues = append(scorer.winningHandValues, cards.FaceValue(face))
	}

	// The cards the neural net can use to make its hand.
	scorer.availableCards = cards.NewUnshuffledDecks(scorer.DeckCount, scorer.JokersPerDeck)

	// Prepare the scorer for doing CRC checksums.
	scorer.crcTable = crc64.MakeTable(crc64.ECMA)

	return scorer, error(nil)
}

// CardAnalysis is how the neural net evaluates a single card
type CardAnalysis struct {
	Card     cards.Card // Which card is it?
	Priority float64    // How important is this card compared to others?
}

// ByPriority implements sort.Interface to sort *descending* by Priority.
// Example: sort.Sort(ByPriority(analyses))
type ByPriority []CardAnalysis

func (a ByPriority) Len() int           { return len(a) }
func (a ByPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool { return a[i].Priority > a[j].Priority }

// Score determines the score of a single specimen in a generation.
func (s *Scorer) Score(neuralNet genetic.NeatNeuralNet, population []genetic.NeatNeuralNet, neuralNetIndex int) (score float64, bonus float64, outcomes []float64) {

	// Evealuate each possible card with the neural net.
	var analyses []CardAnalysis
	for i, card := range s.availableCards {

		// The neural net only takes floats.
		var value float64 = float64(card.Value)

		// When there are multiple decks, the neural net needs to distinguish between the
		// same card in each deck. Just included the index of the card in the pool of availbable cards
		// as the unique id for the card.
		//
		// These ids are fragile and only suitable in the experiment that is running. If a winning neural
		// net is asked to choose from an available card pool that is any other composition or order it
		// will fail in some way. We're imparting to the experiment intimate knowledge of this specific
		// card pool.
		var cardId float64 = float64(i)

		// Suits are strings, but need to be floats none-the-less.
		// Hash them to an arbitrary number, the only consistency being that the same
		// suit will always hash the the same number.
		var suitHash uint64 = crc64.Checksum([]byte(card.Suit), s.crcTable)
		var suit float64 = float64(suitHash)

		// Technically, the card.Value already captures all the information encompassed in the card.Face, but
		// for sake of example, it will also be passed into the neural net.

		// Faces are strings, but need to be floats none-the-less.
		// Hash them to an arbitrary number, the only consistency being that the same
		// face will always hash the the same number.
		var faceHash uint64 = crc64.Checksum([]byte(card.Suit), s.crcTable)
		var face float64 = float64(faceHash)

		// Run the neural net on these inputs.
		var inputs map[string]float64 = map[string]float64{
			"CardId": cardId,
			"Suit":   suit,
			"Face":   face,
			"Value":  value,
		}
		var outputs map[string]float64 = neuralNet.Compute(inputs)

		// Capture the analysis for this card.
		analyses = append(analyses, CardAnalysis{
			Card:     card,
			Priority: outputs["Priority"],
		})
	}

	// The neural net has rated every card. Sort them, higher priority first.
	sort.Sort(ByPriority(analyses))

	// Now build the hand until we have the number of desired cards.
	var handCards []cards.Card
	for i := 0; i < s.handSize; i++ {
		// For other reporting we want the cards.
		handCards = append(handCards, analyses[i].Card)
	}

	// We want to capture a heatmap of what cards are being picked across all hands in a single generation.
	for _, card := range handCards {
		// Note that this card has been seen in this generation.
		s.generationCardDensities[card.Suit+" "+card.Face]++
	}

	// The neural net has built the hand.
	// Now, how did it do?

	// We're going to use hypervolume indicators to analyze multiple outcomes (instead of a single score).

	// There is one outcome value for each possible card. The goal is to reach an outcome of 0.0 which means we've matched
	// the face of the card. For each value up or down from the card, add 1.0 to the score, moving it farther away from the desired score.
	// The hypercube dimensions have a reference value of 20.0 which is greater than they could possibly be (which is correct).

	// The max hypercube volume will be (20.0 - 0.0)^4 = 160000.0 so that's the hypercube volume means a perfect match of cards.

	outcomes = nil
	for i, winningValue := range s.winningHandValues {
		outcomes = append(outcomes, math.Abs(float64(winningValue)-float64(handCards[i].Value)))
	}

	// Return just the multiple outcomes and let the hypervolume code evaluate them.
	return 0.0, 0.0, outcomes
}

func (s *Scorer) GenerationStart(generationNum uint64) {
	// Reset the captured card densities for this new generation.
	s.generationCardDensities = map[string]int{}
	// No card has been seen yet for this generation.
	for _, card := range s.availableCards {
		s.generationCardDensities[card.Suit+" "+card.Face] = 0
	}
}

func (s *Scorer) GenerationDetails() (bytes []byte) {
	var err error
	// When recording details of a generation, given enough information that we could do a heat map of which cards were most often picked.
	if bytes, err = json.Marshal(s.generationCardDensities); err != nil {
		log.Println(err)
	}
	return bytes
}
