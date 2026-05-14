package snapshotter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/snapshotter"
)

func TestSaveAndLoad(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	tmp := filepath.Join(t.TempDir(), "snap.json")

	if err := snapshotter.Save(tmp, "v1", env); err != nil {
		t.Fatalf("Save: %v", err)
	}
	snap, err := snapshotter.Load(tmp)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if snap.Label != "v1" {
		t.Errorf("label: got %q, want %q", snap.Label, "v1")
	}
	if snap.Env["PORT"] != "8080" {
		t.Errorf("PORT: got %q, want %q", snap.Env["PORT"], "8080")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := snapshotter.Load("/nonexistent/snap.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "bad.json")
	os.WriteFile(tmp, []byte("not json"), 0o644)
	_, err := snapshotter.Load(tmp)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestCompare_AddedRemovedChanged(t *testing.T) {
	base := &snapshotter.Snapshot{Env: map[string]string{
		"A": "1", "B": "2", "C": "3",
	}}
	head := &snapshotter.Snapshot{Env: map[string]string{
		"A": "1", "B": "changed", "D": "4",
	}}
	added, removed, changed := snapshotter.Compare(base, head)

	if len(added) != 1 || added[0] != "D" {
		t.Errorf("added: got %v, want [D]", added)
	}
	if len(removed) != 1 || removed[0] != "C" {
		t.Errorf("removed: got %v, want [C]", removed)
	}
	if len(changed) != 1 || changed[0] != "B" {
		t.Errorf("changed: got %v, want [B]", changed)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	env := map[string]string{"X": "1"}
	base := &snapshotter.Snapshot{Env: env}
	head := &snapshotter.Snapshot{Env: map[string]string{"X": "1"}}
	added, removed, changed := snapshotter.Compare(base, head)
	if len(added)+len(removed)+len(changed) != 0 {
		t.Errorf("expected no changes, got added=%v removed=%v changed=%v", added, removed, changed)
	}
}
