package interpolator

import (
	"os"
	"testing"
)

func TestInterpolate_NoReferences(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	got, err := Interpolate(env, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Errorf("expected unchanged values, got %v", got)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "URL": "http://${HOST}:8080"}
	got, err := Interpolate(env, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["URL"] != "http://localhost:8080" {
		t.Errorf("expected expanded URL, got %q", got["URL"])
	}
}

func TestInterpolate_BareStyle(t *testing.T) {
	env := map[string]string{"NAME": "world", "GREETING": "hello $NAME"}
	got, err := Interpolate(env, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["GREETING"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", got["GREETING"])
	}
}

func TestInterpolate_MissingVariable_LeavesPlaceholder(t *testing.T) {
	env := map[string]string{"VAL": "${UNDEFINED}"}
	got, err := Interpolate(env, Options{FailOnMissing: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["VAL"] != "${UNDEFINED}" {
		t.Errorf("expected placeholder preserved, got %q", got["VAL"])
	}
}

func TestInterpolate_FailOnMissing(t *testing.T) {
	env := map[string]string{"VAL": "${MISSING}"}
	_, err := Interpolate(env, Options{FailOnMissing: true})
	if err == nil {
		t.Fatal("expected error for missing variable, got nil")
	}
}

func TestInterpolate_FallbackToOS(t *testing.T) {
	os.Setenv("OS_VAR", "from-os")
	defer os.Unsetenv("OS_VAR")

	env := map[string]string{"COMPUTED": "${OS_VAR}-suffix"}
	got, err := Interpolate(env, Options{FallbackToOS: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["COMPUTED"] != "from-os-suffix" {
		t.Errorf("expected 'from-os-suffix', got %q", got["COMPUTED"])
	}
}

func TestInterpolate_OriginalUnmodified(t *testing.T) {
	env := map[string]string{"A": "1", "B": "${A}"}
	_, err := Interpolate(env, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["B"] != "${A}" {
		t.Errorf("original map was modified")
	}
}

func TestInterpolate_MultipleRefsInValue(t *testing.T) {
	env := map[string]string{
		"PROTO": "https",
		"HOST":  "example.com",
		"PORT":  "443",
		"URL":   "${PROTO}://${HOST}:${PORT}",
	}
	got, err := Interpolate(env, Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["URL"] != "https://example.com:443" {
		t.Errorf("expected full URL, got %q", got["URL"])
	}
}
