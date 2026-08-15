package strategy

import "github.com/tkmru/bjlint-cli/internal/pkg/blackjack"

func softRecommendation(total, dealer int) Recommendation {
	switch total {
	case 13, 14:
		if dealer == 5 || dealer == 6 {
			return doubleOr(blackjack.ActionHit)
		}
	case 15, 16:
		if dealer >= 4 && dealer <= 6 {
			return doubleOr(blackjack.ActionHit)
		}
	case 17:
		if dealer >= 3 && dealer <= 6 {
			return doubleOr(blackjack.ActionHit)
		}
	case 18:
		if dealer >= 2 && dealer <= 6 {
			return doubleOr(blackjack.ActionStand)
		}
		if dealer == 7 || dealer == 8 {
			return action(blackjack.ActionStand)
		}
	case 19:
		// H17 difference: soft 19 doubles against dealer 6.
		if dealer == 6 {
			return doubleOr(blackjack.ActionStand)
		}
		return action(blackjack.ActionStand)
	default:
		if total >= 20 {
			return action(blackjack.ActionStand)
		}
	}
	return action(blackjack.ActionHit)
}
