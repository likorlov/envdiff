package differ

import (
	"testing"
)

func TestDiff_Added(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar", "BAZ": "qux"}

	result := Diff(base, target)

	if len(result.Added) != 1 {
		t.Fatalf("expected 1 added key, got %d", len(result.Added))
	}
	if result.Added["BAZ"] != "qux" {
		t.Errorf("expected Added[BAZ]=qux, got %q", result.Added["BAZ"])
	}
}

func TestDiff_Removed(t *testing.T) {
	base := map[string]string{"FOO": "bar", "OLD": "val"}
	target := map[string]string{"FOO": "bar"}

	result := Diff(base, target)

	if len(result.Removed) != 1 {
		t.Fatalf("expected 1 removed key, got %d", len(result.Removed))
	}
	if result.Removed["OLD"] != "val" {
		t.Errorf("expected Removed[OLD]=val, got %q", result.Removed["OLD"])
	}
}

func TestDiff_Changed(t *testing.T) {
	base := map[string]string{"FOO": "old"}
	target := map[string]string{"FOO": "new"}

	result := Diff(base, target)

	if len(result.Changed) != 1 {
		t.Fatalf("expected 1 changed key, got %d", len(result.Changed))
	}
	pair := result.Changed["FOO"]
	if pair[0] != "old" || pair[1] != "new" {
		t.Errorf("expected Changed[FOO]=[old, new], got %v", pair)
	}
}

func TestDiff_Unchanged(t *testing.T) {
	base := map[string]string{"FOO": "bar"}
	target := map[string]string{"FOO": "bar"}

	result := Diff(base, target)

	if len(result.Unchanged) != 1 {
		t.Fatalf("expected 1 unchanged key, got %d", len(result.Unchanged))
	}
}

func TestDiff_HasDifferences(t *testing.T) {
	base := map[string]string{"A": "1"}
	target := map[string]string{"A": "1"}
	result := Diff(base, target)
	if result.HasDifferences() {
		t.Error("expected no differences for identical maps")
	}

	target["B"] = "2"
	result = Diff(base, target)
	if !result.HasDifferences() {
		t.Error("expected differences when target has extra key")
	}
}

func TestDiff_EmptyMaps(t *testing.T) {
	result := Diff(map[string]string{}, map[string]string{})
	if result.HasDifferences() {
		t.Error("expected no differences for two empty maps")
	}
}
