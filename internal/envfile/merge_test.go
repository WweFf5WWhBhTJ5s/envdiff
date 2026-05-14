package envfile

import (
	"testing"
)

func makeMergeEnv(pairs ...string) EnvFile {
	if len(pairs)%2 != 0 {
		panic("pairs must be even")
	}
	entries := make([]EnvEntry, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		entries = append(entries, EnvEntry{Key: pairs[i], Value: pairs[i+1]})
	}
	return EnvFile{Entries: entries}
}

func TestMerge_NoConflict(t *testing.T) {
	base := makeMergeEnv("A", "1", "B", "2")
	incoming := makeMergeEnv("A", "1", "C", "3")

	res, err := Merge(base, incoming, MergeOptions{Strategy: MergeStrategyOurs})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(res.Conflicts))
	}
	if len(res.File.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(res.File.Entries))
	}
}

func TestMerge_ConflictOurs(t *testing.T) {
	base := makeMergeEnv("A", "base")
	incoming := makeMergeEnv("A", "theirs")

	res, err := Merge(base, incoming, MergeOptions{Strategy: MergeStrategyOurs})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(res.Conflicts))
	}
	if res.File.Entries[0].Value != "base" {
		t.Errorf("expected 'base', got %q", res.File.Entries[0].Value)
	}
}

func TestMerge_ConflictTheirs(t *testing.T) {
	base := makeMergeEnv("A", "base")
	incoming := makeMergeEnv("A", "theirs")

	res, err := Merge(base, incoming, MergeOptions{Strategy: MergeStrategyTheirs})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.File.Entries[0].Value != "theirs" {
		t.Errorf("expected 'theirs', got %q", res.File.Entries[0].Value)
	}
}

func TestMerge_ConflictError(t *testing.T) {
	base := makeMergeEnv("A", "base")
	incoming := makeMergeEnv("A", "theirs")

	_, err := Merge(base, incoming, MergeOptions{Strategy: MergeStrategyError})
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestMerge_IncomingOnlyKeys(t *testing.T) {
	base := makeMergeEnv("A", "1")
	incoming := makeMergeEnv("B", "2", "C", "3")

	res, err := Merge(base, incoming, MergeOptions{Strategy: MergeStrategyOurs})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.File.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(res.File.Entries))
	}
}

func TestMerge_EmptyBase(t *testing.T) {
	base := EnvFile{}
	incoming := makeMergeEnv("X", "10")

	res, err := Merge(base, incoming, MergeOptions{Strategy: MergeStrategyOurs})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.File.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(res.File.Entries))
	}
}
