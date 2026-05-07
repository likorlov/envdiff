package exporter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/exporter"
)

var sampleEnv = map[string]string{
	"APP_ENV":  "production",
	"DB_HOST":  "localhost",
	"DB_PASS":  "s3cr3t p@ss",
	"LOG_LEVEL": "info",
}

func TestExport_Dotenv(t *testing.T) {
	out, err := exporter.Export(sampleEnv, exporter.FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got:\n%s", out)
	}
	// value with space should be quoted
	if !strings.Contains(out, `DB_PASS="s3cr3t p@ss"`) {
		t.Errorf("expected quoted DB_PASS in output, got:\n%s", out)
	}
}

func TestExport_JSON(t *testing.T) {
	out, err := exporter.Export(sampleEnv, exporter.FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"APP_ENV": "production"`) {
		t.Errorf("expected JSON key in output, got:\n%s", out)
	}
	if !strings.HasSuffix(strings.TrimSpace(out), "}") {
		t.Errorf("expected JSON to end with }, got:\n%s", out)
	}
}

func TestExport_Shell(t *testing.T) {
	out, err := exporter.Export(sampleEnv, exporter.FormatShell)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export APP_ENV=") {
		t.Errorf("expected export statement in output, got:\n%s", out)
	}
	if !strings.Contains(out, "export DB_HOST=") {
		t.Errorf("expected export DB_HOST in output, got:\n%s", out)
	}
}

func TestExport_UnknownFormat(t *testing.T) {
	_, err := exporter.Export(sampleEnv, exporter.Format("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format, got nil")
	}
}

func TestExport_EmptyEnv(t *testing.T) {
	for _, fmt := range []exporter.Format{exporter.FormatDotenv, exporter.FormatJSON, exporter.FormatShell} {
		out, err := exporter.Export(map[string]string{}, fmt)
		if err != nil {
			t.Errorf("format %s: unexpected error: %v", fmt, err)
		}
		if fmt == exporter.FormatDotenv && out != "" {
			t.Errorf("format %s: expected empty output, got %q", fmt, out)
		}
	}
}
