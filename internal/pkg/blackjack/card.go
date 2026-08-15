package blackjack

import "fmt"

type Rank uint8

const (
	Ace Rank = iota + 1
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

var AllRanks = []Rank{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}

func (r Rank) String() string {
	names := [...]string{"?", "A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
	if int(r) >= len(names) {
		return "?"
	}
	return names[r]
}

func (r Rank) Value() int {
	if r >= Ten {
		return 10
	}
	return int(r)
}

type Suit uint8

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

var AllSuits = []Suit{Clubs, Diamonds, Hearts, Spades}

func (s Suit) String() string {
	n := [...]string{"Clubs", "Diamonds", "Hearts", "Spades"}
	if int(s) >= len(n) {
		return "?"
	}
	return n[s]
}

type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string { return fmt.Sprintf("[%s]", c.Rank) }
