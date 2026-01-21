package cli

import (
	"fmt"
	"io"

	"github.com/aron98/agent-devenv/internal/errs"
)

const (
	exitNotImplemented = 1
	exitUsage          = 2
)

type CLI struct {
	out io.Writer
	err io.Writer
}

func New(out io.Writer, err io.Writer) *CLI {
	return &CLI{out: out, err: err}
}

func (c *CLI) Run(args []string) error {
	if len(args) == 0 {
		c.printUsage(c.err)
		return errs.New(exitUsage, "no command provided", "run 'devenv help' for usage")
	}

	cmd := args[0]
	switch cmd {
	case "help", "-h", "--help":
		c.printUsage(c.out)
		return nil
	case "init", "up", "down", "list", "status", "ports", "service", "agent", "stdio":
		return c.notImplemented(cmd)
	default:
		c.printUsage(c.err)
		return errs.New(exitUsage, fmt.Sprintf("unknown command: %s", cmd), "run 'devenv help' for usage")
	}
}

func (c *CLI) notImplemented(cmd string) error {
	return errs.New(exitNotImplemented, fmt.Sprintf("%s: not implemented", cmd), "see IMPLEMENTATION_PLAN.md for status")
}

func (c *CLI) printUsage(out io.Writer) {
	fmt.Fprintln(out, "Usage: devenv <command> [options]")
	fmt.Fprintln(out, "Commands: init, up, down, list, status, ports, service, agent, stdio")
}
