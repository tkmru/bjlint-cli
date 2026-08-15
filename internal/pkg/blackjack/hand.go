package blackjack

type Hand struct {
	Cards     []Card
	FromSplit bool
}

func NewHand(cards ...Card) Hand { return Hand{Cards: append([]Card(nil), cards...)} }

func (h Hand) Value() (total int, soft bool) {
	aces := 0
	for _, c := range h.Cards {
		total += c.Rank.Value()
		if c.Rank == Ace {
			aces++
		}
	}
	if aces > 0 && total+10 <= 21 {
		total += 10
		soft = true
	}
	return
}

func (h Hand) IsBust() bool { n, _ := h.Value(); return n > 21 }
func (h Hand) IsBlackjack() bool {
	if h.FromSplit || len(h.Cards) != 2 {
		return false
	}
	n, _ := h.Value()
	return n == 21
}
func (h Hand) CanSplit() bool {
	return len(h.Cards) == 2 && h.Cards[0].Rank.Value() == h.Cards[1].Rank.Value()
}
func (h Hand) PairValue() (int, bool) {
	if !h.CanSplit() {
		return 0, false
	}
	return h.Cards[0].Rank.Value(), true
}
