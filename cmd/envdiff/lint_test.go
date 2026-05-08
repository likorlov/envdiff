package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvForLint(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return p
}

func TestRunLint_Clean(t *testing.T) {
	p := writeTempEnvForLint(t, "FOO=bar\nBAR=baz\n")
	if err := runLint(p, false); err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestRunLint_DuplicateKey_Strict(t *testing.T) {
	p := writeTempEnvForLint(t, "FOO=1\nFOO=2\n")
	if err := runLint(p, true); err == nil {
		t.Error("expected error for duplicate key in strict mode")
	}
}

func TestRunLint_DuplicateKey_NonStrict(t *testing.T) {
	// duplicate-key is still severe even without strict flag
	p := writeTempEnvForLint(t, "FOO=1\nFOO=2\n")
	if err := runLint(p, false); err == nil {
		t.Error("expected error for severe duplicate-key violation")
	}
}

func TestRunLint_TrailingSpace_NonStrict(t *testing.T) {
	// trailing space alone is not severe; non-strict should succeed
	p := writeTempEnvForLint(t, "FOO=bar   \n")
	if err := runLint(p, false); err != nil {
		t.Errorf("expected no error in non-strict mode for trailing space, got: %v", err)
	}
}

func TestRunLint_TrailingSpace_Strict(t *testing.T) {
	p := writeTempEnvForLint(t, "FOO=bar   \n")
	if err := runLint(p, true); err == nil {
		t.Error("expected error in strict mode for trailing space")
	}
}

func TestRunLint_MissingFile(t *testing.T) {
	if err := runLint("/nonexistent/.env", false); err == nil {
		t.Error("expected error for missing file")
	}
}
