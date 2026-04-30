package differ

import "sort"

// DiffResult holds the categorized differences between two env maps.
type DiffResult struct {
	Added   map[string]string // keys present in target but not in base
	Removed map[string]string // keys present in base but not in target
	Changed map[string][2]string // keys present in both but with different values [base, target]
	Unchanged map[string]string // keys present in both with the same value
}

// Diff compares two environment variable maps (base vs target) and returns a DiffResult.
func Diff(base, target map[string]string) DiffResult {
	result := DiffResult{
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Changed:   make(map[string][2]string),
		Unchanged: make(map[string]string),
	}

	for key, baseVal := range base {
		if targetVal, ok := target[key]; ok {
			if baseVal == targetVal {
				result.Unchanged[key] = baseVal
			} else {
				result.Changed[key] = [2]string{baseVal, targetVal}
			}
		} else {
			result.Removed[key] = baseVal
		}
	}

	for key, targetVal := range target {
		if _, ok := base[key]; !ok {
			result.Added[key] = targetVal
		}
	}

	return result
}

// SortedKeys returns the keys of a map in sorted order.
func SortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// HasDifferences returns true if the DiffResult contains any added, removed, or changed keys.
func (d DiffResult) HasDifferences() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}
