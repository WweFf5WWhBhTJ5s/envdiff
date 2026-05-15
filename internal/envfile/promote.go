package envfile

import (
	"fmt"
)

// PromoteOptions controls how promotion between environments behaves.
type PromoteOptions struct {
	// OnlyKeys restricts promotion to a specific set of keys. If empty, all keys are promoted.
	OnlyKeys []string
	// SkipKeys excludes specific keys from promotion.
	SkipKeys []string
	// OverwriteExisting controls whether existing keys in the target are overwritten.
	OverwriteExisting bool
	// DryRun returns the result without writing to the target file.
	DryRun bool
}

// PromoteResult describes the outcome of a promotion operation.
type PromoteResult struct {
	Promoted []string
	Skipped  []string
	Overwritten []string
}

// Promote copies keys from src EnvFile into dst EnvFile according to opts.
// It returns a PromoteResult describing what happened and optionally writes
// the updated dst to disk (unless DryRun is set).
func Promote(src, dst *EnvFile, opts PromoteOptions) (*EnvFile, PromoteResult, error) {
	result := PromoteResult{}

	skipSet := make(map[string]bool, len(opts.SkipKeys))
	for _, k := range opts.SkipKeys {
		skipSet[k] = true
	}

	onlySet := make(map[string]bool, len(opts.OnlyKeys))
	for _, k := range opts.OnlyKeys {
		onlySet[k] = true
	}

	// Build a mutable copy of dst entries.
	outEntries := make([]EnvEntry, len(dst.Entries))
	copy(outEntries, dst.Entries)

	dstIndex := make(map[string]int, len(outEntries))
	for i, e := range outEntries {
		dstIndex[e.Key] = i
	}

	for _, entry := range src.Entries {
		key := entry.Key

		if skipSet[key] {
			result.Skipped = append(result.Skipped, key)
			continue
		}
		if len(onlySet) > 0 && !onlySet[key] {
			result.Skipped = append(result.Skipped, key)
			continue
		}

		if idx, exists := dstIndex[key]; exists {
			if !opts.OverwriteExisting {
				result.Skipped = append(result.Skipped, key)
				continue
			}
			outEntries[idx].Value = entry.Value
			result.Overwritten = append(result.Overwritten, key)
		} else {
			outEntries = append(outEntries, entry)
			dstIndex[key] = len(outEntries) - 1
			result.Promoted = append(result.Promoted, key)
		}
	}

	out := &EnvFile{
		Path:    dst.Path,
		Entries: outEntries,
	}

	if !opts.DryRun && dst.Path != "" {
		if err := writeEnvFile(out); err != nil {
			return nil, result, fmt.Errorf("promote: write %s: %w", dst.Path, err)
		}
	}

	return out, result, nil
}
