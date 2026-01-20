package errs

import (
	"errors"
	"fmt"
)

type CLIError struct {
	Code    int
	Message string
	Hint    string
}

func (err CLIError) Error() string {
	if err.Hint == "" {
		return err.Message
	}
	return fmt.Sprintf("%s\nhint: %s", err.Message, err.Hint)
}

func New(code int, message string, hint string) error {
	return CLIError{Code: code, Message: message, Hint: hint}
}

func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	var cliErr CLIError
	if errors.As(err, &cliErr) {
		if cliErr.Code == 0 {
			return 1
		}
		return cliErr.Code
	}

	return 1
}
