# bjlint-cli

A blackjack CLI trainer that warns you when your play deviates from basic strategy.

Once installed, `bjlint` works entirely offline, so you can practice on the flight to Las Vegas without Wi-Fi. Because it runs in your terminal, it also blends in with command-line work—perfect for playing a discreet hand at the office without anyone noticing.

## Overview

`bjlint` lets you play continuous blackjack rounds while checking every legal choice against basic strategy. A deviation produces a warning but is still executed, so the trainer never takes the decision away from you.

## Features

- 6-deck blackjack
- H17
- DAS
- Split up to 4 hands
- Blackjack pays 3:2
- $10 fixed bet
- $100 starting bankroll
- Basic Strategy lint
- Interactive completion powered by go-prompt
- In-place color dashboard with card suits and round summaries

## Installation

```bash
$ go install github.com/tkmru/bjlint-cli/cmd/bjlint@latest
$ bjlint
```

For local development:

```bash
$ go run ./cmd/bjlint
```

## Usage

Available actions are shown automatically. Use the up and down arrow keys to select one and press Enter to play it. Typing remains available: commands are case-insensitive, and you can use `h`/`hit`, `s`/`stand`, `d`/`double`, `p`/`split`, or `q`/`quit`/`exit`.

The game redraws the current table in place after every command instead of appending each state to the terminal. The previous round's result remains visible above the newly dealt hand until the next action.

```text
$ bjlint

bjlint-cli
Blackjack Basic Strategy Trainer

Rules:
  6 decks / H17 / DAS / Blackjack 3:2

Bankroll: $100.00
Bet:       $10.00

Dealer:
  [6] [?]

Player:
  [10] [6]

Total: 16
Basic Strategy: STAND

> hit

⚠ WARNING: Basic Strategy recommends STAND.
You chose HIT.

Player:
  [10] [6] [4]

Total: 20
Basic Strategy: STAND

> stand
```

## Rules

The game uses six decks, dealer hits soft 17, double after split, no surrender, no insurance, and a 3:2 natural-blackjack payout. The fixed bet is $10 from a starting bankroll of $100. Resplitting is allowed up to four hands. Split aces receive one card each and cannot be hit or resplit. The shoe is replaced between rounds after 75% penetration; it is never replaced during a round.

## Basic Strategy

The table is specifically for 6D / H17 / DAS / No Surrender. Pair, soft-total, and hard-total decisions are maintained separately. When a table entry says “double, otherwise hit/stand” or “split, otherwise hit,” the recommendation automatically uses the fallback when the primary action is unavailable.

## Strategy Lint

A legal choice that differs from the effective recommendation displays a warning and then proceeds. An illegal choice—such as doubling a three-card hand—is rejected without changing game state.

## Interactive Completion

Currently available actions are shown as soon as the prompt opens. Press Down to begin selecting, use Up and Down to move, and press Enter to execute the selected action. You can also type a command or use Tab completion. Completion and execution both use the domain model's available-action set, so hiding an option is not used as a substitute for validation.

## Development

The executable entry point is in `cmd/bjlint`, CLI integration is in `internal/app/bjlint`, and dependency-free game and strategy logic are in `internal/pkg`.

## Testing

```bash
$ go fmt ./...
$ go vet ./...
$ go test ./...
$ go build ./cmd/bjlint
```

Unit tests cover hand values and aces, blackjack detection, H17 dealer play, shoe composition and cut-card behavior, bankroll settlement, double/split rules, strategy tables and fallbacks, parsing, lint behavior, and completion suggestions.

## License

Beerware License. See [LICENSE](LICENSE).
