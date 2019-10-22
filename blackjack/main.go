package main

import (
	"fmt"
	"gophercises/deck"
	"strconv"
	"strings"
)

// A Hand is a set of cards held by a player
type Hand []deck.Card

// A Stage represents the stage of the game
type Stage uint8

const (
	PlayerTurn Stage = iota
	DealerTurn
	Finished
)

// A GameState represents the stage of the game, the deck and the players hands
type GameState struct {
	Deck           []deck.Card
	Stage          Stage
	Player, Dealer Hand
	PlayerMoney    int
	PlayerBet      int
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	ret, cards := cards[0], cards[1:]
	return ret, cards
}

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i, c := range h {
		strs[i] = c.String()
	}
	return strings.Join(strs, ", ")
}

// DealerString returns the stringified dealer's hand
func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

// Score returns the value of a blackjack hand
func (h Hand) Score() (score int) {
	score = h.MinScore()
	if score > 11 {
		return score
	}

	for _, c := range h {
		if c.Rank == deck.Ace {
			// Adds 10 to current value of ace (1) to be 11
			return score + 10
		}
	}

	return score
}

// MinScore return the minimum value of a blackjack hand, considering aces as 1
func (h Hand) MinScore() (score int) {
	for _, c := range h {
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

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:        make([]deck.Card, len(gs.Deck)),
		Stage:       gs.Stage,
		Player:      make([]deck.Card, len(gs.Player)),
		Dealer:      make([]deck.Card, len(gs.Dealer)),
		PlayerMoney: gs.PlayerMoney,
		PlayerBet:   gs.PlayerBet,
	}
	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Player)
	copy(ret.Dealer, gs.Dealer)
	return ret
}

// CurrentPlayer returns the current players hand
func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.Stage {
	case PlayerTurn:
		return &gs.Player
	case DealerTurn:
		return &gs.Dealer
	default:
		panic("it isn't any player's turn")
	}
}

// Suffle shuffles the deck of the gamestate
func Suffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Deck(3), deck.Suffle)
	return ret
}

// Deal deals 2 cards to the player and the dealer, in the right order
func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make([]deck.Card, 0, 5)
	ret.Dealer = make([]deck.Card, 0, 5)

	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Deck = draw(ret.Deck)
		ret.Player = append(ret.Player, card)
		card, ret.Deck = draw(ret.Deck)
		ret.Dealer = append(ret.Dealer, card)
	}
	ret.Stage = PlayerTurn
	return ret
}

// Hit a new card for the current player
func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.CurrentPlayer()
	var card deck.Card
	card, ret.Deck = draw(ret.Deck)
	*hand = append(*hand, card)
	if hand.Score() > 21 {
		return Stand(ret)
	}
	return ret
}

// Stand the current players hand and updates the game stage
func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.Stage++
	return ret
}

// EndHand ends the current game and displays the outcome
func EndHand(gs GameState) GameState {
	ret := clone(gs)
	pScore, dScore := ret.Player.Score(), ret.Dealer.Score()
	fmt.Println("=== FINAL HANDS ===")
	fmt.Printf("Player: %s. Score: %d\n", ret.Player, pScore)
	fmt.Printf("Dealer: %s. Score: %d\n", ret.Dealer, dScore)

	switch {
	case pScore > 21:
		fmt.Println("You busted.")
		ret = lose(ret)
	case dScore > 21:
		fmt.Println("Dealer busted.")
		ret = win(ret)
	case pScore > dScore:
		fmt.Println("You won!")
		ret = win(ret)
	case pScore < dScore:
		fmt.Println("You lose")
		ret = lose(ret)
	case pScore == dScore:
		fmt.Println("Draw")
		ret = exAequo(ret)
	}
	fmt.Println()

	ret.Player = nil
	ret.Dealer = nil
	return ret
}

func win(gs GameState) GameState {
	ret := clone(gs)
	ret.PlayerMoney += 2 * ret.PlayerBet
	ret.PlayerBet = 0
	return ret
}

func lose(gs GameState) GameState {
	ret := clone(gs)
	ret.PlayerBet = 0
	return ret
}

func exAequo(gs GameState) GameState {
	ret := clone(gs)
	ret.PlayerMoney += ret.PlayerBet
	ret.PlayerBet = 0
	return ret
}

func main() {
	var gs GameState
	gs = Suffle(gs)
	numGames := 3
	minBet := 5
	gs.PlayerMoney = 100

	for i := 0; i < numGames && gs.PlayerMoney >= minBet; i++ {
		gs = Deal(gs)
		var input string

		for input == "" {
			fmt.Printf("\nHow much do you want to bet? %d-%d$\n", minBet, gs.PlayerMoney)
			fmt.Scanf("%s\n", &input)
			bet, err := strconv.Atoi(input)
			if err != nil {
				fmt.Println("Invalid bet amount")
				input = ""
			} else if bet > gs.PlayerMoney {
				fmt.Println("You don't have enough money")
				input = ""
			} else if bet < minBet {
				fmt.Println("You have to bet more money")
				input = ""
			} else {
				gs.PlayerBet = bet
				gs.PlayerMoney -= bet
			}
		}

		for gs.Stage == PlayerTurn {
			fmt.Println("----------------------------------------")
			fmt.Println("Player:", gs.Player)
			fmt.Println("Dealer:", gs.Dealer.DealerString())

			fmt.Println("\nWhat do you want to do? (h)it, (s)tand")
			fmt.Scanf("%s\n", &input)

			switch input {
			case "h":
				gs = Hit(gs)
			case "s":
				gs = Stand(gs)
			default:
				fmt.Println("Invalid action")

			}
		}

		for gs.Stage == DealerTurn {
			dScore := gs.Dealer.Score()
			dMinScore := gs.Dealer.MinScore()
			if dScore <= 16 || (dScore == 17 && dMinScore < 17) {
				gs = Hit(gs)
			} else {
				gs = Stand(gs)
			}
		}

		gs = EndHand(gs)
	}
	fmt.Printf("You now have %d$\n", gs.PlayerMoney)
}
