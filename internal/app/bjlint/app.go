package bjlint

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	prompt "github.com/c-bata/go-prompt"
	"github.com/tkmru/bjlint-cli/internal/pkg/blackjack"
)

type App struct {
	game       *blackjack.Game
	out        io.Writer
	quit       bool
	finalShown bool
	screenMode bool
	message    string
}

const actionPrompt = "YOUR MOVE › "

func New() (*App, error) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	g, err := blackjack.NewGame(blackjack.DefaultRules(), rng)
	if err != nil {
		return nil, err
	}
	return &App{game: g, out: os.Stdout}, nil
}
func NewWithGame(g *blackjack.Game, out io.Writer) *App { return &App{game: g, out: out} }

func (a *App) Run() {
	terminal := captureTerminal()
	if terminal != nil {
		defer terminal.close()
	}

	a.enterScreen()
	a.message = a.completeRound()
	var parser *trackingParser
	for !a.quit && a.game.Phase != blackjack.PhaseGameOver {
		a.renderScreen(a.message)
		a.message = ""
		if parser == nil {
			parser = &trackingParser{
				ConsoleParser: prompt.NewStandardInputParser(),
				terminal:      terminal,
			}
		}
		parser.lastKey = prompt.NotDefined
		p := prompt.New(
			func(string) {},
			a.complete,
			prompt.OptionPrefix(actionPrompt),
			prompt.OptionPrefixTextColor(prompt.Cyan),
			prompt.OptionShowCompletionAtStart(),
			prompt.OptionCompletionOnDown(),
			prompt.OptionParser(parser),
		)
		input := p.Input()
		if terminal != nil {
			terminal.restore()
		}
		if parser.lastKey == prompt.ControlD {
			a.quit = true
			continue
		}
		a.Execute(input)
	}
	a.leaveScreen()
	a.renderFinal()
}

type trackingParser struct {
	prompt.ConsoleParser
	lastKey  prompt.Key
	terminal *terminalSnapshot
}

func (p *trackingParser) Read() ([]byte, error) {
	b, err := p.ConsoleParser.Read()
	if err == nil {
		p.lastKey = prompt.GetKey(b)
	}
	return b, err
}

func (p *trackingParser) TearDown() error {
	err := p.ConsoleParser.TearDown()
	if p.terminal != nil {
		p.terminal.restore()
	}
	return err
}

func (a *App) completeRound() string {
	var completed []string
	// A newly dealt round can immediately settle because of player/dealer
	// blackjack. Keep advancing until the player actually has a decision.
	for a.game.Phase == blackjack.PhaseRoundComplete {
		completed = append(completed, a.settlementText())
		if err := a.game.StartRound(); err != nil {
			a.quit = true
			completed = append(completed, fmt.Sprintf("Error: %v", err))
			break
		}
	}
	if a.game.Phase == blackjack.PhaseGameOver {
		completed = append(completed, fmt.Sprintf("Game over.\n\nFinal bankroll: %s", a.game.Bankroll))
	}
	return strings.Join(completed, "\n\n")
}

func (a *App) renderFinal() {
	if a.finalShown {
		return
	}
	if a.game.Phase == blackjack.PhaseGameOver {
		fmt.Fprintln(a.out, "Game over.")
	}
	if a.message != "" {
		fmt.Fprintf(a.out, "\n%s\n", a.message)
	}
	fmt.Fprintf(a.out, "\nFinal bankroll: %s\nThanks for playing.\n", a.game.DisplayBankroll())
	a.finalShown = true
}

func (a *App) enterScreen() {
	a.screenMode = true
}

func (a *App) leaveScreen() {
	if !a.screenMode {
		return
	}
	fmt.Fprint(a.out, "\x1b[2J\x1b[H")
	a.screenMode = false
}
