package strategy

import "github.com/tkmru/bjlint-cli/internal/pkg/blackjack"

func pairRecommendation(pair, dealer int) Recommendation {
	switch pair {
	case 1, 8:
		return splitOr(blackjack.ActionHit)
	case 2, 3:
		if dealer >= 2 && dealer <= 7 {
			return splitOr(blackjack.ActionHit)
		}
	case 4:
		if dealer == 5 || dealer == 6 {
			return splitOr(blackjack.ActionHit)
		}
	case 5:
		return hardRecommendation(10, dealer)
	case 6:
		if dealer >= 2 && dealer <= 6 {
			return splitOr(blackjack.ActionHit)
		}
	case 7:
		if dealer >= 2 && dealer <= 7 {
			return splitOr(blackjack.ActionHit)
		}
	case 9:
		if (dealer >= 2 && dealer <= 6) || dealer == 8 || dealer == 9 {
			return splitOr(blackjack.ActionStand)
		}
		return action(blackjack.ActionStand)
	case 10:
		return action(blackjack.ActionStand)
	}
	return hardRecommendation(pair*2, dealer)
}
