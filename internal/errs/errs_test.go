package errs

import (
	"errors"
	"strings"
	"testing"
)

func TestExitCodeNil(t *testing.T) {
	if ExitCode(nil) != 0 {
		t.Fatalf("ExitCode(nil) = %d, want 0", ExitCode(nil))
	}
}

func TestExitCodeCLIError(t *testing.T) {
	err := New(2, "bad input", "fix it")
	if ExitCode(err) != 2 {
		t.Fatalf("ExitCode(CLIError) = %d, want 2", ExitCode(err))
	}
}

func TestExitCodeCLIErrorZero(t *testing.T) {
	err := CLIError{Message: "missing code"}
	if ExitCode(err) != 1 {
		t.Fatalf("ExitCode(CLIError) = %d, want 1", ExitCode(err))
	}
}

func TestExitCodeUnknown(t *testing.T) {
	err := errors.New("boom")
	if ExitCode(err) != 1 {
		t.Fatalf("ExitCode(unknown) = %d, want 1", ExitCode(err))
	}
}

func TestCLIErrorIncludesHint(t *testing.T) {
	err := New(2, "bad input", "fix it")
	if !strings.Contains(err.Error(), "hint: fix it") {
		t.Fatalf("error missing hint")
	}
}
