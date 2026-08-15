//go:build windows

package bjlint

type terminalSnapshot struct{}

func captureTerminal() *terminalSnapshot { return nil }
func (t *terminalSnapshot) restore()     {}
func (t *terminalSnapshot) close()       {}
