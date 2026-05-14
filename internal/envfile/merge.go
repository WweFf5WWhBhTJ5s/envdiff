package envfile

import "fmt"

// MergeStrategy controls how conflicts are resolved during a merge.
type MergeStrategy int

const (
	// MergeStrategyOurs keeps the value from the base file on conflict.
	MergeStrategyOurs MergeStrategy = iota
	// MergeStrategyTheirs keeps the value from the incoming file on conflict.
	MergeStrategyTheirs
	// MergeStrategyError returns an error on any conflict.
	MergeStrategyError
)

// MergeOptions configures the merge behaviour.
type MergeOptions struct {
	Strategy        MergeStrategy
	IncludeComments bool
}

// MergeResult holds the merged EnvFile and a log of decisions made.
type MergeResult struct {
	File      EnvFile
	Conflicts []MergeConflict
}

// MergeConflict records a key where both files had differing values.
type MergeConflict struct {
	Key      string
	BaseVal  string
	TheirVal string
	Chosen   string
}

// Merge combines base and incoming EnvFiles according to opts.
// Keys present only in incoming are always added.
// Keys present in both are resolved via the chosen strategy.
func Merge(base, incoming EnvFile, opts MergeOptions) (MergeResult, error) {
	seen := make(map[string]bool, len(base.Entries))
	result := MergeResult{}

	merged := make([]EnvEntry, 0, len(base.Entries))

	for _, entry := range base.Entries {
		seen[entry.Key] = true
		incomingVal, exists := findKey(incoming, entry.Key)
		if !exists || entry.Value == incomingVal {
			merged = append(merged, entry)
			continue
		}

		// Conflict detected.
		switch opts.Strategy {
		case MergeStrategyOurs:
			result.Conflicts = append(result.Conflicts, MergeConflict{
				Key: entry.Key, BaseVal: entry.Value,
				TheirVal: incomingVal, Chosen: entry.Value,
			})
			merged = append(merged, entry)
		case MergeStrategyTheirs:
			result.Conflicts = append(result.Conflicts, MergeConflict{
				Key: entry.Key, BaseVal: entry.Value,
				TheirVal: incomingVal, Chosen: incomingVal,
			})
			merged = append(merged, EnvEntry{Key: entry.Key, Value: incomingVal, Comment: entry.Comment})
		case MergeStrategyError:
			return MergeResult{}, fmt.Errorf("merge conflict on key %q: %q vs %q", entry.Key, entry.Value, incomingVal)
		}
	}

	// Append keys that only exist in incoming.
	for _, entry := range incoming.Entries {
		if !seen[entry.Key] {
			merged = append(merged, entry)
		}
	}

	result.File = EnvFile{Entries: merged}
	return result, nil
}

func findKey(f EnvFile, key string) (string, bool) {
	for _, e := range f.Entries {
		if e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}
