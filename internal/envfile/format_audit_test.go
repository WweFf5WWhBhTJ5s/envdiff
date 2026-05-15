package envfile

import (
	"strings"
	"testing"
	"time"
)

func sampleAuditLog() *AuditLog {
	ts := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	return &AuditLog{
		Entries: []AuditEntry{
			{Timestamp: ts, Key: "DB_HOST", Action: AuditChanged, OldValue: "localhost", NewValue: "prod.db", Source: "ci"},
			{Timestamp: ts, Key: "NEW_KEY", Action: AuditAdded, OldValue: "", NewValue: "hello", Source: "ci"},
		},
	}
}

func TestFormatAudit_Text(t *testing.T) {
	out := FormatAuditLog(sampleAuditLog(), "text")
	if !strings.Contains(out, "DB_HOST") {
		t.Error("text output missing DB_HOST")
	}
	if !strings.Contains(out, "changed") {
		t.Error("text output missing action 'changed'")
	}
	if !strings.Contains(out, "source: ci") {
		t.Error("text output missing source label")
	}
}

func TestFormatAudit_Table(t *testing.T) {
	out := FormatAuditLog(sampleAuditLog(), "table")
	if !strings.Contains(out, "TIMESTAMP") {
		t.Error("table output missing header")
	}
	if !strings.Contains(out, "NEW_KEY") {
		t.Error("table output missing NEW_KEY")
	}
}

func TestFormatAudit_JSON(t *testing.T) {
	out := FormatAuditLog(sampleAuditLog(), "json")
	if !strings.Contains(out, `"action"`) {
		t.Error("json output missing 'action' field")
	}
	if !strings.Contains(out, `"DB_HOST"`) {
		t.Error("json output missing DB_HOST key")
	}
	if !strings.Contains(out, `"source"`) {
		t.Error("json output missing 'source' field")
	}
}

func TestFormatAudit_DefaultIsText(t *testing.T) {
	out1 := FormatAuditLog(sampleAuditLog(), "")
	out2 := FormatAuditLog(sampleAuditLog(), "text")
	if out1 != out2 {
		t.Error("default format should equal text format")
	}
}

func TestFormatAudit_EmptyLog(t *testing.T) {
	empty := &AuditLog{}
	for _, fmt := range []string{"text", "table", "json"} {
		out := FormatAuditLog(empty, fmt)
		if out == "" {
			t.Errorf("format %s returned empty string for empty log", fmt)
		}
	}
}
