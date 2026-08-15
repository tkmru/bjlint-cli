package main

import (
	"fmt"
	"os"

	"github.com/tkmru/bjlint-cli/internal/app/bjlint"
)

func main() {
	app, err := bjlint.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "bjlint:", err)
		os.Exit(1)
	}
	app.Run()
}
