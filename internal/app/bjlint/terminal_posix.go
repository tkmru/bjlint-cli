//go:build !windows

package bjlint

import (
	"os"

	"github.com/pkg/term/termios"
	"golang.org/x/sys/unix"
)

// terminalSnapshot keeps an independent copy of the terminal settings from
// before go-prompt switches /dev/tty to raw mode. go-prompt v0.2.6 mutates its
// saved termios value in place, so its own TearDown cannot reliably restore
// echo and canonical input.
type terminalSnapshot struct {
	tty      *os.File
	original unix.Termios
}

func captureTerminal() *terminalSnapshot {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return nil
	}
	state, err := termios.Tcgetattr(tty.Fd())
	if err != nil {
		_ = tty.Close()
		return nil
	}
	return &terminalSnapshot{tty: tty, original: *state}
}

func (t *terminalSnapshot) restore() {
	_ = termios.Tcsetattr(t.tty.Fd(), termios.TCSANOW, &t.original)
}

func (t *terminalSnapshot) close() {
	t.restore()
	_ = t.tty.Close()
}
