package blackjack

import (
	"errors"
	"fmt"
	"math/rand"
)

type Phase uint8

const (
	PhasePlayerTurn Phase = iota
	PhaseDealerTurn
	PhaseRoundComplete
	PhaseGameOver
)

type Outcome uint8

const (
	OutcomeLose Outcome = iota
	OutcomePush
	OutcomeWin
	OutcomeBlackjack
)

func (o Outcome) String() string {
	switch o {
	case OutcomeLose:
		return "LOSE"
	case OutcomePush:
		return "PUSH"
	case OutcomeWin:
		return "WIN"
	case OutcomeBlackjack:
		return "BLACKJACK"
	default:
		return "UNKNOWN"
	}
}

type PlayerHand struct {
	Hand      Hand
	Wager     Money
	Stood     bool
	Doubled   bool
	SplitAces bool
}
type HandResult struct {
	Hand        int
	Outcome     Outcome
	Profit      Money
	Wager       Money
	PlayerTotal int
	DealerTotal int
}

type Game struct {
	Rules       Rules
	Bankroll    Money
	Shoe        *Shoe
	PlayerHands []*PlayerHand
	Dealer      Hand
	CurrentHand int
	Phase       Phase
	Results     []HandResult
	rng         *rand.Rand
	reshuffles  int
}

func NewGame(rules Rules, rng *rand.Rand) (*Game, error) {
	if rules.Bet <= 0 || rules.StartingBankroll < 0 {
		return nil, errors.New("invalid money configuration")
	}
	g := &Game{Rules: rules, Bankroll: rules.StartingBankroll, rng: rng}
	g.Shoe = NewShoe(rules.Decks, rng)
	if err := g.StartRound(); err != nil {
		return nil, err
	}
	return g, nil
}

// NewGameWithShoe is intended for deterministic games and tests.
func NewGameWithShoe(rules Rules, shoe *Shoe) (*Game, error) {
	g := &Game{Rules: rules, Bankroll: rules.StartingBankroll, Shoe: shoe}
	if err := g.StartRound(); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Game) draw() (Card, error) {
	c, ok := g.Shoe.Draw()
	if !ok {
		return Card{}, errors.New("shoe exhausted during round")
	}
	return c, nil
}

func (g *Game) StartRound() error {
	if g.Bankroll < g.Rules.Bet {
		g.Phase = PhaseGameOver
		return nil
	}
	// The cut card is checked only between rounds; a live round is never reshuffled.
	if g.Shoe == nil || g.Shoe.NeedsShuffle() {
		g.Shoe = NewShoe(g.Rules.Decks, g.rng)
		g.reshuffles++
	}
	g.Bankroll -= g.Rules.Bet
	g.PlayerHands = []*PlayerHand{{Wager: g.Rules.Bet}}
	g.Dealer = Hand{}
	g.Results = nil
	g.CurrentHand = 0
	g.Phase = PhasePlayerTurn
	for i := 0; i < 2; i++ {
		p, err := g.draw()
		if err != nil {
			return err
		}
		g.PlayerHands[0].Hand.Cards = append(g.PlayerHands[0].Hand.Cards, p)
		d, err := g.draw()
		if err != nil {
			return err
		}
		g.Dealer.Cards = append(g.Dealer.Cards, d)
	}
	if g.PlayerHands[0].Hand.IsBlackjack() || g.Dealer.IsBlackjack() {
		g.settle()
	}
	return nil
}

func (g *Game) Current() *PlayerHand {
	if g.Phase != PhasePlayerTurn || g.CurrentHand >= len(g.PlayerHands) {
		return nil
	}
	return g.PlayerHands[g.CurrentHand]
}

func (g *Game) AvailableActions() ActionSet {
	var set ActionSet
	h := g.Current()
	if h == nil || h.Stood || h.Hand.IsBust() {
		return set
	}
	set.Add(ActionHit)
	set.Add(ActionStand)
	if len(h.Hand.Cards) == 2 && g.Bankroll >= h.Wager && (!h.Hand.FromSplit || g.Rules.DoubleAfterSplit) && !h.SplitAces {
		set.Add(ActionDouble)
	}
	if h.Hand.CanSplit() && len(g.PlayerHands) < g.Rules.MaxHands && g.Bankroll >= h.Wager && !(h.Hand.Cards[0].Rank == Ace && h.Hand.FromSplit && !g.Rules.ResplitAces) && !h.SplitAces {
		set.Add(ActionSplit)
	}
	return set
}

func (g *Game) Act(action Action) error {
	if !g.AvailableActions().Has(action) {
		return fmt.Errorf("%s is not available for this hand", action)
	}
	h := g.Current()
	switch action {
	case ActionHit:
		c, err := g.draw()
		if err != nil {
			return err
		}
		h.Hand.Cards = append(h.Hand.Cards, c)
		if h.Hand.IsBust() {
			h.Stood = true
			return g.advance()
		}
	case ActionStand:
		h.Stood = true
		return g.advance()
	case ActionDouble:
		g.Bankroll -= h.Wager
		h.Wager *= 2
		h.Doubled = true
		c, err := g.draw()
		if err != nil {
			return err
		}
		h.Hand.Cards = append(h.Hand.Cards, c)
		h.Stood = true
		return g.advance()
	case ActionSplit:
		return g.split()
	}
	return nil
}

func (g *Game) split() error {
	h := g.Current()
	g.Bankroll -= h.Wager
	left, right := h.Hand.Cards[0], h.Hand.Cards[1]
	isAces := left.Rank == Ace
	h.Hand = Hand{Cards: []Card{left}, FromSplit: true}
	h.Stood = false
	h.SplitAces = isAces
	next := &PlayerHand{Hand: Hand{Cards: []Card{right}, FromSplit: true}, Wager: h.Wager, SplitAces: isAces}
	g.PlayerHands = append(g.PlayerHands, nil)
	copy(g.PlayerHands[g.CurrentHand+2:], g.PlayerHands[g.CurrentHand+1:])
	g.PlayerHands[g.CurrentHand+1] = next
	for _, ph := range []*PlayerHand{h, next} {
		c, err := g.draw()
		if err != nil {
			return err
		}
		ph.Hand.Cards = append(ph.Hand.Cards, c)
		if isAces {
			ph.Stood = true
		}
	}
	if isAces {
		return g.advance()
	}
	return nil
}

func (g *Game) advance() error {
	for g.CurrentHand+1 < len(g.PlayerHands) {
		g.CurrentHand++
		if !g.PlayerHands[g.CurrentHand].Stood {
			return nil
		}
	}
	allBust := true
	for _, h := range g.PlayerHands {
		if !h.Hand.IsBust() {
			allBust = false
			break
		}
	}
	if !allBust {
		g.Phase = PhaseDealerTurn
		for DealerShouldHit(g.Dealer, g.Rules) {
			c, err := g.draw()
			if err != nil {
				return err
			}
			g.Dealer.Cards = append(g.Dealer.Cards, c)
		}
	}
	g.settle()
	return nil
}

func (g *Game) settle() {
	g.Results = nil
	dealerTotal, _ := g.Dealer.Value()
	dealerBJ := g.Dealer.IsBlackjack()
	dealerBust := g.Dealer.IsBust()
	for i, h := range g.PlayerHands {
		pt, _ := h.Hand.Value()
		outcome := OutcomeLose
		payout := Money(0)
		switch {
		case h.Hand.IsBust():
			outcome = OutcomeLose
		case h.Hand.IsBlackjack() && dealerBJ:
			outcome = OutcomePush
			payout = h.Wager
		case h.Hand.IsBlackjack():
			outcome = OutcomeBlackjack
			payout = h.Wager + (h.Wager * g.Rules.BlackjackNumerator / g.Rules.BlackjackDenominator)
		case dealerBJ:
			outcome = OutcomeLose
		case dealerBust || pt > dealerTotal:
			outcome = OutcomeWin
			payout = 2 * h.Wager
		case pt == dealerTotal:
			outcome = OutcomePush
			payout = h.Wager
		default:
			outcome = OutcomeLose
		}
		g.Bankroll += payout
		g.Results = append(g.Results, HandResult{Hand: i + 1, Outcome: outcome, Profit: payout - h.Wager, Wager: h.Wager, PlayerTotal: pt, DealerTotal: dealerTotal})
	}
	g.Phase = PhaseRoundComplete
}

func (g *Game) ReshuffleCount() int { return g.reshuffles }

// DisplayBankroll includes stakes committed to an unfinished round. Settlement
// still operates only on Bankroll, avoiding a mix of stake and net accounting.
func (g *Game) DisplayBankroll() Money {
	total := g.Bankroll
	if g.Phase == PhasePlayerTurn || g.Phase == PhaseDealerTurn {
		for _, hand := range g.PlayerHands {
			total += hand.Wager
		}
	}
	return total
}
