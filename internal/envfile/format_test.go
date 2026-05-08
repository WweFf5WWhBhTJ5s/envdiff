package envfile

import (
	"strings"
	"testing"
)

func sampleDiffResults() []DiffResult {
	return []DiffResult{
		{Key: "APP_NAME", Status: StatusUnchanged, OldValue: "myapp", NewValue: "myapp"},
		{Key: "DB_HOST", Status: StatusChanged, OldValue: "localhost", NewValue: "prod-db"},
		{Key: "NEW_KEY", Status: StatusAdded, NewValue: "newval"},
		{Key: "OLD_KEY", Status: StatusRemoved, OldValue: "oldval"},
	}
}

func TestFormatDiff_Text(t *testing.T) {
	var sb strings.Builder
	err := FormatDiff(&sb, sampleDiffResults(), FormatText)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "+ NEW_KEY=newval") {
		t.Errorf("missing added line, got:\n%s", out)
	}
	if !strings.Contains(out, "- OLD_KEY=oldval") {
		t.Errorf("missing removed line, got:\n%s", out)
	}
	if !strings.Contains(out, "~ DB_HOST: localhost -> prod-db") {
		t.Errorf("missing changed line, got:\n%s", out)
	}
	if !strings.Contains(out, "  APP_NAME=myapp") {
		t.Errorf("missing unchanged line, got:\n%s", out)
	}
}

func TestFormatDiff_Table(t *testing.T) {
	var sb strings.Builder
	err := FormatDiff(&sb, sampleDiffResults(), FormatTable)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "STATUS") || !strings.Contains(out, "KEY") {
		t.Errorf("missing table header, got:\n%s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("missing key in table, got:\n%s", out)
	}
}

func TestFormatDiff_Dotenv(t *testing.T) {
	var sb strings.Builder
	err := FormatDiff(&sb, sampleDiffResults(), FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if strings.Contains(out, "OLD_KEY") {
		t.Errorf("removed keys should not appear in dotenv output, got:\n%s", out)
	}
	if !strings.Contains(out, "NEW_KEY=newval") {
		t.Errorf("added key missing from dotenv output, got:\n%s", out)
	}
	if !strings.Contains(out, "APP_NAME=myapp") {
		t.Errorf("unchanged key missing from dotenv output, got:\n%s", out)
	}
}

func TestFormatDiff_DefaultIsText(t *testing.T) {
	var sb strings.Builder
	err := FormatDiff(&sb, sampleDiffResults(), OutputFormat("unknown"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sb.Len() == 0 {
		t.Error("expected non-empty output for unknown format (should default to text)")
	}
}
