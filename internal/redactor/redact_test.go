package redactor_test

import (
	"regexp"
	"testing"

	"github.com/yourorg/envdiff/internal/redactor"
)

func TestRedact_SensitiveKeysAreRedacted(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"API_KEY":     "abc123",
		"APP_NAME":    "envdiff",
		"AUTH_TOKEN":  "tok_xyz",
	}

	result := redactor.Redact(env, redactor.Options{})

	if result["DB_PASSWORD"] != "***REDACTED***" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "***REDACTED***" {
		t.Errorf("expected API_KEY to be redacted, got %q", result["API_KEY"])
	}
	if result["AUTH_TOKEN"] != "***REDACTED***" {
		t.Errorf("expected AUTH_TOKEN to be redacted, got %q", result["AUTH_TOKEN"])
	}
	if result["APP_NAME"] != "envdiff" {
		t.Errorf("expected APP_NAME to be unchanged, got %q", result["APP_NAME"])
	}
}

func TestRedact_CustomPlaceholder(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "topsecret"}
	result := redactor.Redact(env, redactor.Options{Placeholder: "<hidden>"})
	if result["SECRET_KEY"] != "<hidden>" {
		t.Errorf("expected custom placeholder, got %q", result["SECRET_KEY"])
	}
}

func TestRedact_ExtraPatterns(t *testing.T) {
	env := map[string]string{
		"STRIPE_PUBLISHABLE": "pk_live_abc",
		"APP_PORT":           "8080",
	}
	extra := []*regexp.Regexp{regexp.MustCompile(`(?i)stripe`)}
	result := redactor.Redact(env, redactor.Options{ExtraPatterns: extra})
	if result["STRIPE_PUBLISHABLE"] != "***REDACTED***" {
		t.Errorf("expected STRIPE_PUBLISHABLE to be redacted, got %q", result["STRIPE_PUBLISHABLE"])
	}
	if result["APP_PORT"] != "8080" {
		t.Errorf("expected APP_PORT unchanged, got %q", result["APP_PORT"])
	}
}

func TestRedact_OriginalEnvUnmodified(t *testing.T) {
	env := map[string]string{"DB_PASSWORD": "original"}
	redactor.Redact(env, redactor.Options{})
	if env["DB_PASSWORD"] != "original" {
		t.Error("Redact must not modify the original map")
	}
}

func TestRedact_EmptyEnv(t *testing.T) {
	result := redactor.Redact(map[string]string{}, redactor.Options{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestIsSensitiveKey(t *testing.T) {
	cases := []struct {
		key       string
		sensitive bool
	}{
		{"DB_PASSWORD", true},
		{"api_key", true},
		{"PRIVATE_KEY", true},
		{"APP_NAME", false},
		{"PORT", false},
		{"DATABASE_URL", true},
	}
	for _, tc := range cases {
		got := redactor.IsSensitiveKey(tc.key)
		if got != tc.sensitive {
			t.Errorf("IsSensitiveKey(%q) = %v, want %v", tc.key, got, tc.sensitive)
		}
	}
}
