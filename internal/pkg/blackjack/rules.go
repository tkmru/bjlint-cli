package blackjack

type Rules struct {
	Decks                int
	DealerHitsSoft17     bool
	DoubleAfterSplit     bool
	MaxHands             int
	ResplitAces          bool
	HitSplitAces         bool
	BlackjackNumerator   Money
	BlackjackDenominator Money
	Bet                  Money
	StartingBankroll     Money
	Penetration          float64
}

func DefaultRules() Rules {
	return Rules{Decks: 6, DealerHitsSoft17: true, DoubleAfterSplit: true, MaxHands: 4,
		ResplitAces: false, HitSplitAces: false, BlackjackNumerator: 3, BlackjackDenominator: 2,
		Bet: 1000, StartingBankroll: 10000, Penetration: .75}
}

func DealerShouldHit(h Hand, rules Rules) bool {
	total, soft := h.Value()
	return total < 17 || (total == 17 && soft && rules.DealerHitsSoft17)
}
