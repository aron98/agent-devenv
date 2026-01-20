package config

import "strings"

func NormalizeEnvName(input string) string {
	name := strings.ToLower(strings.TrimSpace(input))
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "--", "-")
	for strings.Contains(name, "--") {
		name = strings.ReplaceAll(name, "--", "-")
	}
	name = strings.Trim(name, "-")
	return name
}
