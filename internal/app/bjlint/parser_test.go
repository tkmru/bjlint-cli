package bjlint

import (
	"github.com/tkmru/bjlint-cli/internal/pkg/blackjack"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		in     string
		action blackjack.Action
		quit   bool
	}{{"h", blackjack.ActionHit, false}, {"H", blackjack.ActionHit, false}, {"hit", blackjack.ActionHit, false}, {"HIT", blackjack.ActionHit, false}, {" hit ", blackjack.ActionHit, false}, {"s", blackjack.ActionStand, false}, {"stand", blackjack.ActionStand, false}, {"d", blackjack.ActionDouble, false}, {"double", blackjack.ActionDouble, false}, {"p", blackjack.ActionSplit, false}, {"split", blackjack.ActionSplit, false}, {"q", 0, true}, {"quit", 0, true}, {"exit", 0, true}}
	for _, tt := range tests {
		got, err := Parse(tt.in)
		if err != nil {
			t.Fatal(err)
		}
		if got.Quit != tt.quit {
			t.Errorf("%q quit", tt.in)
		}
		if !tt.quit && (got.Action == nil || *got.Action != tt.action) {
			t.Errorf("%q action=%v", tt.in, got.Action)
		}
	}
	if _, err := Parse("foo"); err == nil {
		t.Fatal("unknown command accepted")
	}
}
