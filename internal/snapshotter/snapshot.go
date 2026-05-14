package snapshotter

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a captured state of an environment file at a point in time.
type Snapshot struct {
	Label     string            `json:"label"`
	Timestamp time.Time         `json:"timestamp"`
	Env       map[string]string `json:"env"`
}

// Save writes a snapshot of env to the given file path as JSON.
func Save(path, label string, env map[string]string) error {
	snap := Snapshot{
		Label:     label,
		Timestamp: time.Now().UTC(),
		Env:       env,
	}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshotter: marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("snapshotter: write %q: %w", path, err)
	}
	return nil
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshotter: read %q: %w", path, err)
	}
	var snap Snapshot
	if err := json.Unmarshal(data, &snap); err != nil {
		return nil, fmt.Errorf("snapshotter: unmarshal: %w", err)
	}
	return &snap, nil
}

// Compare returns the keys added, removed, or changed between two snapshots.
func Compare(base, head *Snapshot) (added, removed, changed []string) {
	for k, hv := range head.Env {
		if bv, ok := base.Env[k]; !ok {
			added = append(added, k)
		} else if bv != hv {
			changed = append(changed, k)
		}
	}
	for k := range base.Env {
		if _, ok := head.Env[k]; !ok {
			removed = append(removed, k)
		}
	}
	return
}
