//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// Suit represents the suit of a playing card (among Spade, Diamond, Club, Heart)
type Suit uint8

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

// Rank represents the rank of a playing car, from Ace to King
type Rank uint8

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

const (
	minRank = Ace
	maxRank = King
)

// Card represents a playing card with a Value and a Suite
type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

// New returns a deck of 52 Cards
func New(options ...func([]Card) []Card) []Card {
	var deck []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			deck = append(deck, Card{Suit: suit, Rank: rank})
		}
	}

	for _, opt := range options {
		deck = opt(deck)
	}

	return deck
}

// DefaultSort a deck of card in the standard order
func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}

// Less returns the less function for sorting the cards by default
func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) <= absRank(cards[j])
	}
}

// Sort allows user to specify custom sorting function
func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		sort.Slice(cards, less(cards))
		return cards
	}
}

// Suffle the deck of cards randomly
func Suffle(cards []Card) []Card {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return cards
}

// Jokers adds two jockers to the deck of cards
func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		ret := cards
		for i := 0; i < n; i++ {
			// Add a rank to differenciate the jokers eventualy
			ret = append(ret, Card{Rank: Rank(i), Suit: Joker})
		}
		return ret
	}
}

// Filter out the cards identified by the fltr function
func Filter(fltr func(Card) bool) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for _, card := range cards {
			if !fltr(card) {
				ret = append(ret, card)
			}
		}
		return ret
	}
}

// Deck creates a deck composed of n duplicated decks
func Deck(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		var ret []Card
		for i := 0; i < n; i++ {
			ret = append(ret, cards...)
		}
		return ret
	}
}
