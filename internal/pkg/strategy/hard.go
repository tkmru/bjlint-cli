package strategy

import "github.com/tkmru/bjlint-cli/internal/pkg/blackjack"

func hardRecommendation(total, dealer int) Recommendation {
	switch {
	case total <= 8:
		return action(blackjack.ActionHit)
	case total == 9:
		if dealer >= 3 && dealer <= 6 {
			return doubleOr(blackjack.ActionHit)
		}
	case total == 10:
		if dealer >= 2 && dealer <= 9 {
			return doubleOr(blackjack.ActionHit)
		}
	case total == 11:
		// H17 difference: double 11 against an ace (represented as 1).
		return doubleOr(blackjack.ActionHit)
	case total == 12:
		if dealer >= 4 && dealer <= 6 {
			return action(blackjack.ActionStand)
		}
	case total >= 13 && total <= 16:
		if dealer >= 2 && dealer <= 6 {
			return action(blackjack.ActionStand)
		}
	case total >= 17:
		return action(blackjack.ActionStand)
	}
	return action(blackjack.ActionHit)
}
