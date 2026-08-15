package blackjack

import "testing"

func TestDefaultBetIsTenDollars(t *testing.T) {
	if got := DefaultRules().Bet; got != 1000 {
		t.Fatalf("default bet=%s, want $10.00", got)
	}
}

func TestDealerH17(t *testing.T) {
	tests := []struct {
		cards []Card
		hit   bool
	}{{[]Card{tc(Ace), tc(Six)}, true}, {[]Card{tc(Ace), tc(Six), tc(Ten)}, false}, {[]Card{tc(Ten), tc(Seven)}, false}, {[]Card{tc(Nine), tc(Seven)}, true}, {[]Card{tc(Ace), tc(Seven)}, false}}
	for _, tt := range tests {
		if got := DealerShouldHit(NewHand(tt.cards...), DefaultRules()); got != tt.hit {
			t.Errorf("%v hit=%v", tt.cards, got)
		}
	}
}
