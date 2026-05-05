package reconciler_test

import (
	"testing"

	"github.com/user/envdiff/internal/reconciler"
)

func TestPlan_Add(t *testing.T) {
	src := map[string]string{"A": "1"}
	dst := map[string]string{"A": "1", "B": "2"}

	steps := reconciler.Plan(src, dst)
	if len(steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(steps))
	}
	if steps[0].Action != reconciler.ActionAdd || steps[0].Key != "B" || steps[0].Value != "2" {
		t.Errorf("unexpected step: %+v", steps[0])
	}
}

func TestPlan_Remove(t *testing.T) {
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{"A": "1"}

	steps := reconciler.Plan(src, dst)
	if len(steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(steps))
	}
	if steps[0].Action != reconciler.ActionRemove || steps[0].Key != "B" {
		t.Errorf("unexpected step: %+v", steps[0])
	}
}

func TestPlan_Update(t *testing.T) {
	src := map[string]string{"A": "old"}
	dst := map[string]string{"A": "new"}

	steps := reconciler.Plan(src, dst)
	if len(steps) != 1 {
		t.Fatalf("expected 1 step, got %d", len(steps))
	}
	if steps[0].Action != reconciler.ActionUpdate || steps[0].Key != "A" || steps[0].Value != "new" {
		t.Errorf("unexpected step: %+v", steps[0])
	}
}

func TestPlan_NoChanges(t *testing.T) {
	src := map[string]string{"A": "1"}
	steps := reconciler.Plan(src, src)
	if len(steps) != 0 {
		t.Errorf("expected no steps, got %d", len(steps))
	}
}

func TestApply(t *testing.T) {
	src := map[string]string{"A": "1", "B": "old"}
	dst := map[string]string{"A": "1", "B": "new", "C": "3"}

	steps := reconciler.Plan(src, dst)
	result := reconciler.Apply(src, steps)

	if result["B"] != "new" {
		t.Errorf("expected B=new, got %s", result["B"])
	}
	if result["C"] != "3" {
		t.Errorf("expected C=3, got %s", result["C"])
	}
	if result["A"] != "1" {
		t.Errorf("expected A=1, got %s", result["A"])
	}
}

func TestFormatSteps_Empty(t *testing.T) {
	out := reconciler.FormatSteps(nil)
	if out != "No changes needed." {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormatSteps_NonEmpty(t *testing.T) {
	steps := []reconciler.Step{
		{Action: reconciler.ActionAdd, Key: "X", Value: "1"},
		{Action: reconciler.ActionRemove, Key: "Y"},
		{Action: reconciler.ActionUpdate, Key: "Z", Value: "2"},
	}
	out := reconciler.FormatSteps(steps)
	expected := "+ X=1\n- Y\n~ Z=2\n"
	if out != expected {
		t.Errorf("expected %q, got %q", expected, out)
	}
}
