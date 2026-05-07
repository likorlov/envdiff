package exporter

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Format represents the output format for exporting env vars.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
	FormatShell  Format = "shell"
)

// Export serializes a map of environment variables into the given format.
func Export(env map[string]string, format Format) (string, error) {
	switch format {
	case FormatDotenv:
		return exportDotenv(env), nil
	case FormatJSON:
		return exportJSON(env)
	case FormatShell:
		return exportShell(env), nil
	default:
		return "", fmt.Errorf("unsupported export format: %q", format)
	}
}

func sortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func exportDotenv(env map[string]string) string {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		v := env[k]
		if strings.ContainsAny(v, " \t\n#") {
			fmt.Fprintf(&sb, "%s=%q\n", k, v)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}
	return sb.String()
}

func exportJSON(env map[string]string) (string, error) {
	ordered := make(map[string]string, len(env))
	for k, v := range env {
		ordered[k] = v
	}
	b, err := json.MarshalIndent(ordered, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(b) + "\n", nil
}

func exportShell(env map[string]string) string {
	var sb strings.Builder
	for _, k := range sortedKeys(env) {
		v := env[k]
		fmt.Fprintf(&sb, "export %s=%q\n", k, v)
	}
	return sb.String()
}
