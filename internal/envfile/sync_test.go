package envfile

import (
	"os"
	"path/filepath"
	"testing"
)

func makeEnvFileFromMap(keys []string, vals map[string]string) *EnvFile {
	ef := &EnvFile{}
	for _, k := range keys {
		ef.Entries = append(ef.Entries, EnvEntry{Key: k, Value: vals[k]})
	}
	return ef
}

func TestSync_AddMissing(t *testing.T) {
	source := makeEnvFileFromMap([]string{"A", "B", "C"}, map[string]string{"A": "1", "B": "2", "C": "3"})
	target := makeEnvFileFromMap([]string{"A"}, map[string]string{"A": "1"})

	tmpFile := filepath.Join(t.TempDir(), ".env")
	result, err := Sync(source, target, tmpFile, SyncOptions{Mode: SyncAddMissing})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(result.Added))
	}
	if len(result.Removed) != 0 {
		t.Errorf("expected 0 removed, got %d", len(result.Removed))
	}
}

func TestSync_Full_RemovesExtra(t *testing.T) {
	source := makeEnvFileFromMap([]string{"A", "B"}, map[string]string{"A": "1", "B": "2"})
	target := makeEnvFileFromMap([]string{"A", "B", "X"}, map[string]string{"A": "1", "B": "2", "X": "99"})

	tmpFile := filepath.Join(t.TempDir(), ".env")
	result, err := Sync(source, target, tmpFile, SyncOptions{Mode: SyncFull})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Removed) != 1 || result.Removed[0] != "X" {
		t.Errorf("expected X to be removed, got %v", result.Removed)
	}
}

func TestSync_DryRun_DoesNotWrite(t *testing.T) {
	source := makeEnvFileFromMap([]string{"A", "B"}, map[string]string{"A": "1", "B": "2"})
	target := makeEnvFileFromMap([]string{"A"}, map[string]string{"A": "1"})

	tmpFile := filepath.Join(t.TempDir(), ".env")
	_, err := Sync(source, target, tmpFile, SyncOptions{Mode: SyncAddMissing, DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, statErr := os.Stat(tmpFile); !os.IsNotExist(statErr) {
		t.Error("expected file to not be created during dry run")
	}
}

func TestSync_KeepValues(t *testing.T) {
	source := makeEnvFileFromMap([]string{"A"}, map[string]string{"A": "new"})
	target := makeEnvFileFromMap([]string{"A"}, map[string]string{"A": "old"})

	tmpFile := filepath.Join(t.TempDir(), ".env")
	result, err := Sync(source, target, tmpFile, SyncOptions{Mode: SyncAddMissing, KeepValues: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Updated) != 0 {
		t.Errorf("expected 0 updated with KeepValues=true, got %d", len(result.Updated))
	}

	parsed, err := Parse(tmpFile)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if parsed.Entries[0].Value != "old" {
		t.Errorf("expected old value to be kept, got %q", parsed.Entries[0].Value)
	}
}
