package reconciler

import (
	"fmt"
	"sort"
	"strings"

	"github.com/user/envdiff/internal/differ"
)

// Action represents the type of reconciliation action.
type Action string

const (
	ActionAdd    Action = "ADD"
	ActionRemove Action = "REMOVE"
	ActionUpdate Action = "UPDATE"
)

// Step describes a single reconciliation step.
type Step struct {
	Action Action
	Key    string
	Value  string // target value for ADD/UPDATE; empty for REMOVE
}

// Plan returns an ordered list of steps to transform src into dst.
// It uses the Diff result to determine what needs to change.
func Plan(src, dst map[string]string) []Step {
	diff := differ.Diff(src, dst)
	var steps []Step

	keys := make([]string, 0, len(diff))
	for k := range diff {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		e := diff[k]
		switch {
		case e.Added:
			steps = append(steps, Step{Action: ActionAdd, Key: k, Value: e.DstVal})
		case e.Removed:
			steps = append(steps, Step{Action: ActionRemove, Key: k})
		case e.Changed:
			steps = append(steps, Step{Action: ActionUpdate, Key: k, Value: e.DstVal})
		}
	}

	return steps
}

// Apply applies the reconciliation steps to src and returns the resulting map.
func Apply(src map[string]string, steps []Step) map[string]string {
	result := make(map[string]string, len(src))
	for k, v := range src {
		result[k] = v
	}

	for _, s := range steps {
		switch s.Action {
		case ActionAdd, ActionUpdate:
			result[s.Key] = s.Value
		case ActionRemove:
			delete(result, s.Key)
		}
	}

	return result
}

// FormatSteps returns a human-readable summary of the reconciliation plan.
func FormatSteps(steps []Step) string {
	if len(steps) == 0 {
		return "No changes needed."
	}
	var sb strings.Builder
	for _, s := range steps {
		switch s.Action {
		case ActionAdd:
			fmt.Fprintf(&sb, "+ %s=%s\n", s.Key, s.Value)
		case ActionRemove:
			fmt.Fprintf(&sb, "- %s\n", s.Key)
		case ActionUpdate:
			fmt.Fprintf(&sb, "~ %s=%s\n", s.Key, s.Value)
		}
	}
	return sb.String()
}
