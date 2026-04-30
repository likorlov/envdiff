package parser

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// EnvMap represents a parsed environment file as a map of key-value pairs.
type EnvMap map[string]string

// ParseFile reads and parses a .env file, returning an EnvMap.
// It skips blank lines and comments (lines starting with '#').
// It returns an error if the file cannot be opened or contains malformed lines.
func ParseFile(path string) (EnvMap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: cannot open file %q: %w", path, err)
	}
	defer f.Close()

	return parse(bufio.NewScanner(f), path)
}

// ParseString parses env content from a raw string, useful for testing.
func ParseString(content string) (EnvMap, error) {
	return parse(bufio.NewScanner(strings.NewReader(content)), "<string>")
}

func parse(scanner *bufio.Scanner, source string) (EnvMap, error) {
	env := make(EnvMap)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip blank lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Strip optional inline comments
		if idx := strings.Index(line, " #"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("parser: %s:%d: malformed line %q", source, lineNum, line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Strip surrounding quotes from value
		value = stripQuotes(value)

		if key == "" {
			return nil, fmt.Errorf("parser: %s:%d: empty key", source, lineNum)
		}

		env[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: scanning %s: %w", source, err)
	}

	return env, nil
}

func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
