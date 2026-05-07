package sorter

import (
	"sort"
	"strings"
)

// SortOrder defines the ordering strategy for env var keys.
type SortOrder string

const (
	SortAlpha      SortOrder = "alpha"
	SortAlphaDesc  SortOrder = "alpha-desc"
	SortLength     SortOrder = "length"
	SortLengthDesc SortOrder = "length-desc"
)

// SortedEnv returns a slice of "KEY=VALUE" strings sorted by the given order.
func SortedEnv(env map[string]string, order SortOrder) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	switch order {
	case SortAlphaDesc:
		sort.Slice(keys, func(i, j int) bool {
			return strings.ToLower(keys[i]) > strings.ToLower(keys[j])
		})
	case SortLength:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) == len(keys[j]) {
				return keys[i] < keys[j]
			}
			return len(keys[i]) < len(keys[j])
		})
	case SortLengthDesc:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) == len(keys[j]) {
				return keys[i] > keys[j]
			}
			return len(keys[i]) > len(keys[j])
		})
	default: // SortAlpha
		sort.Slice(keys, func(i, j int) bool {
			return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
		})
	}

	result := make([]string, 0, len(keys))
	for _, k := range keys {
		result = append(result, k+"="+env[k])
	}
	return result
}

// SortedKeys returns keys from an env map sorted by the given order.
func SortedKeys(env map[string]string, order SortOrder) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}

	switch order {
	case SortAlphaDesc:
		sort.Slice(keys, func(i, j int) bool {
			return strings.ToLower(keys[i]) > strings.ToLower(keys[j])
		})
	case SortLength:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) == len(keys[j]) {
				return keys[i] < keys[j]
			}
			return len(keys[i]) < len(keys[j])
		})
	case SortLengthDesc:
		sort.Slice(keys, func(i, j int) bool {
			if len(keys[i]) == len(keys[j]) {
				return keys[i] > keys[j]
			}
			return len(keys[i]) > len(keys[j])
		})
	default:
		sort.Slice(keys, func(i, j int) bool {
			return strings.ToLower(keys[i]) < strings.ToLower(keys[j])
		})
	}
	return keys
}
