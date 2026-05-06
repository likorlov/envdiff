package formatter

import (
	"strings"
	"testing"

	"github.com/envdiff/internal/differ"
)

var sampleEntries = []differ.Entry{
	{Key: "APP_ENV", Kind: differ.Unchanged, FromValue: "production", ToValue: "production"},
	{Key: "DB_HOST", Kind: differ.Changed, FromValue: "localhost", ToValue: "db.prod.internal"},
	{Key: "NEW_KEY", Kind: differ.Added, FromValue: "", ToValue: "new_value"},
	{Key: "OLD_KEY", Kind: differ.Removed, FromValue: "old_value", ToValue: ""},
}

func TestFormatDiff_Text(t *testing.T) {
	out, err := FormatDiff(sampleEntries, FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "+ NEW_KEY=new_value") {
		t.Errorf("expected added line, got:\n%s", out)
	}
	if !strings.Contains(out, "- OLD_KEY=old_value") {
		t.Errorf("expected removed line, got:\n%s", out)
	}
	if !strings.Contains(out, "~ DB_HOST: localhost -> db.prod.internal") {
		t.Errorf("expected changed line, got:\n%s", out)
	}
	if !strings.Contains(out, "  APP_ENV=production") {
		t.Errorf("expected unchanged line, got:\n%s", out)
	}
}

func TestFormatDiff_Text_NoDiff(t *testing.T) {
	out, err := FormatDiff([]differ.Entry{}, FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "No differences found.") {
		t.Errorf("expected no-diff message, got: %s", out)
	}
}

func TestFormatDiff_Dotenv(t *testing.T) {
	out, err := FormatDiff(sampleEntries, FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "NEW_KEY=new_value") {
		t.Errorf("expected added key in dotenv output, got:\n%s", out)
	}
	if !strings.Contains(out, "# removed: OLD_KEY") {
		t.Errorf("expected removed comment in dotenv output, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST=db.prod.internal") {
		t.Errorf("expected updated value in dotenv output, got:\n%s", out)
	}
}

func TestFormatDiff_JSON(t *testing.T) {
	out, err := FormatDiff(sampleEntries, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(out, "[") {
		t.Errorf("expected JSON array, got: %s", out)
	}
	if !strings.Contains(out, `"key":"NEW_KEY"`) {
		t.Errorf("expected NEW_KEY in JSON output, got:\n%s", out)
	}
}

func TestFormatDiff_JSON_Empty(t *testing.T) {
	out, err := FormatDiff([]differ.Entry{}, FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(out) != "[]" {
		t.Errorf("expected empty JSON array, got: %s", out)
	}
}

func TestFormatDiff_UnknownFormat(t *testing.T) {
	_, err := FormatDiff(sampleEntries, OutputFormat("xml"))
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}
