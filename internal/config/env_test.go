package config

import (
	"strings"
	"testing"
)

func TestNormalizeEnvName(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{input: "Feature Login", expected: "feature-login"},
		{input: "  My--Env  ", expected: "my-env"},
		{input: "multi   space", expected: "multi-space"},
		{input: "--Name--", expected: "name"},
	}

	for _, testCase := range cases {
		t.Run(testCase.input, func(t *testing.T) {
			got := NormalizeEnvName(testCase.input)
			if got != testCase.expected {
				t.Fatalf("NormalizeEnvName(%q) = %q, want %q", testCase.input, got, testCase.expected)
			}
		})
	}
}

func TestValidateEnvName(t *testing.T) {
	valid := []string{
		"a",
		"abc-123",
		"a_b",
		"a-b_c",
	}
	invalid := []string{
		"",
		"A",
		"a b",
		"-bad",
		strings.Repeat("a", 64),
	}

	for _, name := range valid {
		t.Run("valid_"+name, func(t *testing.T) {
			if !ValidateEnvName(name) {
				t.Fatalf("ValidateEnvName(%q) = false, want true", name)
			}
		})
	}

	for _, name := range invalid {
		t.Run("invalid_"+name, func(t *testing.T) {
			if ValidateEnvName(name) {
				t.Fatalf("ValidateEnvName(%q) = true, want false", name)
			}
		})
	}
}
