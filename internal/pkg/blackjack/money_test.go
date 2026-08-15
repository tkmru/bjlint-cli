package blackjack

import "testing"

func TestMoneyString(t *testing.T) {
	for _, tt := range []struct {
		m Money
		s string
	}{{500, "$5.00"}, {750, "$7.50"}, {10000, "$100.00"}, {-500, "-$5.00"}} {
		if got := tt.m.String(); got != tt.s {
			t.Errorf("got %s want %s", got, tt.s)
		}
	}
}
