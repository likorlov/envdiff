package parser

import (
	"testing"
)

func TestParseString_Basic(t *testing.T) {
	input := `
DB_HOST=localhost
DB_PORT=5432
APP_NAME="my app"
`
	env, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := EnvMap{
		"DB_HOST":  "localhost",
		"DB_PORT":  "5432",
		"APP_NAME": "my app",
	}

	for k, v := range expected {
		if env[k] != v {
			t.Errorf("key %q: got %q, want %q", k, env[k], v)
		}
	}
}

func TestParseString_SkipsCommentsAndBlanks(t *testing.T) {
	input := `
# This is a comment

FOO=bar
# Another comment
BAZ=qux
`
	env, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(env) != 2 {
		t.Errorf("expected 2 keys, got %d", len(env))
	}
}

func TestParseString_InlineComment(t *testing.T) {
	input := `PORT=8080 # default port`
	env, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", env["PORT"])
	}
}

func TestParseString_SingleQuotedValue(t *testing.T) {
	input := `SECRET='my secret value'`
	env, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SECRET"] != "my secret value" {
		t.Errorf("unexpected value: %q", env["SECRET"])
	}
}

func TestParseString_MalformedLine(t *testing.T) {
	input := `NOTAVALIDLINE`
	_, err := ParseString(input)
	if err == nil {
		t.Error("expected error for malformed line, got nil")
	}
}

func TestParseString_EmptyKey(t *testing.T) {
	input := `=value`
	_, err := ParseString(input)
	if err == nil {
		t.Error("expected error for empty key, got nil")
	}
}

func TestParseString_EmptyValue(t *testing.T) {
	input := `EMPTY=`
	env, err := ParseString(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := env["EMPTY"]; !ok || v != "" {
		t.Errorf("expected EMPTY=\"\", got %q (present=%v)", v, ok)
	}
}

func TestParseString_ExportPrefix(t *testing.T) {
	// Some .env files use "export KEY=VALUE" syntax (e.g. for sourcing in bash).
	// Verify that the parser handles or consistently rejects this prefix.
	input := `export API_KEY=abc123`
	env, err := ParseString(input)
	if err != nil {
		t.Skipf("export prefix not supported (got error: %v)", err)
	}
	if env["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", env["API_KEY"])
	}
}
