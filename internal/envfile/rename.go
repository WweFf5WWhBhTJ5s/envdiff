package envfile

import (
	"fmt"
	"strings"
)

// RenameResult describes the outcome of a single key rename operation.
type RenameResult struct {
	OldKey string
	NewKey string
	Renamed bool
	Reason  string
}

// RenameOptions controls the behavior of the Rename function.
type RenameOptions struct {
	// DryRun prevents any mutation; results still describe what would happen.
	DryRun bool
	// ErrorOnMissing causes Rename to return an error if OldKey is not found.
	ErrorOnMissing bool
	// ErrorOnCollision causes Rename to return an error if NewKey already exists.
	ErrorOnCollision bool
}

// Rename renames one or more keys in an EnvFile according to the provided
// mapping (oldKey -> newKey). The original order of entries is preserved.
func Rename(ef EnvFile, mapping map[string]string, opts RenameOptions) (EnvFile, []RenameResult, error) {
	results := make([]RenameResult, 0, len(mapping))

	// Build a quick lookup of existing keys.
	existing := make(map[string]bool, len(ef.Entries))
	for _, e := range ef.Entries {
		existing[e.Key] = true
	}

	for oldKey, newKey := range mapping {
		newKey = strings.TrimSpace(newKey)
		if newKey == "" {
			return ef, results, fmt.Errorf("rename: new key for %q must not be empty", oldKey)
		}

		if !existing[oldKey] {
			if opts.ErrorOnMissing {
				return ef, results, fmt.Errorf("rename: key %q not found", oldKey)
			}
			results = append(results, RenameResult{OldKey: oldKey, NewKey: newKey, Renamed: false, Reason: "key not found"})
			continue
		}

		if existing[newKey] && oldKey != newKey {
			if opts.ErrorOnCollision {
				return ef, results, fmt.Errorf("rename: key %q already exists", newKey)
			}
			results = append(results, RenameResult{OldKey: oldKey, NewKey: newKey, Renamed: false, Reason: "target key already exists"})
			continue
		}

		results = append(results, RenameResult{OldKey: oldKey, NewKey: newKey, Renamed: true})
	}

	if opts.DryRun {
		return ef, results, nil
	}

	// Apply renames that succeeded.
	renamed := make(map[string]string)
	for _, r := range results {
		if r.Renamed {
			renamed[r.OldKey] = r.NewKey
		}
	}

	updated := make([]EnvEntry, len(ef.Entries))
	for i, e := range ef.Entries {
		if newKey, ok := renamed[e.Key]; ok {
			e.Key = newKey
		}
		updated[i] = e
	}

	return EnvFile{Path: ef.Path, Entries: updated}, results, nil
}
