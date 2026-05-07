package merger

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMerge_NoConflicts(t *testing.T) {
	base := map[string]string{"A": "1", "B": "2"}
	incoming := map[string]string{"C": "3"}
	res := Merge(base, incoming, StrategyOurs)
	want := map[string]string{"A": "1", "B": "2", "C": "3"}
	if diff := cmp.Diff(want, res.Env); diff != "" {
		t.Errorf("Env mismatch (-want +got):\n%s", diff)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
}

func TestMerge_StrategyOurs(t *testing.T) {
	base := map[string]string{"A": "base"}
	incoming := map[string]string{"A": "theirs"}
	res := Merge(base, incoming, StrategyOurs)
	if res.Env["A"] != "base" {
		t.Errorf("expected 'base', got %q", res.Env["A"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0].Resolved != "base" {
		t.Errorf("unexpected conflicts: %+v", res.Conflicts)
	}
}

func TestMerge_StrategyTheirs(t *testing.T) {
	base := map[string]string{"A": "base"}
	incoming := map[string]string{"A": "theirs"}
	res := Merge(base, incoming, StrategyTheirs)
	if res.Env["A"] != "theirs" {
		t.Errorf("expected 'theirs', got %q", res.Env["A"])
	}
	if len(res.Conflicts) != 1 || res.Conflicts[0].Resolved != "theirs" {
		t.Errorf("unexpected conflicts: %+v", res.Conflicts)
	}
}

func TestMerge_StrategyUnion(t *testing.T) {
	base := map[string]string{"A": "base", "B": "shared"}
	incoming := map[string]string{"A": "theirs", "C": "new"}
	res := Merge(base, incoming, StrategyUnion)
	if res.Env["A"] != "theirs" {
		t.Errorf("A: expected 'theirs', got %q", res.Env["A"])
	}
	if res.Env["B"] != "shared" {
		t.Errorf("B: expected 'shared', got %q", res.Env["B"])
	}
	if res.Env["C"] != "new" {
		t.Errorf("C: expected 'new', got %q", res.Env["C"])
	}
}

func TestMerge_ConflictsSorted(t *testing.T) {
	base := map[string]string{"Z": "1", "A": "1", "M": "1"}
	incoming := map[string]string{"Z": "2", "A": "2", "M": "2"}
	res := Merge(base, incoming, StrategyOurs)
	keys := make([]string, len(res.Conflicts))
	for i, c := range res.Conflicts {
		keys[i] = c.Key
	}
	expected := []string{"A", "M", "Z"}
	if diff := cmp.Diff(expected, keys); diff != "" {
		t.Errorf("conflict order mismatch (-want +got):\n%s", diff)
	}
}

func TestMerge_EmptyInputs(t *testing.T) {
	res := Merge(map[string]string{}, map[string]string{}, StrategyOurs)
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %v", res.Env)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts")
	}
}
