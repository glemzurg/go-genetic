package main

import (
	"encoding/json"
	"github.com/glemzurg/go-genetic"
	"glemzurg/games/l5r"
	"hash/crc64"
	"io/ioutil"
	"log"
	"sort"
	"strings"
)

// Scorer is the game-specific experiment scoring for a running experiment.
type Scorer struct {

	// Public members will be recorded in the configuration for the experiment.
	DeckCount int // How many decks together create the card pool.
	JokersPerDeck int // How many jokers should be included in each deck.
	WinningHand       []string // The faces of the winning hand (suit ignored), order is important.
	MaxNoveltyFingerprints int       // The number of fingerprints to remember when tracking novely search outcomes.

	// Private members will not be recorded in the configuration for the experiment.
	handSize int // The size of the winning hand we are trying to match.
	winningHandFaceCount map[string]int // Have a distilled version of the winning hand face counts.
	availableCards    []cards.Card // Cards that can be used to build our hands.
	noveltyTally            genetic.NoveltySearchTally // Keep track of how many times we've see each outcome.
	generationCardDensities map[string]int             // The number of times we've seen a card in the current generation, keyed by card suite + face.

	// We want to hash strings as integers quickly.
	// Neural nets only take numeric input.
	// We'll use the CRC checksums as the integer values for the strings.
	crcTable *crc64.Table
}

// LoadConfig loads the json filename as a new configuration.
func LoadConfig(filename string, availableCards l5r.AvailableCards) (Scorer, error) {
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

	// The cards we can use.
	scorer.availableCards = availableCards


	// Put the cards in a format we can use well.
	scorer.cardLookup = map[string]l5r.DynastyCard{}
	for _, card := range scorer.AvailableCards.Dynasty {
		scorer.cardLookup[card.Id] = card
	}

	// Initialize the novelty tally.
	scorer.noveltyTally = genetic.NewNoveltySearchTally(scorer.MaxNoveltyFingerprints)

	// Prepare the scorer for doing CRC checksums.
	scorer.crcTable = crc64.MakeTable(crc64.ECMA)

	return scorer, error(nil)
}

// toFaceCounts converts a hand to a form we can use for fuzzy-scoring.
func toFaceCounts(faces []string) map[string]int {
	var ok bool
	var faceCounts map[string]int = map[string]int{}
	// Make sure all map entries exist.
	for _, face := range faces {
		faceCounts[face] = 0
	}
	// Count all the faces.
	for _, face := range faces {
		faceCounts[face]++
	}
	return faceCounts
}

// toFingerprint creates a string that we can use for novelty detection.
funct toFingerprint(faceCounts map[string]int) string {
	// Get the keys of the map.
	var uniqueFaces []string
	for face := range faceCounts {
		uniqueFaces = append(uniqueFaces, face)
	}
	sort.Strings(uniqueFaces)
	var fingerprint string = strings.Join(uniqueFaces, "-")
	return fingerprint
}

// CardAnalysis is how the neural net evaluates a single card
type CardAnalysis struct {
	Card   cards.Card // Which card is it?
	Priority float64 // How important is this card compared to others?
}

// ByPriority implements sort.Interface to sort ascending by GeneId.
// Example: sort.Sort(ByPriority(analyses))
type ByPriority []CardAnalysis

func (a ByPriority) Len() int           { return len(a) }
func (a ByPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool { return a[i].Priority > a[j].Priority }

// Score determines the score of a single specimen in a generation.
func (s *Scorer) Score(neuralNet genetic.NeatNeuralNet, population []genetic.NeatNeuralNet, neuralNetIndex int) (score float64, bonus float64, outcomes []float64) {

	// Evealuate each possible card with the neural net.
	var analyses []CardAnalysis
	for _, card := range s.availableCards {

		// The neural net only takes floats.
		var value float64 = float64(card.Value)

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
			"Suit": suit,
			"Face":  face,
			"Value": value,
		}
		var outputs map[string]float64 = neuralNet.Compute(inputs)

		// Capture the analysis for this card.
		analyses = append(analyses, CardAnalysis{
			Card:   card,
			Priority: outputs["Priority"],
		})
	}

	// The neural net has rated every card. Sort them, higher priority first.
	sort.Sort(ByPriority(analyses))

	// Now build the hand until we have the number of desired cards.
	var handCards []Card
	var hand []string
	for i := 0; i < s.handSize; i++ {
		// The hand is just the faces.
		hand = append(hand, analysis[i].Card.Face)
		// For other reporting we want the cards.
		handCards = append(handCards, analysis[i].Card)
	}

	// We want to capture a heatmap of what cards are being picked across all hands in a generation.
	for _, card := range handCards {
		// Note that this card has been seen in this generation.
		s.generationCardDensities[card.Suit + " " + card.Face]++
	}

	// The neural net has built the hand.
	// Now, how did it do?
	var score float64
	
	// The order of the cards matters in this hand.
	// 
	// For each face that the neural net finds, give it 10 points.
	// For each face that the neural net finds in the right spot, give it +40 points.
	//
	// If hte hand is a perfect match the score will be 50 * hand size.
	//
	// Additionally, we want to reward the neural net for trying something new (novelty search).
	//
	// Give a bonus of 50 divided by the times this particular combinations of cards has been seen before (sorting irrelevant)
	
	// First, give the score for the right hand composition (but maybe in the wrong order).
	var handFaceCounts map[string]int = toFaceCounts(hand)
	for face, count := range handFaceCounts {

		// Does this face exist in the winning hand?
		var ok bool
		if winningCount, ok = s.winningHandFaceCounts[face]; ok {
			
			// This face does exist in the winning hand.
			// We only get counted for each that is a match. Too many faces doesn't give more points.
			// Use whichever count is lower as the match count.
			var matchCount int = handFaceCount
			if winningCount < matchCount {
				matchCount = winningCount
			}
			
			// Give the score!
			score += 10.0 * float64(matchCount)
		}
	}
	
	// Now let's reward the neural net for putting the cards in the right places.
	for i, face := range hand {
		if face = winningHand[i] {
			score += 40.0
		}
	}
	
	// Now let's give a bonus for exploring new combinations.
	// Reward seeing a hand with cards we've never seen before.
	// In this case, the fingerprint of the hand we just looked at will be a sorted distillation of just the faces used.
	var fingerprint string = toFingerprint(handFaceCounts)
	var seen int = s.noveltyTally.Seen(fingerprint)
	var bonus float64 = 50.0 / float64(seen)

	return score, bonus, nil // No multi-outcomes.
}

func (s *Scorer) GenerationStart(generationNum uint64) {
	// Reset the captured card densities for this new generation.
	s.generationCardDensities = map[string]int{}
	// No card has been seen yet for this generation.
	for _, card := range s.availableCards {
		s.generationCardDensities[card.Suit + " " + card.Face] = 0
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
