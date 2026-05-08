package linter

import (
	"strings"
	"testing"
)

func violations(t *testing.T, lines []string) []Violation {
	t.Helper()
	return Lint(lines)
}

func TestLint_Clean(t *testing.T) {
	v := violations(t, []string{"FOO=bar", "BAR_BAZ=123"})
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestLint_DuplicateKey(t *testing.T) {
	v := violations(t, []string{"FOO=1", "FOO=2"})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Rule != RuleDuplicateKey {
		t.Errorf("expected rule %s, got %s", RuleDuplicateKey, v[0].Rule)
	}
	if v[0].Line != 2 {
		t.Errorf("expected line 2, got %d", v[0].Line)
	}
}

func TestLint_TrailingSpace(t *testing.T) {
	v := violations(t, []string{"FOO=bar   "})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Rule != RuleTrailingSpace {
		t.Errorf("expected rule %s, got %s", RuleTrailingSpace, v[0].Rule)
	}
}

func TestLint_KeyNamingConvention(t *testing.T) {
	v := violations(t, []string{"foo_bar=value"})
	hasNaming := false
	for _, vv := range v {
		if vv.Rule == RuleKeyNamingConvention {
			hasNaming = true
		}
	}
	if !hasNaming {
		t.Error("expected key-naming-convention violation")
	}
}

func TestLint_NoValue(t *testing.T) {
	v := violations(t, []string{"FOO="})
	hasNoVal := false
	for _, vv := range v {
		if vv.Rule == RuleNoValue {
			hasNoVal = true
		}
	}
	if !hasNoVal {
		t.Error("expected no-value violation")
	}
}

func TestLint_SingleQuotePreference(t *testing.T) {
	v := violations(t, []string{"FOO='bar'"})
	hasQuote := false
	for _, vv := range v {
		if vv.Rule == RuleQuotingStyle {
			hasQuote = true
		}
	}
	if !hasQuote {
		t.Error("expected quoting-style violation")
	}
}

func TestLint_SkipsComments(t *testing.T) {
	v := violations(t, []string{"# this is a comment", "", "FOO=bar"})
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestFormatViolations_Empty(t *testing.T) {
	out := FormatViolations(nil)
	if out != "no lint violations found" {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatViolations_NonEmpty(t *testing.T) {
	v := []Violation{{Line: 1, Key: "foo", Rule: RuleDuplicateKey, Message: "dup"}}
	out := FormatViolations(v)
	if !strings.Contains(out, "duplicate-key") {
		t.Errorf("expected rule name in output, got: %q", out)
	}
}
