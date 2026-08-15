package bjlint

import (
	"bytes"
	"strings"
	"testing"

	"github.com/tkmru/bjlint-cli/internal/pkg/blackjack"
)

func appCard(r blackjack.Rank) blackjack.Card {
	return blackjack.Card{Rank: r, Suit: blackjack.Clubs}
}

func testApp(t *testing.T, cards []blackjack.Card) (*App, *bytes.Buffer) {
	t.Helper()
	g, err := blackjack.NewGameWithShoe(blackjack.DefaultRules(), blackjack.NewShoeFromCards(cards))
	if err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	return NewWithGame(g, &out), &out
}

func TestExecutorWarnsButExecutesDeviation(t *testing.T) {
	a, _ := testApp(t, []blackjack.Card{appCard(blackjack.Ten), appCard(blackjack.Six), appCard(blackjack.Six), appCard(blackjack.Ten), appCard(blackjack.Two)})
	a.Execute("hit")
	if !strings.Contains(a.message, "recommended: STAND") || len(a.game.Current().Hand.Cards) != 3 {
		t.Fatalf("warning or execution missing: %q", a.message)
	}
}

func TestExecutorRejectsUnavailableActionWithoutMutation(t *testing.T) {
	a, _ := testApp(t, []blackjack.Card{appCard(blackjack.Four), appCard(blackjack.Ten), appCard(blackjack.Four), appCard(blackjack.Six), appCard(blackjack.Two)})
	a.Execute("hit")
	before := len(a.game.Current().Hand.Cards)
	a.Execute("double")
	if !strings.Contains(a.message, "Invalid action: DOUBLE") || len(a.game.Current().Hand.Cards) != before {
		t.Fatalf("invalid action changed state or message: %q", a.message)
	}
}

func TestExecutorUsesEffectiveFallbackWithoutWarning(t *testing.T) {
	a, _ := testApp(t, []blackjack.Card{appCard(blackjack.Ace), appCard(blackjack.Six), appCard(blackjack.Two), appCard(blackjack.Ten), appCard(blackjack.Four), appCard(blackjack.Two)})
	a.Execute("hit") // A,2 -> A,2,4 (soft 17); double is now unavailable.
	a.Execute("hit")
	if strings.Contains(a.message, "WARNING") {
		t.Fatalf("fallback hit produced warning: %q", a.message)
	}
}

func TestExecutorUnknownCommandDoesNotChangeState(t *testing.T) {
	a, _ := testApp(t, []blackjack.Card{appCard(blackjack.Four), appCard(blackjack.Ten), appCard(blackjack.Four), appCard(blackjack.Six)})
	before := a.game.Shoe.Remaining()
	a.Execute("foo")
	if !strings.Contains(a.message, "Unknown command: foo") || a.game.Shoe.Remaining() != before {
		t.Fatal("unknown command changed state or message")
	}
}

func TestExecutorRedrawsInteractiveScreen(t *testing.T) {
	a, out := testApp(t, []blackjack.Card{appCard(blackjack.Four), appCard(blackjack.Ten), appCard(blackjack.Four), appCard(blackjack.Six)})
	a.screenMode = true
	a.Execute("foo")
	if out.Len() != 0 {
		t.Fatalf("executor wrote while go-prompt owned the screen: %q", out.String())
	}
	a.renderScreen(a.message)
	if !strings.HasPrefix(out.String(), "\x1b[2J\x1b[H") || !strings.Contains(out.String(), "DEALER") {
		t.Fatalf("interactive screen was not fully redrawn: %q", out.String())
	}
}

func TestExecutorSkipsImmediatelySettledNextRound(t *testing.T) {
	cards := []blackjack.Card{
		// Round 1: player 18 pushes dealer 18.
		appCard(blackjack.Ten), appCard(blackjack.Ten), appCard(blackjack.Eight), appCard(blackjack.Eight),
		// Round 2: dealer blackjack, settled without prompting.
		appCard(blackjack.Nine), appCard(blackjack.Ace), appCard(blackjack.Seven), appCard(blackjack.King),
	}
	a, _ := testApp(t, cards)
	a.Execute("stand")
	if a.game.Phase != blackjack.PhasePlayerTurn {
		t.Fatalf("prompt would have no actions in phase %v", a.game.Phase)
	}
	if !a.game.AvailableActions().Has(blackjack.ActionHit) || !strings.Contains(a.message, "blackjack") {
		t.Fatalf("automatic settlement was not retained: %q", a.message)
	}
}
