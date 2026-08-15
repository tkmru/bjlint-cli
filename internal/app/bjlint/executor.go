package bjlint

import (
	"fmt"

	"github.com/tkmru/bjlint-cli/internal/pkg/strategy"
)

func (a *App) Execute(input string) {
	cmd, err := Parse(input)
	if err != nil {
		a.message = capitalize(err.Error())
		return
	}
	if cmd.Quit {
		a.quit = true
		return
	}
	if cmd.Action == nil {
		return
	}
	available := a.game.AvailableActions()
	if !available.Has(*cmd.Action) {
		a.message = fmt.Sprintf("Invalid action: %s is not available for this hand.", cmd.Action.String())
		return
	}
	h := a.game.Current()
	rec := strategy.Recommend(h.Hand, a.game.Dealer.Cards[0], a.game.Rules, available)
	a.message = ""
	if strategy.Deviates(rec, *cmd.Action) {
		a.message = fmt.Sprintf("⚠ %s deviates from Basic Strategy — recommended: %s", cmd.Action.String(), rec.Primary)
	}
	if err := a.game.Act(*cmd.Action); err != nil {
		a.message = fmt.Sprintf("Error: %v", err)
		return
	}
	if settlement := a.completeRound(); settlement != "" {
		if a.message != "" {
			a.message += "\n\n"
		}
		a.message += settlement
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	b := []byte(s)
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 32
	}
	return string(b)
}
