package bjlint

import (
	"fmt"
	"strings"

	"github.com/tkmru/bjlint-cli/internal/pkg/blackjack"
)

type Command struct {
	Action *blackjack.Action
	Quit   bool
}

func Parse(input string) (Command, error) {
	s := strings.ToLower(strings.TrimSpace(input))
	var a blackjack.Action
	switch s {
	case "":
		return Command{}, nil
	case "h", "hit":
		a = blackjack.ActionHit
	case "s", "stand":
		a = blackjack.ActionStand
	case "d", "double":
		a = blackjack.ActionDouble
	case "p", "split":
		a = blackjack.ActionSplit
	case "q", "quit", "exit":
		return Command{Quit: true}, nil
	default:
		return Command{}, fmt.Errorf("unknown command: %s", strings.TrimSpace(input))
	}
	return Command{Action: &a}, nil
}
