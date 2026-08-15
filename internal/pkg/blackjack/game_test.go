package blackjack

import "testing"

func gameWith(cards []Card, bankroll Money) (*Game, error) {
	r := DefaultRules()
	r.StartingBankroll = bankroll
	return NewGameWithShoe(r, NewShoeFromCards(cards))
}

func TestBankrollSettlement(t *testing.T) {
	tests := []struct {
		name   string
		cards  []Card
		action Action
		want   Money
	}{
		{"win", []Card{tc(Ten), tc(Nine), tc(Ten), tc(Eight)}, ActionStand, 11000},
		{"loss", []Card{tc(Ten), tc(Ten), tc(Six), tc(Nine)}, ActionStand, 9000},
		{"push", []Card{tc(Ten), tc(Ten), tc(Eight), tc(Eight)}, ActionStand, 10000},
		{"blackjack", []Card{tc(Ace), tc(Nine), tc(King), tc(Eight)}, ActionStand, 11500},
		{"double win", []Card{tc(Five), tc(Nine), tc(Six), tc(Eight), tc(King)}, ActionDouble, 12000},
		{"double loss", []Card{tc(Five), tc(Ten), tc(Six), tc(Nine), tc(Five)}, ActionDouble, 8000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := gameWith(tt.cards, 10000)
			if err != nil {
				t.Fatal(err)
			}
			if g.Phase == PhasePlayerTurn {
				if err = g.Act(tt.action); err != nil {
					t.Fatal(err)
				}
			}
			if g.Bankroll != tt.want {
				t.Fatalf("bankroll=%s want %s", g.Bankroll, tt.want)
			}
		})
	}
}

func TestDoubleRules(t *testing.T) {
	g, err := gameWith([]Card{tc(Five), tc(Ten), tc(Six), tc(Seven), tc(Ten)}, 10000)
	if err != nil {
		t.Fatal(err)
	}
	if !g.AvailableActions().Has(ActionDouble) {
		t.Fatal("two-card double unavailable")
	}
	if err = g.Act(ActionHit); err != nil {
		t.Fatal(err)
	}
	if g.AvailableActions().Has(ActionDouble) {
		t.Fatal("three-card double available")
	}
	g, err = gameWith([]Card{tc(Five), tc(Ten), tc(Six), tc(Seven)}, 500)
	if err != nil {
		t.Fatal(err)
	}
	if g.AvailableActions().Has(ActionDouble) {
		t.Fatal("double available without bankroll")
	}
}

func TestDoubleDrawsOnceAndStands(t *testing.T) {
	g, err := gameWith([]Card{tc(Five), tc(Ten), tc(Six), tc(Seven), tc(Ten)}, 10000)
	if err != nil {
		t.Fatal(err)
	}
	before := g.Shoe.Remaining()
	if err = g.Act(ActionDouble); err != nil {
		t.Fatal(err)
	}
	if len(g.PlayerHands[0].Hand.Cards) != 3 || g.Shoe.Remaining() != before-1 || g.PlayerHands[0].Wager != 2000 || !g.PlayerHands[0].Stood {
		t.Fatalf("bad double state: %+v", g.PlayerHands[0])
	}
}

func TestSplitAccountingAndIndependentHands(t *testing.T) {
	g, err := gameWith([]Card{tc(Eight), tc(Ten), tc(Eight), tc(Seven), tc(Ten), tc(Nine)}, 10000)
	if err != nil {
		t.Fatal(err)
	}
	if err = g.Act(ActionSplit); err != nil {
		t.Fatal(err)
	}
	if len(g.PlayerHands) != 2 || g.Bankroll != 8000 {
		t.Fatal("split stake or hand count")
	}
	if !g.AvailableActions().Has(ActionDouble) {
		t.Fatal("DAS unavailable")
	}
	if err = g.Act(ActionStand); err != nil {
		t.Fatal(err)
	}
	if g.CurrentHand != 1 {
		t.Fatal("did not advance")
	}
	if err = g.Act(ActionStand); err != nil {
		t.Fatal(err)
	}
	if g.Bankroll != 11000 || len(g.Results) != 2 {
		t.Fatalf("bankroll=%s results=%v", g.Bankroll, g.Results)
	}
}

func TestSplitAces(t *testing.T) {
	g, err := gameWith([]Card{tc(Ace), tc(Nine), tc(Ace), tc(Eight), tc(King), tc(Nine)}, 10000)
	if err != nil {
		t.Fatal(err)
	}
	if err = g.Act(ActionSplit); err != nil {
		t.Fatal(err)
	}
	if g.Phase != PhaseRoundComplete {
		t.Fatal("split aces did not auto-stand")
	}
	for _, h := range g.PlayerHands {
		if len(h.Hand.Cards) != 2 || !h.SplitAces || h.Hand.IsBlackjack() {
			t.Fatalf("bad split ace hand: %+v", h)
		}
	}
	if g.Bankroll != 12000 {
		t.Fatalf("split-ace 21 received blackjack payout: %s", g.Bankroll)
	}
}

func TestMaximumHandsAndInsufficientBankroll(t *testing.T) {
	cards := []Card{tc(Eight), tc(Ten), tc(Eight), tc(Seven), tc(Eight), tc(Eight), tc(Eight), tc(Eight), tc(Two), tc(Three)}
	g, err := gameWith(cards, 20000)
	if err != nil {
		t.Fatal(err)
	}
	for len(g.PlayerHands) < 4 {
		if err = g.Act(ActionSplit); err != nil {
			t.Fatal(err)
		}
	}
	if g.AvailableActions().Has(ActionSplit) {
		t.Fatal("split allowed beyond four hands")
	}
	g, err = gameWith([]Card{tc(Eight), tc(Ten), tc(Eight), tc(Seven)}, 500)
	if err != nil {
		t.Fatal(err)
	}
	if g.AvailableActions().Has(ActionSplit) {
		t.Fatal("split available without additional stake")
	}
}

func TestAllPlayerHandsBustSkipsDealerPlay(t *testing.T) {
	g, err := gameWith([]Card{tc(Ten), tc(Five), tc(Six), tc(Six), tc(King), tc(Ten)}, 10000)
	if err != nil {
		t.Fatal(err)
	}
	before := g.Shoe.Remaining()
	if err = g.Act(ActionHit); err != nil {
		t.Fatal(err)
	}
	if len(g.Dealer.Cards) != 2 || g.Shoe.Remaining() != before-1 {
		t.Fatal("dealer drew after every player hand had busted")
	}
}

func TestDealerBlackjackSettlement(t *testing.T) {
	g, err := gameWith([]Card{tc(Ace), tc(Ace), tc(King), tc(King)}, 10000)
	if err != nil {
		t.Fatal(err)
	}
	if g.Phase != PhaseRoundComplete || g.Results[0].Outcome != OutcomePush || g.Bankroll != 10000 {
		t.Fatalf("both blackjacks: phase=%v result=%v bankroll=%s", g.Phase, g.Results, g.Bankroll)
	}
	g, err = gameWith([]Card{tc(Ten), tc(Ace), tc(Nine), tc(King)}, 10000)
	if err != nil {
		t.Fatal(err)
	}
	if g.Results[0].Outcome != OutcomeLose || g.Bankroll != 9000 {
		t.Fatalf("dealer blackjack not settled as loss: %v %s", g.Results, g.Bankroll)
	}
}
