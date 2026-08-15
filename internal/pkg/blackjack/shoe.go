package blackjack

import "math/rand"

type Shoe struct {
	cards []Card
	next  int
	decks int
	cut   int
}

func NewShoe(decks int, rng *rand.Rand) *Shoe {
	if decks <= 0 {
		decks = 1
	}
	cards := make([]Card, 0, decks*52)
	for i := 0; i < decks; i++ {
		for _, s := range AllSuits {
			for _, r := range AllRanks {
				cards = append(cards, Card{r, s})
			}
		}
	}
	if rng != nil {
		rng.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	}
	return &Shoe{cards: cards, decks: decks, cut: int(float64(len(cards)) * .75)}
}

// NewShoeFromCards creates a deterministic shoe whose cards are drawn in order.
func NewShoeFromCards(cards []Card) *Shoe {
	return &Shoe{cards: append([]Card(nil), cards...), decks: 1, cut: int(float64(len(cards)) * .75)}
}
func (s *Shoe) Draw() (Card, bool) {
	if s.next >= len(s.cards) {
		return Card{}, false
	}
	c := s.cards[s.next]
	s.next++
	return c, true
}
func (s *Shoe) Remaining() int     { return len(s.cards) - s.next }
func (s *Shoe) Size() int          { return len(s.cards) }
func (s *Shoe) Used() int          { return s.next }
func (s *Shoe) NeedsShuffle() bool { return s.next >= s.cut }
func (s *Shoe) Cards() []Card      { return append([]Card(nil), s.cards...) }
