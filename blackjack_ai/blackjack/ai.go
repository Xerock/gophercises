package blackjack

import (
	"fmt"
	"gophercises/deck"
)

// AI interface defines the behaviour of a blackjack AI
type AI interface {
	Play(hand []deck.Card, dealer deck.Card) Move
	// Bet an amount of money, knowing if the deck has recently been suffled
	Bet(shuffled bool) int
	Results(hand []deck.Card, dealer []deck.Card)
}

type dealerAI struct{}

// Play returns a Move for the dealer
func (ai dealerAI) Play(hand []deck.Card, dealer deck.Card) Move {
	dScore := Score(hand...)
	if dScore <= 16 || dScore == 17 && Soft(hand...) {
		return MoveHit
	}
	return MoveStand
}

func (ai dealerAI) Bet(shuffled bool) int {
	// Do nothing
	return 1
}

// Results ...
func (ai dealerAI) Results(hand []deck.Card, dealer []deck.Card) {
	// Do nothing
}

type humanAI struct{}

// HumanAI returns a human player
func HumanAI() AI {
	return humanAI{}
}

// Play returns a Move for a Human player
func (ai humanAI) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer)

		var input string
		fmt.Println("\nWhat do you want to do? (h)it, (s)tand, (d)ouble")
		fmt.Scanf("%s\n", &input)

		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		case "d":
			return MoveDouble
		default:
			fmt.Println("Invalid action")

		}
	}
}

func (ai humanAI) Bet(shuffled bool) int {
	fmt.Println("----------------------------------------")
	if shuffled {
		fmt.Println("The deck was just shuffled")
	}
	fmt.Println("What whould you like to bet?")
	var bet int
	fmt.Scanf("%d\n", &bet)
	return bet
}

// Results ...
func (ai humanAI) Results(hand []deck.Card, dealer []deck.Card) {
	fmt.Println("=== FINAL HANDS ===")
	fmt.Printf("Player: %s.\n", hand)
	fmt.Printf("Dealer: %s.\n", dealer)
}

type basicAI struct{}

// BasicAI returns a basic ai
func BasicAI() AI {
	return basicAI{}
}

func (ai basicAI) Play(hand []deck.Card, dealer deck.Card) Move {
	score := Score(hand...)
	dScore := Score(dealer)
	switch {
	case (score == 10 || score == 11) && len(hand) == 2:
		return MoveDouble
	case dScore >= 5 && dScore <= 6:
		return MoveStand
	case Score(hand...) <= 13:
		return MoveHit
	default:
		return MoveStand
	}
}

func (ai basicAI) Bet(shuffled bool) int {
	//rand.Seed(time.Now().Unix())
	//minBet := 100
	//return rand.Intn(1+minBet*10) + 100
	return 100
}

func (ai basicAI) Results(hand []deck.Card, dealer []deck.Card) {
	// do nothing
}
