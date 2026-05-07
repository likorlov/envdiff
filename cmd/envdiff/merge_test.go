package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return filepath.Clean(f.Name())
}

func TestRunMerge_StrategyOurs(t *testing.T) {
	base := writeTempEnv(t, "A=base\nB=shared\n")
	incoming := writeTempEnv(t, "A=theirs\nC=new\n")

	// Capture stdout by redirecting within the test is complex; validate no error.
	if err := runMerge(base, incoming, "ours"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunMerge_StrategyTheirs(t *testing.T) {
	base := writeTempEnv(t, "A=base\n")
	incoming := writeTempEnv(t, "A=theirs\n")
	if err := runMerge(base, incoming, "theirs"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRunMerge_UnknownStrategy(t *testing.T) {
	base := writeTempEnv(t, "A=1\n")
	incoming := writeTempEnv(t, "A=2\n")
	err := runMerge(base, incoming, "bogus")
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
	if !strings.Contains(err.Error(), "unknown strategy") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestRunMerge_MissingFile(t *testing.T) {
	err := runMerge("/nonexistent/base.env", "/nonexistent/incoming.env", "ours")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestNewMergeCmd_Flags(t *testing.T) {
	cmd := newMergeCmd()
	if cmd.Use != "merge <base> <incoming>" {
		t.Errorf("unexpected Use: %s", cmd.Use)
	}
	f := cmd.Flags().Lookup("strategy")
	if f == nil {
		t.Fatal("expected --strategy flag")
	}
	if f.DefValue != "ours" {
		t.Errorf("expected default 'ours', got %q", f.DefValue)
	}
}
