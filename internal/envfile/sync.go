package envfile

import (
	"fmt"
	"os"
	"strings"
)

// SyncMode controls how missing keys are handled during sync.
type SyncMode int

const (
	// SyncAddMissing adds keys present in source but missing in target.
	SyncAddMissing SyncMode = iota
	// SyncFull adds missing keys and removes extra keys from target.
	SyncFull
)

// SyncOptions configures the sync operation.
type SyncOptions struct {
	Mode        SyncMode
	DryRun      bool
	KeepValues  bool // if true, keep target values for changed keys
}

// SyncResult describes what was changed during a sync.
type SyncResult struct {
	Added   []string
	Removed []string
	Updated []string
}

// Sync applies changes from source EnvFile to target EnvFile based on options.
// If DryRun is true, no file is written; only the SyncResult is returned.
func Sync(source, target *EnvFile, targetPath string, opts SyncOptions) (*SyncResult, error) {
	result := &SyncResult{}

	// Build a mutable copy of target entries.
	entries := make(map[string]string, len(target.Entries))
	for _, e := range target.Entries {
		entries[e.Key] = e.Value
	}

	// Add or update keys from source.
	for _, se := range source.Entries {
		if _, exists := entries[se.Key]; !exists {
			entries[se.Key] = se.Value
			result.Added = append(result.Added, se.Key)
		} else if !opts.KeepValues && entries[se.Key] != se.Value {
			entries[se.Key] = se.Value
			result.Updated = append(result.Updated, se.Key)
		}
	}

	// Remove keys not in source if SyncFull.
	if opts.Mode == SyncFull {
		sourceKeys := make(map[string]struct{}, len(source.Entries))
		for _, se := range source.Entries {
			sourceKeys[se.Key] = struct{}{}
		}
		for _, te := range target.Entries {
			if _, exists := sourceKeys[te.Key]; !exists {
				delete(entries, te.Key)
				result.Removed = append(result.Removed, te.Key)
			}
		}
	}

	if opts.DryRun {
		return result, nil
	}

	return result, writeEnvFile(targetPath, source, entries)
}

// writeEnvFile writes the merged entries to disk, preserving source key order.
func writeEnvFile(path string, source *EnvFile, entries map[string]string) error {
	var sb strings.Builder
	seen := make(map[string]struct{})

	for _, se := range source.Entries {
		val, ok := entries[se.Key]
		if !ok {
			continue
		}
		sb.WriteString(fmt.Sprintf("%s=%s\n", se.Key, val))
		seen[se.Key] = struct{}{}
	}

	// Append any keys that exist in entries but not in source order.
	for k, v := range entries {
		if _, ok := seen[k]; !ok {
			sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
		}
	}

	return os.WriteFile(path, []byte(sb.String()), 0o644)
}
