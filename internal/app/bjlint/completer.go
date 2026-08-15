package bjlint

import (
	"strings"

	prompt "github.com/c-bata/go-prompt"
	"github.com/tkmru/bjlint-cli/internal/pkg/blackjack"
)

func Suggestions(actions blackjack.ActionSet) []prompt.Suggest {
	all := []struct {
		a blackjack.Action
		s prompt.Suggest
	}{
		{blackjack.ActionHit, prompt.Suggest{Text: "hit", Description: "Hit another card"}},
		{blackjack.ActionStand, prompt.Suggest{Text: "stand", Description: "Stand with current hand"}},
		{blackjack.ActionDouble, prompt.Suggest{Text: "double", Description: "Double the bet and draw one card"}},
		{blackjack.ActionSplit, prompt.Suggest{Text: "split", Description: "Split the pair"}},
	}
	out := make([]prompt.Suggest, 0, 5)
	for _, x := range all {
		if actions.Has(x.a) {
			out = append(out, x.s)
		}
	}
	out = append(out, prompt.Suggest{Text: "quit", Description: "Quit bjlint"})
	return out
}

func (a *App) complete(d prompt.Document) []prompt.Suggest {
	word := strings.ToLower(d.GetWordBeforeCursor())
	return prompt.FilterHasPrefix(Suggestions(a.game.AvailableActions()), word, true)
}
