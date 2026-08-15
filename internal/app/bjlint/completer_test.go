package bjlint

import (
	"github.com/tkmru/bjlint-cli/internal/pkg/blackjack"
	"reflect"
	"testing"
)

func suggestionNames(s blackjack.ActionSet) []string {
	x := Suggestions(s)
	out := make([]string, len(x))
	for i, v := range x {
		out[i] = v.Text
	}
	return out
}
func aset(xs ...blackjack.Action) blackjack.ActionSet {
	var s blackjack.ActionSet
	for _, x := range xs {
		s.Add(x)
	}
	return s
}
func TestSuggestions(t *testing.T) {
	tests := []struct {
		name string
		set  blackjack.ActionSet
		want []string
	}{{"two cards", aset(blackjack.ActionHit, blackjack.ActionStand, blackjack.ActionDouble), []string{"hit", "stand", "double", "quit"}}, {"pair", aset(blackjack.ActionHit, blackjack.ActionStand, blackjack.ActionDouble, blackjack.ActionSplit), []string{"hit", "stand", "double", "split", "quit"}}, {"three cards or no bankroll", aset(blackjack.ActionHit, blackjack.ActionStand), []string{"hit", "stand", "quit"}}}
	for _, tt := range tests {
		if got := suggestionNames(tt.set); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%s got %v", tt.name, got)
		}
	}
}
