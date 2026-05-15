package envfile

import (
	"testing"
)

func makeAuditDiff() []DiffResult {
	return []DiffResult{
		{Key: "APP_NAME", Status: StatusUnchanged, OldValue: "myapp", NewValue: "myapp"},
		{Key: "DB_HOST", Status: StatusChanged, OldValue: "localhost", NewValue: "prod.db"},
		{Key: "NEW_KEY", Status: StatusAdded, OldValue: "", NewValue: "hello"},
		{Key: "OLD_KEY", Status: StatusRemoved, OldValue: "bye", NewValue: ""},
		{Key: "DB_PASSWORD", Status: StatusChanged, OldValue: "secret", NewValue: "newsecret"},
	}
}

func TestAudit_SkipsUnchanged(t *testing.T) {
	results := makeAuditDiff()
	log := Audit(results, "test", nil)
	for _, e := range log.Entries {
		if e.Key == "APP_NAME" {
			t.Error("unchanged key should not appear in audit log")
		}
	}
}

func TestAudit_RecordsActions(t *testing.T) {
	results := makeAuditDiff()
	log := Audit(results, "test", nil)

	expected := map[string]AuditAction{
		"DB_HOST":  AuditChanged,
		"NEW_KEY":  AuditAdded,
		"OLD_KEY":  AuditRemoved,
	}

	for _, e := range log.Entries {
		if want, ok := expected[e.Key]; ok {
			if e.Action != want {
				t.Errorf("key %s: got action %s, want %s", e.Key, e.Action, want)
			}
		}
	}
}

func TestAudit_MasksSecrets(t *testing.T) {
	results := makeAuditDiff()
	masker := NewMasker()
	log := Audit(results, "test", masker)

	for _, e := range log.Entries {
		if e.Key == "DB_PASSWORD" {
			if e.OldValue == "secret" || e.NewValue == "newsecret" {
				t.Error("sensitive values should be masked in audit log")
			}
			return
		}
	}
	t.Error("DB_PASSWORD entry not found in audit log")
}

func TestAudit_SourceLabel(t *testing.T) {
	results := makeAuditDiff()
	log := Audit(results, "staging->production", nil)
	for _, e := range log.Entries {
		if e.Source != "staging->production" {
			t.Errorf("expected source 'staging->production', got %q", e.Source)
		}
	}
}

func TestAudit_TimestampSet(t *testing.T) {
	results := makeAuditDiff()
	log := Audit(results, "test", nil)
	for _, e := range log.Entries {
		if e.Timestamp.IsZero() {
			t.Errorf("entry for %s has zero timestamp", e.Key)
		}
	}
}
