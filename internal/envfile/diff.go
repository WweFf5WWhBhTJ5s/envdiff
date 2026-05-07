package envfile

// DiffKind represents the type of difference between two env files.
type DiffKind string

const (
	Added    DiffKind = "added"    // key exists in right but not left
	Removed  DiffKind = "removed"  // key exists in left but not right
	Changed  DiffKind = "changed"  // key exists in both but values differ
	Unchanged DiffKind = "unchanged" // key exists in both with same value
)

// DiffEntry represents a single difference between two env files.
type DiffEntry struct {
	Key       string
	Kind      DiffKind
	LeftValue string
	RightValue string
}

// Result holds the complete diff between two env files.
type Result struct {
	Left    string
	Right   string
	Entries []DiffEntry
}

// Diff computes the difference between two EnvFile instances.
func Diff(left, right *EnvFile) *Result {
	result := &Result{
		Left:  left.Path,
		Right: right.Path,
	}

	seen := make(map[string]bool)

	for _, entry := range left.Entries {
		seen[entry.Key] = true
		if re, ok := right.Index[entry.Key]; ok {
			if entry.Value == re.Value {
				result.Entries = append(result.Entries, DiffEntry{
					Key:        entry.Key,
					Kind:       Unchanged,
					LeftValue:  entry.Value,
					RightValue: re.Value,
				})
			} else {
				result.Entries = append(result.Entries, DiffEntry{
					Key:        entry.Key,
					Kind:       Changed,
					LeftValue:  entry.Value,
					RightValue: re.Value,
				})
			}
		} else {
			result.Entries = append(result.Entries, DiffEntry{
				Key:       entry.Key,
				Kind:      Removed,
				LeftValue: entry.Value,
			})
		}
	}

	for _, entry := range right.Entries {
		if !seen[entry.Key] {
			result.Entries = append(result.Entries, DiffEntry{
				Key:        entry.Key,
				Kind:       Added,
				RightValue: entry.Value,
			})
		}
	}

	return result
}

// HasChanges returns true if the diff contains any added, removed, or changed entries.
func (r *Result) HasChanges() bool {
	for _, e := range r.Entries {
		if e.Kind != Unchanged {
			return true
		}
	}
	return false
}
