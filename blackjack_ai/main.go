package main

import (
	"fmt"
	"gophercises/blackjack_ai/blackjack"
	"gophercises/deck"
)

type betterAI struct {
	score int
	seen  int
	decks int
}

func (ai *betterAI) Play(hand []deck.Card, dealer deck.Card) blackjack.Move {
	score := blackjack.Score(hand...)
	dScore := blackjack.Score(dealer)
	switch {
	case (score == 10 || score == 11) && len(hand) == 2:
		return blackjack.MoveDouble
	case dScore >= 5 && dScore <= 6:
		return blackjack.MoveStand
	case blackjack.Score(hand...) <= 13:
		return blackjack.MoveHit
	default:
		return blackjack.MoveStand
	}
}

func (ai *betterAI) Bet(shuffled bool) int {
	minBet := 100
	if shuffled {
		ai.score = 0
		ai.seen = 0
	}
	trueScore := ai.score / ((52*ai.decks - ai.seen) / 52)
	switch {
	case trueScore > 14:
		return 1000 * minBet
	case trueScore > 8:
		return 50 * minBet
	default:
		return minBet
	}
}

func (ai *betterAI) Results(hand []deck.Card, dealer []deck.Card) {
	for _, card := range dealer {
		ai.count(card)
	}
	for _, card := range hand {
		ai.count(card)
	}
}

func (ai *betterAI) count(c deck.Card) {
	score := blackjack.Score(c)
	switch {
	case score >= 10:
		ai.score--
	case score <= 6:
		ai.score++
	}
	ai.seen++
}

func main() {
	opts := blackjack.Options{
		Decks:           3,
		Hands:           50000,
		BlackjackPayout: 1.5,
	}
	g := blackjack.New(opts)

	winings := g.Play(&betterAI{
		score: 0,
		seen:  0,
		decks: 3,
	})

	fmt.Println(winings)
}
