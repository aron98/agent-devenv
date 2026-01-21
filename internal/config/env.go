package config

import (
	"regexp"
	"strings"
)

var envNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-_]{0,62}$`)

func NormalizeEnvName(input string) string {
	name := strings.ToLower(strings.TrimSpace(input))
	name = strings.ReplaceAll(name, " ", "-")
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}
	name = strings.Trim(name, "-")
	return name
}

func ValidateEnvName(name string) bool {
	return envNamePattern.MatchString(name)
}
