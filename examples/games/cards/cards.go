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

// FaceValue calcultes the numeric value from the face.
func FaceValue(face string) int {
	switch face {
	case "joker":
		return JOKER_VALUE
	case "ace":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	case "7":
		return 7
	case "8":
		return 8
	case "9":
		return 9
	case "10":
		return 10
	case "jack":
		return 11
	case "queen":
		return 12
	case "king":
		return 13
	}
	panic(fmt.Sprintf("Unknown card face: '%f'", face))
	return 0 // Should never happen.
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
