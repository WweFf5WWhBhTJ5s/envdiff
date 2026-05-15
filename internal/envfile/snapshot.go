package envfile

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of an EnvFile.
type Snapshot struct {
	Label     string            `json:"label"`
	Timestamp time.Time         `json:"timestamp"`
	Source    string            `json:"source"`
	Entries   map[string]string `json:"entries"`
}

// SnapshotOptions controls snapshot behaviour.
type SnapshotOptions struct {
	Label  string
	Source string
	Masker *Masker
}

// TakeSnapshot captures the current state of an EnvFile into a Snapshot.
func TakeSnapshot(ef EnvFile, opts SnapshotOptions) Snapshot {
	entries := make(map[string]string, len(ef.Entries))
	for _, e := range ef.Entries {
		val := e.Value
		if opts.Masker != nil && opts.Masker.IsSensitive(e.Key) {
			val = opts.Masker.MaskValue(e.Key, val)
		}
		entries[e.Key] = val
	}
	return Snapshot{
		Label:     opts.Label,
		Timestamp: time.Now().UTC(),
		Source:    opts.Source,
		Entries:   entries,
	}
}

// SaveSnapshot writes a Snapshot to a JSON file at path.
func SaveSnapshot(snap Snapshot, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(snap); err != nil {
		return fmt.Errorf("snapshot: encode: %w", err)
	}
	return nil
}

// LoadSnapshot reads a Snapshot from a JSON file at path.
func LoadSnapshot(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: open file: %w", err)
	}
	defer f.Close()
	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return Snapshot{}, fmt.Errorf("snapshot: decode: %w", err)
	}
	return snap, nil
}

// DiffSnapshot compares two snapshots and returns DiffResults.
func DiffSnapshot(base, head Snapshot) []DiffResult {
	baseEnv := snapshotToEnvFile(base)
	headEnv := snapshotToEnvFile(head)
	return Diff(baseEnv, headEnv)
}

func snapshotToEnvFile(s Snapshot) EnvFile {
	ef := EnvFile{Path: s.Source}
	for k, v := range s.Entries {
		ef.Entries = append(ef.Entries, Entry{Key: k, Value: v})
	}
	return ef
}
