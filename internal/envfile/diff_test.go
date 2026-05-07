package envfile

import (
	"testing"
)

func makeEnvFile(path string, pairs map[string]string) *EnvFile {
	env := &EnvFile{
		Path:  path,
		Index: make(map[string]*Entry),
	}
	for k, v := range pairs {
		e := Entry{Key: k, Value: v}
		env.Entries = append(env.Entries, e)
		env.Index[k] = &env.Entries[len(env.Entries)-1]
	}
	return env
}

func TestDiff_Unchanged(t *testing.T) {
	left := makeEnvFile(".env.left", map[string]string{"KEY": "value"})
	right := makeEnvFile(".env.right", map[string]string{"KEY": "value"})

	result := Diff(left, right)
	if result.HasChanges() {
		t.Error("expected no changes")
	}
	if len(result.Entries) != 1 || result.Entries[0].Kind != Unchanged {
		t.Errorf("expected unchanged entry, got %v", result.Entries)
	}
}

func TestDiff_Changed(t *testing.T) {
	left := makeEnvFile(".env.left", map[string]string{"DB_URL": "old"})
	right := makeEnvFile(".env.right", map[string]string{"DB_URL": "new"})

	result := Diff(left, right)
	if !result.HasChanges() {
		t.Error("expected changes")
	}
	entry := result.Entries[0]
	if entry.Kind != Changed {
		t.Errorf("expected Changed, got %s", entry.Kind)
	}
	if entry.LeftValue != "old" || entry.RightValue != "new" {
		t.Errorf("unexpected values: left=%s right=%s", entry.LeftValue, entry.RightValue)
	}
}

func TestDiff_Added(t *testing.T) {
	left := makeEnvFile(".env.left", map[string]string{})
	right := makeEnvFile(".env.right", map[string]string{"NEW_KEY": "val"})

	result := Diff(left, right)
	if len(result.Entries) != 1 || result.Entries[0].Kind != Added {
		t.Errorf("expected added entry, got %v", result.Entries)
	}
}

func TestDiff_Removed(t *testing.T) {
	left := makeEnvFile(".env.left", map[string]string{"OLD_KEY": "val"})
	right := makeEnvFile(".env.right", map[string]string{})

	result := Diff(left, right)
	if len(result.Entries) != 1 || result.Entries[0].Kind != Removed {
		t.Errorf("expected removed entry, got %v", result.Entries)
	}
}

func TestDiff_Mixed(t *testing.T) {
	left := makeEnvFile(".env.left", map[string]string{
		"SHARED": "same",
		"CHANGED": "old",
		"REMOVED": "gone",
	})
	right := makeEnvFile(".env.right", map[string]string{
		"SHARED":  "same",
		"CHANGED": "new",
		"ADDED":   "here",
	})

	result := Diff(left, right)
	if !result.HasChanges() {
		t.Error("expected changes")
	}

	kinds := make(map[string]DiffKind)
	for _, e := range result.Entries {
		kinds[e.Key] = e.Kind
	}

	if kinds["SHARED"] != Unchanged {
		t.Errorf("SHARED should be unchanged")
	}
	if kinds["CHANGED"] != Changed {
		t.Errorf("CHANGED should be changed")
	}
	if kinds["REMOVED"] != Removed {
		t.Errorf("REMOVED should be removed")
	}
	if kinds["ADDED"] != Added {
		t.Errorf("ADDED should be added")
	}
}
