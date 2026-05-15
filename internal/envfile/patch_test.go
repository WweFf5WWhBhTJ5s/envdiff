package envfile

import (
	"testing"
)

func makePatchEnv(pairs ...string) EnvFile {
	ef := EnvFile{Path: "test.env"}
	for i := 0; i+1 < len(pairs); i += 2 {
		ef.Entries = append(ef.Entries, EnvEntry{Key: pairs[i], Value: pairs[i+1]})
	}
	return ef
}

func TestPatch_UpdateExistingKey(t *testing.T) {
	ef := makePatchEnv("HOST", "localhost", "PORT", "8080")
	out, results, err := Patch(ef, map[string]string{"PORT": "9090"}, PatchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := entryValue(out, "PORT"); v != "9090" {
		t.Errorf("expected PORT=9090, got %q", v)
	}
	if len(results) != 1 || results[0].Action != "set" {
		t.Errorf("expected one 'set' result, got %+v", results)
	}
}

func TestPatch_AddsNewKey(t *testing.T) {
	ef := makePatchEnv("HOST", "localhost")
	out, results, err := Patch(ef, map[string]string{"DEBUG": "true"}, PatchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := entryValue(out, "DEBUG"); v != "true" {
		t.Errorf("expected DEBUG=true, got %q", v)
	}
	if len(results) != 1 || results[0].Action != "added" {
		t.Errorf("expected one 'added' result, got %+v", results)
	}
}

func TestPatch_ErrorOnMissing(t *testing.T) {
	ef := makePatchEnv("HOST", "localhost")
	_, _, err := Patch(ef, map[string]string{"MISSING": "val"}, PatchOptions{ErrorOnMissing: true})
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestPatch_ErrorOnUnknown(t *testing.T) {
	ef := makePatchEnv("HOST", "localhost")
	_, _, err := Patch(ef, map[string]string{"UNKNOWN": "val"}, PatchOptions{ErrorOnUnknown: true})
	if err == nil {
		t.Fatal("expected error for unknown key, got nil")
	}
}

func TestPatch_NoMutationOfOriginal(t *testing.T) {
	ef := makePatchEnv("HOST", "localhost")
	_, _, err := Patch(ef, map[string]string{"HOST": "remotehost"}, PatchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entryValue(ef, "HOST") != "localhost" {
		t.Error("original EnvFile was mutated")
	}
}

func TestPatch_MultipleKeys(t *testing.T) {
	ef := makePatchEnv("A", "1", "B", "2", "C", "3")
	out, results, err := Patch(ef, map[string]string{"A": "10", "B": "20"}, PatchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entryValue(out, "A") != "10" || entryValue(out, "B") != "20" {
		t.Errorf("values not updated correctly")
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

func entryValue(ef EnvFile, key string) string {
	for _, e := range ef.Entries {
		if e.Key == key {
			return e.Value
		}
	}
	return ""
}
