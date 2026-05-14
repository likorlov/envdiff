package profiler

import (
	"strings"
	"testing"
)

var sampleProfile = Profile{
	Name:     "production",
	Required: []string{"DATABASE_URL", "SECRET_KEY", "PORT"},
	Optional: []string{"LOG_LEVEL", "DEBUG", "TIMEOUT"},
}

func TestCheck_AllPresent(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"SECRET_KEY":   "abc123",
		"PORT":         "8080",
	}
	violations := Check(env, sampleProfile)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestCheck_MissingRequired(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
	}
	violations := Check(env, sampleProfile)
	if len(violations) != 2 {
		t.Errorf("expected 2 violations, got %d", len(violations))
	}
}

func TestCheck_EmptyRequiredValue(t *testing.T) {
	env := map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"SECRET_KEY":   "   ",
		"PORT":         "8080",
	}
	violations := Check(env, sampleProfile)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "SECRET_KEY" {
		t.Errorf("expected violation for SECRET_KEY, got %s", violations[0].Key)
	}
}

func TestCoverage_Full(t *testing.T) {
	env := map[string]string{
		"LOG_LEVEL": "info",
		"DEBUG":     "false",
		"TIMEOUT":   "30",
	}
	cov := Coverage(env, sampleProfile)
	if cov != 100.0 {
		t.Errorf("expected 100.0, got %.2f", cov)
	}
}

func TestCoverage_Partial(t *testing.T) {
	env := map[string]string{
		"LOG_LEVEL": "info",
	}
	cov := Coverage(env, sampleProfile)
	expected := 100.0 / 3.0
	if cov < expected-0.01 || cov > expected+0.01 {
		t.Errorf("expected ~%.2f, got %.2f", expected, cov)
	}
}

func TestCoverage_NoOptional(t *testing.T) {
	p := Profile{Name: "minimal", Required: []string{"FOO"}}
	cov := Coverage(map[string]string{"FOO": "bar"}, p)
	if cov != 100.0 {
		t.Errorf("expected 100.0 for no optional keys, got %.2f", cov)
	}
}

func TestFormatViolations_Empty(t *testing.T) {
	out := FormatViolations(nil)
	if out != "no profile violations" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatViolations_Sorted(t *testing.T) {
	v := []Violation{
		{Key: "Z_KEY", Message: "required key is missing"},
		{Key: "A_KEY", Message: "required key is missing"},
	}
	out := FormatViolations(v)
	if !strings.Contains(out, "A_KEY") || !strings.Contains(out, "Z_KEY") {
		t.Errorf("expected both keys in output: %q", out)
	}
	if strings.Index(out, "A_KEY") > strings.Index(out, "Z_KEY") {
		t.Error("expected A_KEY before Z_KEY in sorted output")
	}
}
