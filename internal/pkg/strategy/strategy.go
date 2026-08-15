package strategy

import "github.com/tkmru/bjlint-cli/internal/pkg/blackjack"

type Recommendation struct {
	Primary     blackjack.Action
	Fallback    blackjack.Action
	HasFallback bool
}

func (r Recommendation) Effective(available blackjack.ActionSet) blackjack.Action {
	if available.Has(r.Primary) {
		return r.Primary
	}
	if r.HasFallback {
		return r.Fallback
	}
	return r.Primary
}

// Recommend implements the 6-deck, dealer-hits-soft-17, DAS, no-surrender table.
func Recommend(hand blackjack.Hand, dealer blackjack.Card, rules blackjack.Rules, available blackjack.ActionSet) Recommendation {
	_ = rules // Rules is explicit so future tables can be selected without coupling callers.
	var r Recommendation
	if pair, ok := hand.PairValue(); ok {
		r = pairRecommendation(pair, dealer.Rank.Value())
	} else if total, soft := hand.Value(); soft {
		r = softRecommendation(total, dealer.Rank.Value())
	} else {
		r = hardRecommendation(total, dealer.Rank.Value())
	}
	if !available.Has(r.Primary) && r.HasFallback {
		r.Primary = r.Fallback
		r.HasFallback = false
	}
	return r
}

func doubleOr(fallback blackjack.Action) Recommendation {
	return Recommendation{blackjack.ActionDouble, fallback, true}
}
func splitOr(fallback blackjack.Action) Recommendation {
	return Recommendation{blackjack.ActionSplit, fallback, true}
}
func action(a blackjack.Action) Recommendation { return Recommendation{Primary: a} }

func Deviates(recommendation Recommendation, chosen blackjack.Action) bool {
	return recommendation.Primary != chosen
}
