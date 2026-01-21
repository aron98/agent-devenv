package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/aron98/agent-devenv/internal/errs"
)

func TestRunHelp(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := New(stdout, stderr)

	if err := app.Run([]string{"help"}); err != nil {
		t.Fatalf("Run(help) returned error: %v", err)
	}
	if !strings.Contains(stdout.String(), "Usage: devenv") {
		t.Fatalf("stdout missing usage text")
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr should be empty")
	}
}

func TestRunNoArgs(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := New(stdout, stderr)

	err := app.Run([]string{})
	if errs.ExitCode(err) != 2 {
		t.Fatalf("ExitCode = %d, want 2", errs.ExitCode(err))
	}
	if !strings.Contains(err.Error(), "no command provided") {
		t.Fatalf("error missing message")
	}
	if !strings.Contains(stderr.String(), "Usage: devenv") {
		t.Fatalf("stderr missing usage text")
	}
}

func TestRunUnknownCommand(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := New(stdout, stderr)

	err := app.Run([]string{"nope"})
	if errs.ExitCode(err) != 2 {
		t.Fatalf("ExitCode = %d, want 2", errs.ExitCode(err))
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("error missing message")
	}
	if !strings.Contains(stderr.String(), "Usage: devenv") {
		t.Fatalf("stderr missing usage text")
	}
}

func TestRunNotImplemented(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	app := New(stdout, stderr)

	err := app.Run([]string{"init"})
	if errs.ExitCode(err) != 1 {
		t.Fatalf("ExitCode = %d, want 1", errs.ExitCode(err))
	}
	if !strings.Contains(err.Error(), "not implemented") {
		t.Fatalf("error missing message")
	}
}
