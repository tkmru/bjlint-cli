package bjlint

import (
	"bytes"
	"strings"
	"testing"

	"github.com/tkmru/bjlint-cli/internal/pkg/blackjack"
)

func TestCardTokenIncludesRankAndSuit(t *testing.T) {
	a := &App{}
	for _, tt := range []struct {
		card blackjack.Card
		want string
	}{
		{blackjack.Card{Rank: blackjack.Ace, Suit: blackjack.Spades}, "A ♠"},
		{blackjack.Card{Rank: blackjack.Ten, Suit: blackjack.Hearts}, "10♥"},
		{blackjack.Card{Rank: blackjack.King, Suit: blackjack.Diamonds}, "K ♦"},
	} {
		if got := a.cardToken(tt.card); !strings.Contains(got, tt.want) {
			t.Errorf("card token %q does not contain %q", got, tt.want)
		}
	}
}

func TestActionLabelsOnlyShowAvailableActions(t *testing.T) {
	a := &App{}
	set := aset(blackjack.ActionHit, blackjack.ActionStand)
	got := a.actionLabels(set)
	if !strings.Contains(got, "[H]") || !strings.Contains(got, "[S]") || strings.Contains(got, "[D]") || strings.Contains(got, "[P]") {
		t.Fatalf("unexpected labels: %q", got)
	}
}

func TestActionPromptExplainsInputState(t *testing.T) {
	if actionPrompt != "YOUR MOVE › " {
		t.Fatalf("unexpected action prompt: %q", actionPrompt)
	}
}

func TestNoticeBorderHasConsistentWidth(t *testing.T) {
	var out bytes.Buffer
	a := &App{out: &out}
	a.renderNotice("⚠ STAND deviates from Basic Strategy — recommended: DOUBLE\n\nRound complete\nPlayer 14 • Dealer 21 (blackjack) • LOSE -$5.00")
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if width := displayWidth(line); width != panelWidth+2 {
			t.Errorf("notice line width=%d, want 57: %q", width, line)
		}
	}
}

func TestWelcomeBorderHasConsistentWidth(t *testing.T) {
	var out bytes.Buffer
	a := &App{out: &out}
	a.renderWelcome()
	for _, line := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if width := displayWidth(line); width != panelWidth+2 {
			t.Errorf("welcome line width=%d, want %d: %q", width, panelWidth+2, line)
		}
	}
}

func TestDisplayWidthTreatsTerminalSymbolsAsNarrow(t *testing.T) {
	if got := displayWidth("─ •"); got != 3 {
		t.Fatalf("displayWidth(\"─ •\")=%d, want 3", got)
	}
}
