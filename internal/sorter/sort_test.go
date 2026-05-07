package sorter_test

import (
	"testing"

	"github.com/user/envdiff/internal/sorter"
)

var sampleEnv = map[string]string{
	"ZEBRA":     "1",
	"apple":     "2",
	"MID_VAR":   "3",
	"A":         "4",
	"LONG_NAME_KEY": "5",
}

func TestSortedKeys_Alpha(t *testing.T) {
	keys := sorter.SortedKeys(sampleEnv, sorter.SortAlpha)
	if len(keys) != len(sampleEnv) {
		t.Fatalf("expected %d keys, got %d", len(sampleEnv), len(keys))
	}
	for i := 1; i < len(keys); i++ {
		if keys[i-1] > keys[i] {
			t.Errorf("keys not sorted alpha: %q > %q", keys[i-1], keys[i])
		}
	}
}

func TestSortedKeys_AlphaDesc(t *testing.T) {
	keys := sorter.SortedKeys(sampleEnv, sorter.SortAlphaDesc)
	for i := 1; i < len(keys); i++ {
		if keys[i-1] < keys[i] {
			t.Errorf("keys not sorted alpha-desc: %q < %q", keys[i-1], keys[i])
		}
	}
}

func TestSortedKeys_Length(t *testing.T) {
	keys := sorter.SortedKeys(sampleEnv, sorter.SortLength)
	for i := 1; i < len(keys); i++ {
		if len(keys[i-1]) > len(keys[i]) {
			t.Errorf("keys not sorted by length asc: len(%q)=%d > len(%q)=%d",
				keys[i-1], len(keys[i-1]), keys[i], len(keys[i]))
		}
	}
}

func TestSortedKeys_LengthDesc(t *testing.T) {
	keys := sorter.SortedKeys(sampleEnv, sorter.SortLengthDesc)
	for i := 1; i < len(keys); i++ {
		if len(keys[i-1]) < len(keys[i]) {
			t.Errorf("keys not sorted by length desc: len(%q)=%d < len(%q)=%d",
				keys[i-1], len(keys[i-1]), keys[i], len(keys[i]))
		}
	}
}

func TestSortedEnv_ContainsAllPairs(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	result := sorter.SortedEnv(env, sorter.SortAlpha)
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0] != "BAZ=qux" {
		t.Errorf("expected BAZ=qux first, got %q", result[0])
	}
	if result[1] != "FOO=bar" {
		t.Errorf("expected FOO=bar second, got %q", result[1])
	}
}

func TestSortedEnv_EmptyMap(t *testing.T) {
	result := sorter.SortedEnv(map[string]string{}, sorter.SortAlpha)
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestSortedKeys_DefaultFallback(t *testing.T) {
	keys := sorter.SortedKeys(sampleEnv, sorter.SortOrder("unknown"))
	if len(keys) != len(sampleEnv) {
		t.Fatalf("expected %d keys, got %d", len(sampleEnv), len(keys))
	}
}
