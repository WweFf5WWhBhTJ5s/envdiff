package envfile

import "fmt"

// PatchOptions controls how Patch behaves.
type PatchOptions struct {
	// ErrorOnMissing returns an error if a key to update does not exist.
	ErrorOnMissing bool
	// ErrorOnUnknown returns an error if a key in the patch is not in the file.
	ErrorOnUnknown bool
}

// PatchResult describes what happened to a single key during a patch.
type PatchResult struct {
	Key    string
	OldVal string
	NewVal string
	Action string // "set", "added", "skipped"
}

// Patch applies a map of key→value overrides to an EnvFile.
// Keys already present are updated; keys absent are added unless
// ErrorOnMissing is set. Keys in the patch that are unknown when
// ErrorOnUnknown is set cause an error.
func Patch(ef EnvFile, patches map[string]string, opts PatchOptions) (EnvFile, []PatchResult, error) {
	if opts.ErrorOnUnknown {
		for k := range patches {
			if _, found := findEntryIndex(ef, k); !found {
				return ef, nil, fmt.Errorf("patch: unknown key %q", k)
			}
		}
	}

	results := make([]PatchResult, 0, len(patches))
	out := EnvFile{Path: ef.Path, Entries: make([]EnvEntry, len(ef.Entries))}
	copy(out.Entries, ef.Entries)

	for k, newVal := range patches {
		idx, found := findEntryIndex(out, k)
		if found {
			old := out.Entries[idx].Value
			out.Entries[idx].Value = newVal
			results = append(results, PatchResult{Key: k, OldVal: old, NewVal: newVal, Action: "set"})
		} else if opts.ErrorOnMissing {
			return ef, nil, fmt.Errorf("patch: key %q not found", k)
		} else {
			out.Entries = append(out.Entries, EnvEntry{Key: k, Value: newVal})
			results = append(results, PatchResult{Key: k, OldVal: "", NewVal: newVal, Action: "added"})
		}
	}

	return out, results, nil
}

// findEntryIndex returns the index of key k in ef, or -1/false if absent.
func findEntryIndex(ef EnvFile, k string) (int, bool) {
	for i, e := range ef.Entries {
		if e.Key == k {
			return i, true
		}
	}
	return -1, false
}
