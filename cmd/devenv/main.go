package main

import (
	"os"

	"github.com/aron98/agent-devenv/internal/cli"
)

func main() {
	app := cli.New(os.Stdout, os.Stderr)
	os.Exit(app.Run(os.Args[1:]))
}
