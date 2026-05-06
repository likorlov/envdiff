package formatter

import (
	"fmt"
	"strings"

	"github.com/envdiff/internal/differ"
)

// OutputFormat controls how diff results are rendered.
type OutputFormat string

const (
	FormatText  OutputFormat = "text"
	FormatDotenv OutputFormat = "dotenv"
	FormatJSON  OutputFormat = "json"
)

// FormatDiff renders a slice of differ.Entry in the requested format.
func FormatDiff(entries []differ.Entry, format OutputFormat) (string, error) {
	switch format {
	case FormatText:
		return formatText(entries), nil
	case FormatDotenv:
		return formatDotenv(entries), nil
	case FormatJSON:
		return formatJSON(entries), nil
	default:
		return "", fmt.Errorf("unknown format %q", format)
	}
}

func formatText(entries []differ.Entry) string {
	if len(entries) == 0 {
		return "No differences found.\n"
	}
	var sb strings.Builder
	for _, e := range entries {
		switch e.Kind {
		case differ.Added:
			fmt.Fprintf(&sb, "+ %s=%s\n", e.Key, e.ToValue)
		case differ.Removed:
			fmt.Fprintf(&sb, "- %s=%s\n", e.Key, e.FromValue)
		case differ.Changed:
			fmt.Fprintf(&sb, "~ %s: %s -> %s\n", e.Key, e.FromValue, e.ToValue)
		case differ.Unchanged:
			fmt.Fprintf(&sb, "  %s=%s\n", e.Key, e.FromValue)
		}
	}
	return sb.String()
}

func formatDotenv(entries []differ.Entry) string {
	var sb strings.Builder
	for _, e := range entries {
		switch e.Kind {
		case differ.Added, differ.Unchanged:
			fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.ToValue)
		case differ.Changed:
			fmt.Fprintf(&sb, "# changed from: %s\n%s=%s\n", e.FromValue, e.Key, e.ToValue)
		case differ.Removed:
			fmt.Fprintf(&sb, "# removed: %s\n", e.Key)
		}
	}
	return sb.String()
}

func formatJSON(entries []differ.Entry) string {
	if len(entries) == 0 {
		return "[]\n"
	}
	var sb strings.Builder
	sb.WriteString("[\n")
	for i, e := range entries {
		comma := ","
		if i == len(entries)-1 {
			comma = ""
		}
		fmt.Fprintf(&sb, "  {\"key\":%q,\"kind\":%q,\"from\":%q,\"to\":%q}%s\n",
			e.Key, e.Kind, e.FromValue, e.ToValue, comma)
	}
	sb.WriteString("]\n")
	return sb.String()
}
