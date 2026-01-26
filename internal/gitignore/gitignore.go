package gitignore

import (
	"fmt"
	"os"
	"strings"
)

func EnsureEntries(path string, entries []string) (bool, error) {
	if len(entries) == 0 {
		return false, nil
	}

	contents, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("read gitignore: %w", err)
	}

	normalized := strings.ReplaceAll(string(contents), "\r\n", "\n")
	lines := strings.Split(normalized, "\n")

	existing := make(map[string]struct{}, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		existing[trimmed] = struct{}{}
	}

	var missing []string
	for _, entry := range entries {
		trimmed := strings.TrimSpace(entry)
		if trimmed == "" {
			continue
		}
		if _, ok := existing[trimmed]; !ok {
			missing = append(missing, trimmed)
			existing[trimmed] = struct{}{}
		}
	}

	if len(missing) == 0 {
		return false, nil
	}

	var builder strings.Builder
	builder.WriteString(normalized)
	if len(normalized) > 0 && !strings.HasSuffix(normalized, "\n") {
		builder.WriteString("\n")
	}
	for _, entry := range missing {
		builder.WriteString(entry)
		builder.WriteString("\n")
	}

	if err := os.WriteFile(path, []byte(builder.String()), 0o644); err != nil {
		return false, fmt.Errorf("write gitignore: %w", err)
	}

	return true, nil
}
