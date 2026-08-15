package strategy

import (
	"github.com/tkmru/bjlint-cli/internal/pkg/blackjack"
	"testing"
)

func c(r blackjack.Rank) blackjack.Card { return blackjack.Card{Rank: r, Suit: blackjack.Clubs} }
func actions(xs ...blackjack.Action) blackjack.ActionSet {
	var s blackjack.ActionSet
	for _, x := range xs {
		s.Add(x)
	}
	return s
}

var all = actions(blackjack.ActionHit, blackjack.ActionStand, blackjack.ActionDouble, blackjack.ActionSplit)

func TestHardStrategy(t *testing.T) {
	tests := []struct {
		cards  []blackjack.Card
		dealer blackjack.Rank
		want   blackjack.Action
	}{
		{[]blackjack.Card{c(blackjack.Four), c(blackjack.Four)}, blackjack.Ten, blackjack.ActionHit},
		{[]blackjack.Card{c(blackjack.Five), c(blackjack.Six)}, blackjack.Six, blackjack.ActionDouble},
		{[]blackjack.Card{c(blackjack.Ten), c(blackjack.Two)}, blackjack.Four, blackjack.ActionStand},
		{[]blackjack.Card{c(blackjack.Ten), c(blackjack.Six)}, blackjack.Six, blackjack.ActionStand},
		{[]blackjack.Card{c(blackjack.Ten), c(blackjack.Six)}, blackjack.Ten, blackjack.ActionHit},
		{[]blackjack.Card{c(blackjack.Five), c(blackjack.Four)}, blackjack.Three, blackjack.ActionDouble},
		{[]blackjack.Card{c(blackjack.Ten), c(blackjack.Seven)}, blackjack.Ace, blackjack.ActionStand},
	}
	for _, tt := range tests {
		got := Recommend(blackjack.NewHand(tt.cards...), c(tt.dealer), blackjack.DefaultRules(), all).Primary
		if got != tt.want {
			t.Errorf("%v vs %s got %s want %s", tt.cards, tt.dealer, got, tt.want)
		}
	}
}

func TestSoftStrategy(t *testing.T) {
	tests := []struct {
		other, dealer blackjack.Rank
		want          blackjack.Action
	}{{blackjack.Six, blackjack.Six, blackjack.ActionDouble}, {blackjack.Seven, blackjack.Two, blackjack.ActionDouble}, {blackjack.Seven, blackjack.Nine, blackjack.ActionHit}, {blackjack.Eight, blackjack.Six, blackjack.ActionDouble}, {blackjack.Eight, blackjack.Ace, blackjack.ActionStand}}
	for _, tt := range tests {
		got := Recommend(blackjack.NewHand(c(blackjack.Ace), c(tt.other)), c(tt.dealer), blackjack.DefaultRules(), all).Primary
		if got != tt.want {
			t.Errorf("A,%s vs %s got %s", tt.other, tt.dealer, got)
		}
	}
}

func TestPairStrategy(t *testing.T) {
	tests := []struct {
		a, b, dealer blackjack.Rank
		want         blackjack.Action
	}{{blackjack.Eight, blackjack.Eight, blackjack.Ten, blackjack.ActionSplit}, {blackjack.Ace, blackjack.Ace, blackjack.Ten, blackjack.ActionSplit}, {blackjack.Ten, blackjack.Ten, blackjack.Six, blackjack.ActionStand}, {blackjack.Five, blackjack.Five, blackjack.Six, blackjack.ActionDouble}, {blackjack.King, blackjack.Queen, blackjack.Six, blackjack.ActionStand}, {blackjack.Jack, blackjack.Ten, blackjack.Six, blackjack.ActionStand}}
	for _, tt := range tests {
		got := Recommend(blackjack.NewHand(c(tt.a), c(tt.b)), c(tt.dealer), blackjack.DefaultRules(), all).Primary
		if got != tt.want {
			t.Errorf("%s,%s vs %s got %s", tt.a, tt.b, tt.dealer, got)
		}
	}
}

func TestFallbackAndLint(t *testing.T) {
	h := blackjack.NewHand(c(blackjack.Ace), c(blackjack.Six))
	withDouble := Recommend(h, c(blackjack.Six), blackjack.DefaultRules(), actions(blackjack.ActionHit, blackjack.ActionStand, blackjack.ActionDouble))
	if withDouble.Primary != blackjack.ActionDouble {
		t.Fatal(withDouble)
	}
	without := Recommend(h, c(blackjack.Six), blackjack.DefaultRules(), actions(blackjack.ActionHit, blackjack.ActionStand))
	if without.Primary != blackjack.ActionHit {
		t.Fatal(without)
	}
	if Deviates(without, blackjack.ActionHit) {
		t.Fatal("effective fallback incorrectly warned")
	}
	if !Deviates(action(blackjack.ActionStand), blackjack.ActionHit) {
		t.Fatal("deviation missed")
	}
	if Deviates(action(blackjack.ActionStand), blackjack.ActionStand) {
		t.Fatal("correct action warned")
	}
}
