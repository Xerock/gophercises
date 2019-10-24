package blackjack

import (
	"errors"
	"gophercises/deck"
)

// A Stage represents the stage of the game
type Stage uint8

const (
	PlayerTurn Stage = iota
	DealerTurn
	Finished
)

// Options of a blackjack game
type Options struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
}

// New returns a new game
func New(opts Options) Game {
	g := Game{
		stage:    PlayerTurn,
		dealerAI: dealerAI{},
		balance:  0,
	}
	if opts.Decks == 0 {
		opts.Decks = 3
	}
	if opts.Hands == 0 {
		opts.Hands = 100
	}
	if opts.BlackjackPayout == 0. {
		opts.BlackjackPayout = 1.5
	}
	g.nDecks = opts.Decks
	g.nHands = opts.Hands
	g.blackjackPayout = opts.BlackjackPayout

	return g
}

// Game represent a blackjack game
type Game struct {
	nDecks          int
	nHands          int
	blackjackPayout float64

	stage Stage
	deck  []deck.Card

	player    []deck.Card
	playerBet int
	balance   int

	dealer   []deck.Card
	dealerAI AI
}

func (g *Game) currentHand() *[]deck.Card {
	switch g.stage {
	case PlayerTurn:
		return &g.player
	case DealerTurn:
		return &g.dealer
	default:
		panic("it isn't any player's turn")
	}
}

func bet(g *Game, ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)
	if bet < 100 {
		panic("Bet must be at least 100")
	}
	g.playerBet = bet
}

func deal(g *Game) {
	g.player = make([]deck.Card, 0, 5)
	g.dealer = make([]deck.Card, 0, 5)

	var card deck.Card
	for i := 0; i < 2; i++ {
		card, g.deck = draw(g.deck)
		g.player = append(g.player, card)
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
	}
	g.stage = PlayerTurn
}

// Play a game of blackjack
func (g *Game) Play(ai AI) int {
	g.deck = nil
	minCardsLeft := 52 * g.nDecks / 3
	for i := 0; i < g.nHands; i++ {
		shuffled := false
		if len(g.deck) < minCardsLeft {
			g.deck = deck.New(deck.Deck(g.nDecks), deck.Suffle)
			shuffled = true
		}

		bet(g, ai, shuffled)
		deal(g)

		for g.stage == PlayerTurn {
			hand := make([]deck.Card, len(g.player))
			copy(hand, g.player)
			move := ai.Play(hand, g.dealer[0])
			err := move(g)
			if err != nil {
				switch err {
				case errBust:
					MoveStand(g)
				default:
					panic(err)
				}
			}
		}
		for g.stage == DealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}
		endHand(g, ai)
	}
	return g.balance
}

var (
	errBust = errors.New("Hand exceded 21")
)

// Move defines a valid action in a blackjack game
type Move func(*Game) error

// MoveHit executes a hit action on the game
func MoveHit(g *Game) error {
	hand := g.currentHand()
	var card deck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		return errBust
	}
	return nil
}

// MoveStand executes a stand action on the game
func MoveStand(g *Game) error {
	g.stage++
	return nil
}

// MoveDouble executes the double action on the game
func MoveDouble(g *Game) error {
	if len(g.player) != 2 {
		return errors.New("Can only double on a 2 cards hand")
	}
	g.playerBet *= 2
	MoveHit(g)
	return MoveStand(g)
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	ret, cards := cards[0], cards[1:]
	return ret, cards
}

func endHand(g *Game, ai AI) {
	pScore, dScore := Score(g.player...), Score(g.dealer...)
	pBjack, dBjack := Blackjack(g.player...), Blackjack(g.dealer...)
	winning := g.playerBet
	switch {
	case pBjack && dBjack:
		winning = 0
	case dBjack:
		winning = -winning
	case pBjack:
		winning *= int(g.blackjackPayout)
	case pScore > 21:
		winning = -winning
	case dScore > 21:
		// win
	case pScore > dScore:
		// win
	case pScore < dScore:
		winning = -winning
	case pScore == dScore:
		winning = 0
	}
	g.balance += winning
	ai.Results(g.player, g.dealer)

	g.player = nil
	g.dealer = nil
}

// Blackjack returns true if a hand is a blackjack
func Blackjack(hand ...deck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
}

// Score returns the value of a blackjack hand
func Score(hand ...deck.Card) (score int) {
	score = minScore(hand...)
	if score > 11 {
		return score
	}

	for _, c := range hand {
		if c.Rank == deck.Ace {
			// Adds 10 to current value of ace (1) to be 11
			return score + 10
		}
	}

	return score
}

// Soft returns true if the current score is a soft score (ace counted as 11)
func Soft(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)
	return minScore != score
}

func minScore(cards ...deck.Card) (score int) {
	for _, c := range cards {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
