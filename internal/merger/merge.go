package merger

import "sort"

// Strategy defines how conflicts are resolved during a merge.
type Strategy int

const (
	// StrategyOurs keeps the value from the base map on conflict.
	StrategyOurs Strategy = iota
	// StrategyTheirs keeps the value from the incoming map on conflict.
	StrategyTheirs
	// StrategyUnion includes all keys from both maps; conflicts use StrategyTheirs.
	StrategyUnion
)

// Conflict records a key whose value differed between base and incoming.
type Conflict struct {
	Key      string
	BaseVal  string
	TheirVal string
	Resolved string
}

// Result holds the merged environment and any conflicts that were encountered.
type Result struct {
	Env       map[string]string
	Conflicts []Conflict
}

// Merge combines base and incoming according to the given strategy.
// Keys present only in base or only in incoming are always included.
// Conflicts arise when a key exists in both maps with different values.
func Merge(base, incoming map[string]string, strategy Strategy) Result {
	merged := make(map[string]string, len(base))
	for k, v := range base {
		merged[k] = v
	}

	var conflicts []Conflict

	for k, theirVal := range incoming {
		baseVal, exists := merged[k]
		if !exists {
			merged[k] = theirVal
			continue
		}
		if baseVal == theirVal {
			continue
		}
		// Conflict: same key, different values.
		resolved := baseVal
		if strategy == StrategyTheirs || strategy == StrategyUnion {
			resolved = theirVal
			merged[k] = theirVal
		}
		conflicts = append(conflicts, Conflict{
			Key:      k,
			BaseVal:  baseVal,
			TheirVal: theirVal,
			Resolved: resolved,
		})
	}

	// Sort conflicts for deterministic output.
	sort.Slice(conflicts, func(i, j int) bool {
		return conflicts[i].Key < conflicts[j].Key
	})

	return Result{Env: merged, Conflicts: conflicts}
}
