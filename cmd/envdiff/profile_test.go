package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnvForProfile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp env: %v", err)
	}
	return path
}

func TestRunProfile_AllPresent(t *testing.T) {
	path := writeTempEnvForProfile(t, "DATABASE_URL=postgres://localhost\nSECRET_KEY=abc\nPORT=8080\n")
	err := runProfile(path, "prod", []string{"DATABASE_URL", "SECRET_KEY", "PORT"}, nil, false)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestRunProfile_MissingRequired(t *testing.T) {
	path := writeTempEnvForProfile(t, "DATABASE_URL=postgres://localhost\n")
	err := runProfile(path, "prod", []string{"DATABASE_URL", "SECRET_KEY"}, nil, false)
	if err == nil {
		t.Error("expected error for missing required key")
	}
	if !strings.Contains(err.Error(), "SECRET_KEY") {
		t.Errorf("expected error to mention SECRET_KEY, got: %v", err)
	}
}

func TestRunProfile_JSONOutput(t *testing.T) {
	path := writeTempEnvForProfile(t, "FOO=bar\n")
	// Redirect stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	_ = runProfile(path, "test", []string{"FOO"}, []string{"BAR"}, true)

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	tmp := make([]byte, 1024)
	for {
		n, err := r.Read(tmp)
		buf.Write(tmp[:n])
		if err != nil {
			break
		}
	}
	out := buf.String()
	if !strings.Contains(out, `"profile"`) {
		t.Errorf("expected JSON with profile field, got: %s", out)
	}
	if !strings.Contains(out, `"violations"`) {
		t.Errorf("expected JSON with violations field, got: %s", out)
	}
}

func TestRunProfile_MissingFile(t *testing.T) {
	err := runProfile("/nonexistent/.env", "prod", []string{"FOO"}, nil, false)
	if err == nil {
		t.Error("expected error for missing file")
	}
}
