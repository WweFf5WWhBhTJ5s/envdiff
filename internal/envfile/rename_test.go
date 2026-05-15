package envfile

import (
	"testing"
)

func makeRenameEnv(pairs ...string) EnvFile {
	entries := make([]EnvEntry, 0, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		entries = append(entries, EnvEntry{Key: pairs[i], Value: pairs[i+1]})
	}
	return EnvFile{Path: ".env", Entries: entries}
}

func TestRename_BasicRename(t *testing.T) {
	ef := makeRenameEnv("OLD_KEY", "value1", "OTHER", "value2")
	result, results, err := Rename(ef, map[string]string{"OLD_KEY": "NEW_KEY"}, RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || !results[0].Renamed {
		t.Fatalf("expected rename to succeed, got %+v", results)
	}
	if result.Entries[0].Key != "NEW_KEY" {
		t.Errorf("expected NEW_KEY, got %q", result.Entries[0].Key)
	}
	if result.Entries[0].Value != "value1" {
		t.Errorf("expected value1 preserved, got %q", result.Entries[0].Value)
	}
}

func TestRename_MissingKey_NoError(t *testing.T) {
	ef := makeRenameEnv("EXISTING", "val")
	_, results, err := Rename(ef, map[string]string{"MISSING": "NEW"}, RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Renamed {
		t.Error("expected Renamed=false for missing key")
	}
}

func TestRename_MissingKey_ErrorOnMissing(t *testing.T) {
	ef := makeRenameEnv("EXISTING", "val")
	_, _, err := Rename(ef, map[string]string{"MISSING": "NEW"}, RenameOptions{ErrorOnMissing: true})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRename_Collision_ErrorOnCollision(t *testing.T) {
	ef := makeRenameEnv("A", "1", "B", "2")
	_, _, err := Rename(ef, map[string]string{"A": "B"}, RenameOptions{ErrorOnCollision: true})
	if err == nil {
		t.Fatal("expected error on collision")
	}
}

func TestRename_Collision_NoError(t *testing.T) {
	ef := makeRenameEnv("A", "1", "B", "2")
	_, results, err := Rename(ef, map[string]string{"A": "B"}, RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Renamed {
		t.Error("expected Renamed=false on collision without ErrorOnCollision")
	}
}

func TestRename_DryRun_DoesNotMutate(t *testing.T) {
	ef := makeRenameEnv("OLD", "val")
	result, results, err := Rename(ef, map[string]string{"OLD": "NEW"}, RenameOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !results[0].Renamed {
		t.Error("dry run result should still indicate rename would succeed")
	}
	if result.Entries[0].Key != "OLD" {
		t.Errorf("dry run must not mutate entries, got %q", result.Entries[0].Key)
	}
}

func TestRename_PreservesOrder(t *testing.T) {
	ef := makeRenameEnv("A", "1", "B", "2", "C", "3")
	result, _, err := Rename(ef, map[string]string{"B": "Z"}, RenameOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := []string{result.Entries[0].Key, result.Entries[1].Key, result.Entries[2].Key}
	expected := []string{"A", "Z", "C"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("position %d: expected %q got %q", i, expected[i], k)
		}
	}
}
