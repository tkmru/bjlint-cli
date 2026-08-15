package blackjack

import "testing"

func tc(r Rank) Card { return Card{Rank: r, Suit: Clubs} }

func TestHandValue(t *testing.T) {
	tests := []struct {
		name       string
		cards      []Card
		total      int
		soft, bust bool
	}{
		{"soft 17", []Card{tc(Ace), tc(Six)}, 17, true, false},
		{"hard 17", []Card{tc(Ace), tc(Six), tc(Ten)}, 17, false, false},
		{"two aces", []Card{tc(Ace), tc(Ace)}, 12, true, false},
		{"multiple aces 21", []Card{tc(Ace), tc(Ace), tc(Nine)}, 21, true, false},
		{"faces", []Card{tc(King), tc(Queen)}, 20, false, false},
		{"bust", []Card{tc(Ten), tc(Six), tc(Eight)}, 24, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHand(tt.cards...)
			got, soft := h.Value()
			if got != tt.total || soft != tt.soft || h.IsBust() != tt.bust {
				t.Fatalf("got (%d,%v,%v)", got, soft, h.IsBust())
			}
		})
	}
}

func TestBlackjack(t *testing.T) {
	for _, r := range []Rank{Ten, Jack, Queen, King} {
		for _, cards := range [][]Card{{tc(Ace), tc(r)}, {tc(r), tc(Ace)}} {
			if !NewHand(cards...).IsBlackjack() {
				t.Errorf("expected blackjack: %v", cards)
			}
		}
	}
	bad := []Hand{{Cards: []Card{tc(Ace), tc(King)}, FromSplit: true}, NewHand(tc(Ace), tc(Five), tc(Five)), NewHand(tc(Seven), tc(Seven), tc(Seven))}
	for _, h := range bad {
		if h.IsBlackjack() {
			t.Errorf("unexpected blackjack: %v", h)
		}
	}
}

func TestCanSplitByValue(t *testing.T) {
	tests := []struct {
		a, b Rank
		want bool
	}{{Eight, Eight, true}, {Ace, Ace, true}, {King, Queen, true}, {Jack, Ten, true}, {Eight, Seven, false}}
	for _, tt := range tests {
		if got := NewHand(tc(tt.a), tc(tt.b)).CanSplit(); got != tt.want {
			t.Errorf("%s,%s=%v", tt.a, tt.b, got)
		}
	}
}
