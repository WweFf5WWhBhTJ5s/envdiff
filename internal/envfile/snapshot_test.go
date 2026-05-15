package envfile

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeSnapshotEnv(entries map[string]string) EnvFile {
	ef := EnvFile{Path: "test.env"}
	for k, v := range entries {
		ef.Entries = append(ef.Entries, Entry{Key: k, Value: v})
	}
	return ef
}

func TestTakeSnapshot_CapturesEntries(t *testing.T) {
	ef := makeSnapshotEnv(map[string]string{"APP": "prod", "PORT": "8080"})
	snap := TakeSnapshot(ef, SnapshotOptions{Label: "v1", Source: "test.env"})
	if snap.Label != "v1" {
		t.Errorf("expected label v1, got %s", snap.Label)
	}
	if snap.Entries["APP"] != "prod" {
		t.Errorf("expected APP=prod")
	}
	if snap.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestTakeSnapshot_MasksSensitive(t *testing.T) {
	ef := makeSnapshotEnv(map[string]string{"SECRET_KEY": "abc123", "APP": "prod"})
	masker := NewMasker()
	snap := TakeSnapshot(ef, SnapshotOptions{Label: "masked", Masker: masker})
	if snap.Entries["SECRET_KEY"] == "abc123" {
		t.Error("expected SECRET_KEY to be masked")
	}
	if snap.Entries["APP"] != "prod" {
		t.Error("expected APP to be unmasked")
	}
}

func TestSaveAndLoadSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")
	original := Snapshot{
		Label:     "test",
		Timestamp: time.Now().UTC().Truncate(time.Second),
		Source:    "prod.env",
		Entries:   map[string]string{"KEY": "value"},
	}
	if err := SaveSnapshot(original, path); err != nil {
		t.Fatalf("SaveSnapshot: %v", err)
	}
	loaded, err := LoadSnapshot(path)
	if err != nil {
		t.Fatalf("LoadSnapshot: %v", err)
	}
	if loaded.Label != original.Label {
		t.Errorf("label mismatch: got %s", loaded.Label)
	}
	if loaded.Entries["KEY"] != "value" {
		t.Error("entry mismatch")
	}
}

func TestLoadSnapshot_MissingFile(t *testing.T) {
	_, err := LoadSnapshot("/nonexistent/snap.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestDiffSnapshot_DetectsChanges(t *testing.T) {
	base := Snapshot{Entries: map[string]string{"A": "1", "B": "2"}}
	head := Snapshot{Entries: map[string]string{"A": "1", "B": "3", "C": "4"}}
	results := DiffSnapshot(base, head)
	statuses := map[string]DiffStatus{}
	for _, r := range results {
		statuses[r.Key] = r.Status
	}
	if statuses["A"] != StatusUnchanged {
		t.Error("expected A unchanged")
	}
	if statuses["B"] != StatusChanged {
		t.Error("expected B changed")
	}
	if statuses["C"] != StatusAdded {
		t.Error("expected C added")
	}
}

func TestSaveSnapshot_InvalidPath(t *testing.T) {
	snap := Snapshot{Label: "x", Entries: map[string]string{}}
	err := SaveSnapshot(snap, filepath.Join(os.DevNull, "bad", "path.json"))
	if err == nil {
		t.Error("expected error for invalid path")
	}
}
