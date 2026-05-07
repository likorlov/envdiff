package validator

import (
	"strings"
	"testing"
)

func TestValidate_ValidEnv(t *testing.T) {
	env := map[string]string{
		"APP_ENV":  "production",
		"DB_HOST":  "localhost",
		"PORT":     "8080",
	}
	violations := Validate(env, DefaultRules)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %v", len(violations), violations)
	}
}

func TestValidate_LowercaseKey(t *testing.T) {
	env := map[string]string{
		"app_env": "production",
	}
	violations := Validate(env, []Rule{KeyFormat})
	if len(violations) == 0 {
		t.Fatal("expected violation for lowercase key, got none")
	}
	if violations[0].Rule != "key_format" {
		t.Errorf("expected rule key_format, got %s", violations[0].Rule)
	}
}

func TestValidate_KeyWithSpaces(t *testing.T) {
	env := map[string]string{
		"MY KEY": "value",
	}
	violations := Validate(env, []Rule{NoWhitespaceKeys})
	if len(violations) == 0 {
		t.Fatal("expected violation for key with spaces, got none")
	}
}

func TestValidate_EmptyValue(t *testing.T) {
	env := map[string]string{
		"API_KEY": "",
	}
	violations := Validate(env, []Rule{NoEmptyValues})
	if len(violations) == 0 {
		t.Fatal("expected violation for empty value, got none")
	}
	if violations[0].Rule != "no_empty_values" {
		t.Errorf("expected rule no_empty_values, got %s", violations[0].Rule)
	}
}

func TestValidate_MultipleViolations(t *testing.T) {
	env := map[string]string{
		"bad key":  "",
		"GOOD_KEY": "ok",
	}
	rules := []Rule{NoWhitespaceKeys, NoEmptyValues}
	violations := Validate(env, rules)
	// "bad key" violates NoWhitespaceKeys; "" violates NoEmptyValues
	if len(violations) < 2 {
		t.Errorf("expected at least 2 violations, got %d", len(violations))
	}
}

func TestFormatViolations_NoViolations(t *testing.T) {
	out := FormatViolations(nil)
	if out != "no validation issues found" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatViolations_WithViolations(t *testing.T) {
	violations := []Violation{
		{Key: "bad_key", Value: "x", Rule: "key_format", Message: "key must be uppercase alphanumeric with underscores"},
	}
	out := FormatViolations(violations)
	if !strings.Contains(out, "key_format") {
		t.Errorf("expected rule name in output, got: %s", out)
	}
	if !strings.Contains(out, "bad_key") {
		t.Errorf("expected key name in output, got: %s", out)
	}
}
