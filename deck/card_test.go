package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Five, Suit: Spade})
	fmt.Println(Card{Rank: King, Suit: Diamond})
	fmt.Println(Card{Rank: Jack, Suit: Club})
	fmt.Println(Card{Suit: Joker})

	// Output:
	// Ace of Hearts
	// Five of Spades
	// King of Diamonds
	// Jack of Clubs
	// Joker
}

func TestNew(t *testing.T) {
	cards := New()
	// 13 ranks * 4 suits
	if len(cards) != 13*4 {
		t.Error("Wrong number of cards in a new deck.")
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	expected := Card{Rank: Ace, Suit: Spade}
	if cards[0] != expected {
		t.Error("Expected", expected, "as first card. Recieved: ", cards[0])
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	expected := Card{Rank: Ace, Suit: Spade}
	if cards[0] != expected {
		t.Error("Expected", expected, "as first card. Recieved: ", cards[0])
	}
}

func TestJokers(t *testing.T) {
	numOfJokers := 2
	cards := New(Jokers(numOfJokers))
	count := 0
	for _, card := range cards {
		if card.Suit == Joker {
			count++
		}
	}
	if count != numOfJokers {
		t.Error("Expected", numOfJokers, "jokers. Got", count)
	}
}

func TestFilter(t *testing.T) {
	filterAces := func(c Card) bool {
		return c.Rank == Ace
	}
	cards := New(Filter(filterAces))
	for _, card := range cards {
		if card.Rank == Ace {
			t.Errorf("Found %s. Supposed to be filtered", card)
		}
	}
}

func TestDeck(t *testing.T) {
	numOfDeck := 3
	cards := New(Deck(numOfDeck))
	if len(cards) != numOfDeck*52 {
		t.Errorf("Expected %d cards. Got %d", numOfDeck*52, len(cards))
	}
}
