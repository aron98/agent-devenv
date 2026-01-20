package cli

import (
	"fmt"
	"io"
)

const (
	exitOK             = 0
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

func (c *CLI) Run(args []string) int {
	if len(args) == 0 {
		c.printUsage(c.err)
		return exitUsage
	}

	cmd := args[0]
	switch cmd {
	case "help", "-h", "--help":
		c.printUsage(c.out)
		return exitOK
	case "init", "up", "down", "list", "status", "ports", "service", "agent", "stdio":
		return c.notImplemented(cmd)
	default:
		fmt.Fprintf(c.err, "unknown command: %s\n", cmd)
		c.printUsage(c.err)
		return exitUsage
	}
}

func (c *CLI) notImplemented(cmd string) int {
	fmt.Fprintf(c.err, "%s: not implemented\n", cmd)
	return exitNotImplemented
}

func (c *CLI) printUsage(out io.Writer) {
	fmt.Fprintln(out, "Usage: devenv <command> [options]")
	fmt.Fprintln(out, "Commands: init, up, down, list, status, ports, service, agent, stdio")
}
