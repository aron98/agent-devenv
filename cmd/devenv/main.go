package main

import (
	"os"

	"github.com/aron98/agent-devenv/internal/cli"
	"github.com/aron98/agent-devenv/internal/errs"
	"github.com/aron98/agent-devenv/internal/logging"
)

func main() {
	app := cli.New(os.Stdout, os.Stderr)
	logger := logging.Default()

	err := app.Run(os.Args[1:])
	if err != nil {
		logger.Error(err.Error())
	}

	os.Exit(errs.ExitCode(err))
}
