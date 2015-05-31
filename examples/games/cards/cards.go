/*
	Package cards provides the details of a standard 54 card deck (13 of each suit and two jokers)
*/
package cards

import (
	"fmt"
	"strconv"
)

const (
	// The suits.
	SUIT_SPADES   = "spades"
	SUIT_DIAMONDS = "diamonds"
	SUIT_CLUBS    = "clubs"
	SUIT_HEARTS   = "hearts"
	SUIT_JOKERS   = "jokers"

	JOKER_VALUE = -1 // Since golang can zero-initialize structures, don't use 0 for jokers.
	_MAX_JOKERS = 2  // The limit to jockers in a deck.
)

// Card is the basic details of a playing card.
type Card struct {
	Suit  string
	Face  string
	Value int
}

// newCard is a single well-formed card.
func newCard(suit string, value int) Card {
	var face string
	switch value {
	case JOKER_VALUE:
		face = "joker"
	case 1:
		face = "ace"
	case 11:
		face = "jack"
	case 12:
		face = "queen"
	case 13:
		face = "king"
	default:
		if value >= 2 && value <= 10 {
			face = strconv.Itoa(value)
		} else {
			panic(fmt.Sprintf("Unknown card value: %d", value))
		}
	}
	return Card{Suit: suit, Face: face, Value: value}
}

// NewUnshuffledDecks gets a deck of cards (maybe multiple decks)
func NewUnshuffledDecks(deckCount int, jokersPerDeck int) []Card {
	var cards []Card
	for i := 0; i < deckCount; i++ {
		var singleDeckCards []Card = singleDeck(jokersPerDeck)
		cards = append(cards, singleDeckCards...)
	}
	return cards
}

// singleDeck is all the cards in a single deck.
func singleDeck(jokersToInclude int) []Card {
	var cards []Card
	for _, suit := range []string{SUIT_SPADES, SUIT_DIAMONDS, SUIT_CLUBS, SUIT_HEARTS} {
		for i := 1; i <= 13; i++ {
			cards = append(cards, newCard(suit, i))
		}
	}
	for i := 0; i < _MAX_JOKERS && i < jokersToInclude; i++ {
		cards = append(cards, newCard(SUIT_JOKERS, JOKER_VALUE))
	}
	return cards
}
