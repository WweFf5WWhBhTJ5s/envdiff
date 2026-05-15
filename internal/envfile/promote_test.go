package envfile

import (
	"testing"
)

func makePromoteEnv(pairs ...string) *EnvFile {
	ef := &EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		ef.Entries = append(ef.Entries, EnvEntry{Key: pairs[i], Value: pairs[i+1]})
	}
	return ef
}

func TestPromote_AddsNewKeys(t *testing.T) {
	src := makePromoteEnv("FOO", "bar", "BAZ", "qux")
	dst := makePromoteEnv("EXISTING", "val")

	out, res, err := Promote(src, dst, PromoteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d", len(res.Promoted))
	}
	if findKey(out, "FOO") == nil {
		t.Error("expected FOO in output")
	}
	if findKey(out, "EXISTING") == nil {
		t.Error("expected EXISTING preserved in output")
	}
}

func TestPromote_DoesNotOverwriteByDefault(t *testing.T) {
	src := makePromoteEnv("FOO", "new_value")
	dst := makePromoteEnv("FOO", "old_value")

	out, res, err := Promote(src, dst, PromoteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "FOO" {
		t.Errorf("expected FOO skipped, got %v", res.Skipped)
	}
	if e := findKey(out, "FOO"); e == nil || e.Value != "old_value" {
		t.Error("expected FOO to retain old_value")
	}
}

func TestPromote_OverwriteExisting(t *testing.T) {
	src := makePromoteEnv("FOO", "new_value")
	dst := makePromoteEnv("FOO", "old_value")

	out, res, err := Promote(src, dst, PromoteOptions{DryRun: true, OverwriteExisting: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwritten) != 1 || res.Overwritten[0] != "FOO" {
		t.Errorf("expected FOO overwritten, got %v", res.Overwritten)
	}
	if e := findKey(out, "FOO"); e == nil || e.Value != "new_value" {
		t.Error("expected FOO to have new_value")
	}
}

func TestPromote_OnlyKeys(t *testing.T) {
	src := makePromoteEnv("FOO", "1", "BAR", "2", "BAZ", "3")
	dst := makePromoteEnv()

	_, res, err := Promote(src, dst, PromoteOptions{DryRun: true, OnlyKeys: []string{"FOO", "BAZ"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Errorf("expected 2 promoted, got %d: %v", len(res.Promoted), res.Promoted)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "BAR" {
		t.Errorf("expected BAR skipped, got %v", res.Skipped)
	}
}

func TestPromote_SkipKeys(t *testing.T) {
	src := makePromoteEnv("SECRET", "s3cr3t", "PUBLIC", "hello")
	dst := makePromoteEnv()

	_, res, err := Promote(src, dst, PromoteOptions{DryRun: true, SkipKeys: []string{"SECRET"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 1 || res.Promoted[0] != "PUBLIC" {
		t.Errorf("expected PUBLIC promoted, got %v", res.Promoted)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "SECRET" {
		t.Errorf("expected SECRET skipped, got %v", res.Skipped)
	}
}
