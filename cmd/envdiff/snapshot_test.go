package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/snapshotter"
)

func writeTempEnvForSnapshot(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestSnapshotSave(t *testing.T) {
	envFile := writeTempEnvForSnapshot(t, "APP=prod\nPORT=9000\n")
	outFile := filepath.Join(t.TempDir(), "snap.json")

	cmd := newSnapshotCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"save", "--label", "test-snap", envFile, outFile})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "test-snap") {
		t.Errorf("output missing label: %s", buf.String())
	}
	data, _ := os.ReadFile(outFile)
	var snap snapshotter.Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		t.Fatalf("invalid JSON snapshot: %v", err)
	}
	if snap.Env["PORT"] != "9000" {
		t.Errorf("PORT: got %q, want %q", snap.Env["PORT"], "9000")
	}
}

func TestSnapshotDiff_Changes(t *testing.T) {
	dir := t.TempDir()
	baseSnap := filepath.Join(dir, "base.json")
	headSnap := filepath.Join(dir, "head.json")

	snapshotter.Save(baseSnap, "base", map[string]string{"A": "1", "B": "2"})
	snapshotter.Save(headSnap, "head", map[string]string{"A": "changed", "C": "3"})

	cmd := newSnapshotCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"diff", baseSnap, headSnap})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "+ Added") {
		t.Errorf("missing added section: %s", out)
	}
	if !strings.Contains(out, "- Removed") {
		t.Errorf("missing removed section: %s", out)
	}
	if !strings.Contains(out, "~ Changed") {
		t.Errorf("missing changed section: %s", out)
	}
}

func TestSnapshotDiff_NoChanges(t *testing.T) {
	dir := t.TempDir()
	baseSnap := filepath.Join(dir, "base.json")
	headSnap := filepath.Join(dir, "head.json")
	env := map[string]string{"KEY": "value"}
	snapshotter.Save(baseSnap, "base", env)
	snapshotter.Save(headSnap, "head", env)

	cmd := newSnapshotCmd()
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{"diff", baseSnap, headSnap})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected no-diff message: %s", buf.String())
	}
}
