package bjlint

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/tkmru/bjlint-cli/internal/pkg/blackjack"
	"github.com/tkmru/bjlint-cli/internal/pkg/strategy"
)

const (
	panelWidth = 55

	ansiReset  = "\x1b[0m"
	ansiBold   = "\x1b[1m"
	ansiDim    = "\x1b[2m"
	ansiRed    = "\x1b[31m"
	ansiGreen  = "\x1b[32m"
	ansiYellow = "\x1b[33m"
	ansiBlue   = "\x1b[34m"
	ansiCyan   = "\x1b[36m"
)

var terminalWidth = &runewidth.Condition{EastAsianWidth: false}

func displayWidth(text string) int {
	return terminalWidth.StringWidth(text)
}

func (a *App) paint(code, text string) string {
	if !a.screenMode {
		return text
	}
	return code + text + ansiReset
}

func (a *App) renderWelcome() {
	const (
		name  = "Basic Strategy Trainer"
		rules = "6D • H17 • DAS • BJ 3:2"
	)
	top := "─ BJ LINT "
	top += strings.Repeat("─", panelWidth-displayWidth(top))
	bodyWidth := displayWidth(" " + name + "  " + rules)
	padding := strings.Repeat(" ", panelWidth-bodyWidth)

	fmt.Fprintln(a.out, a.paint(ansiCyan+ansiBold, "╭"+top+"╮"))
	fmt.Fprintf(a.out, "│ %s  %s%s│\n", a.paint(ansiBold, name), a.paint(ansiDim, rules), padding)
	fmt.Fprintln(a.out, a.paint(ansiCyan+ansiBold, "╰"+strings.Repeat("─", panelWidth)+"╯"))
}

func (a *App) renderScreen(message string) {
	if a.screenMode {
		fmt.Fprint(a.out, "\x1b[2J\x1b[H")
	}
	a.renderWelcome()
	if message != "" {
		a.renderNotice(message)
	}
	a.renderState()
	// Keep go-prompt on its own line even after a compact round summary.
	fmt.Fprintln(a.out)
}

func (a *App) renderNotice(message string) {
	color := ansiBlue
	title := "LAST ACTION"
	switch {
	case strings.Contains(message, "WARNING"), strings.Contains(message, "deviates"):
		color, title = ansiYellow, "STRATEGY COACH"
		if strings.Contains(message, "Round complete") {
			title = "COACH + LAST ROUND"
		}
	case strings.HasPrefix(message, "Invalid"), strings.HasPrefix(message, "Unknown"), strings.HasPrefix(message, "Error"):
		color, title = ansiRed, "INPUT"
	case strings.Contains(message, "Round complete"):
		color, title = ansiGreen, "LAST ROUND"
	}
	top := "─ " + title + " "
	top += strings.Repeat("─", panelWidth-displayWidth(top))
	fmt.Fprintf(a.out, "\n%s\n", a.paint(color+ansiBold, "┌"+top+"┐"))
	for _, line := range strings.Split(strings.TrimSpace(message), "\n") {
		if line == "Round complete" || line == "" {
			continue
		}
		for _, wrapped := range wrapDisplayLine(line, panelWidth-2) {
			body := " " + wrapped
			body += strings.Repeat(" ", panelWidth-displayWidth(body))
			fmt.Fprintf(a.out, "%s%s%s\n", a.paint(color, "│"), body, a.paint(color, "│"))
		}
	}
	fmt.Fprintln(a.out, a.paint(color+ansiBold, "└"+strings.Repeat("─", panelWidth)+"┘"))
}

func wrapDisplayLine(line string, width int) []string {
	words := strings.Fields(line)
	if len(words) == 0 {
		return nil
	}
	lines := []string{words[0]}
	for _, word := range words[1:] {
		last := len(lines) - 1
		if displayWidth(lines[last])+1+displayWidth(word) <= width {
			lines[last] += " " + word
		} else {
			lines = append(lines, word)
		}
	}
	return lines
}

func suitSymbol(s blackjack.Suit) string {
	switch s {
	case blackjack.Clubs:
		return "♣"
	case blackjack.Diamonds:
		return "♦"
	case blackjack.Hearts:
		return "♥"
	case blackjack.Spades:
		return "♠"
	default:
		return "?"
	}
}

func (a *App) cardToken(c blackjack.Card) string {
	text := fmt.Sprintf("[ %-2s%s ]", c.Rank, suitSymbol(c.Suit))
	if c.Suit == blackjack.Hearts || c.Suit == blackjack.Diamonds {
		return a.paint(ansiRed+ansiBold, text)
	}
	return a.paint(ansiBold, text)
}

func (a *App) cards(h blackjack.Hand) string {
	parts := make([]string, len(h.Cards))
	for i, c := range h.Cards {
		parts[i] = a.cardToken(c)
	}
	return strings.Join(parts, " ")
}

func handTotal(h blackjack.Hand) string {
	total, soft := h.Value()
	switch {
	case h.IsBust():
		return fmt.Sprintf("BUST %d", total)
	case soft:
		return fmt.Sprintf("SOFT %d", total)
	default:
		return fmt.Sprintf("HARD %d", total)
	}
}

func (a *App) renderState() {
	if a.game.Phase == blackjack.PhaseRoundComplete {
		a.renderSettlement(a.out)
		return
	}
	hand := a.game.Current()
	if hand == nil {
		return
	}
	fmt.Fprintf(a.out, "\n  %s  %s    %s  %s    %s  %d/%d\n",
		a.paint(ansiDim, "BANKROLL"), a.paint(ansiGreen+ansiBold, a.game.DisplayBankroll().String()),
		a.paint(ansiDim, "BET"), a.paint(ansiBold, a.game.Rules.Bet.String()),
		a.paint(ansiDim, "SHOE"), a.game.Shoe.Remaining(), a.game.Shoe.Size())

	fmt.Fprintf(a.out, "\n  %s  %s %s\n", a.paint(ansiBlue+ansiBold, "DEALER  "), a.cardToken(a.game.Dealer.Cards[0]), a.paint(ansiDim, "[ ?  ]"))

	playerTitle := "PLAYER"
	if len(a.game.PlayerHands) > 1 {
		playerTitle = fmt.Sprintf("PLAYER  •  HAND %d OF %d", a.game.CurrentHand+1, len(a.game.PlayerHands))
	}
	fmt.Fprintf(a.out, "  %s  %s\n", a.paint(ansiBlue+ansiBold, fmt.Sprintf("%-8s", playerTitle)), a.cards(hand.Hand))
	totalColor := ansiBold
	if hand.Hand.IsBust() {
		totalColor = ansiRed + ansiBold
	} else if soft, _ := hand.Hand.Value(); soft == 21 {
		totalColor = ansiGreen + ansiBold
	}
	fmt.Fprintf(a.out, "  %s  %s\n", a.paint(ansiDim, "TOTAL   "), a.paint(totalColor, handTotal(hand.Hand)))

	if hand.Hand.IsBlackjack() {
		fmt.Fprintf(a.out, "\n  %s\n", a.paint(ansiGreen+ansiBold, "★ BLACKJACK"))
		return
	}
	if hand.Hand.IsBust() {
		return
	}

	recommendation := strategy.Recommend(hand.Hand, a.game.Dealer.Cards[0], a.game.Rules, a.game.AvailableActions())
	fmt.Fprintf(a.out, "  %s  %s\n", a.paint(ansiDim, "STRATEGY"), a.paint(ansiYellow+ansiBold, recommendation.Primary.String()))
	fmt.Fprintf(a.out, "\n  %s  %s\n", a.paint(ansiDim, "ACTIONS "), a.actionLabels(a.game.AvailableActions()))
	fmt.Fprintf(a.out, "  %s\n", a.paint(ansiDim, "↑/↓ select • Enter confirms • Type a key or command • Q quits"))
}

func (a *App) actionLabels(set blackjack.ActionSet) string {
	items := []struct {
		action blackjack.Action
		key    string
		label  string
	}{
		{blackjack.ActionHit, "H", "Hit"},
		{blackjack.ActionStand, "S", "Stand"},
		{blackjack.ActionDouble, "D", "Double"},
		{blackjack.ActionSplit, "P", "Split"},
	}
	labels := make([]string, 0, len(items)+1)
	for _, item := range items {
		if set.Has(item.action) {
			labels = append(labels, a.paint(ansiCyan+ansiBold, "["+item.key+"]")+" "+item.label)
		}
	}
	labels = append(labels, a.paint(ansiDim, "[Q]")+" Quit")
	return strings.Join(labels, "    ")
}

func (a *App) settlementText() string {
	var b bytes.Buffer
	a.renderSettlement(&b)
	return b.String()
}

func (a *App) renderSettlement(out io.Writer) {
	fmt.Fprintln(out, "Round complete")
	dealerTotal, _ := a.game.Dealer.Value()
	dealerLabel := fmt.Sprintf("Dealer %d", dealerTotal)
	if a.game.Dealer.IsBust() {
		dealerLabel += " (bust)"
	} else if a.game.Dealer.IsBlackjack() {
		dealerLabel += " (blackjack)"
	}
	if len(a.game.Results) == 1 {
		result := a.game.Results[0]
		delta := resultDelta(result.Profit)
		fmt.Fprintf(out, "Player %d  •  %s  •  %s %s\n", result.PlayerTotal, dealerLabel, result.Outcome, delta)
		fmt.Fprintf(out, "Bankroll %s", a.game.Bankroll)
		return
	}
	fmt.Fprintf(out, "%s", dealerLabel)
	for i, result := range a.game.Results {
		label := "Player"
		if len(a.game.Results) > 1 {
			label = fmt.Sprintf("Hand %d", i+1)
		}
		fmt.Fprintf(out, "\n%-8s %2d  %-9s %s", label, result.PlayerTotal, result.Outcome, resultDelta(result.Profit))
	}
	fmt.Fprintf(out, "\nBankroll %s", a.game.Bankroll)
}

func resultDelta(profit blackjack.Money) string {
	if profit > 0 {
		return "+" + profit.String()
	}
	if profit < 0 {
		return profit.String()
	}
	return ""
}
