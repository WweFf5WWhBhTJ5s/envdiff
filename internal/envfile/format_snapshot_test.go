package envfile

import (
	"strings"
	"testing"
	"time"
)

func sampleSnapshot() Snapshot {
	return Snapshot{
		Label:     "release-1.2",
		Source:    "prod.env",
		Timestamp: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
		Entries: map[string]string{
			"APP_ENV": "production",
			"PORT":    "8080",
		},
	}
}

func TestFormatSnapshot_Text(t *testing.T) {
	out, err := FormatSnapshot(sampleSnapshot(), "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "release-1.2") {
		t.Error("expected label in text output")
	}
	if !strings.Contains(out, "APP_ENV=production") {
		t.Error("expected APP_ENV entry")
	}
	if !strings.Contains(out, "PORT=8080") {
		t.Error("expected PORT entry")
	}
}

func TestFormatSnapshot_Table(t *testing.T) {
	out, err := FormatSnapshot(sampleSnapshot(), "table")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "KEY") {
		t.Error("expected KEY header")
	}
	if !strings.Contains(out, "APP_ENV") {
		t.Error("expected APP_ENV row")
	}
}

func TestFormatSnapshot_JSON(t *testing.T) {
	out, err := FormatSnapshot(sampleSnapshot(), "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `"label"`) {
		t.Error("expected label field in JSON")
	}
	if !strings.Contains(out, `"release-1.2"`) {
		t.Error("expected label value in JSON")
	}
}

func TestFormatSnapshot_DefaultIsText(t *testing.T) {
	out1, _ := FormatSnapshot(sampleSnapshot(), "")
	out2, _ := FormatSnapshot(sampleSnapshot(), "text")
	if out1 != out2 {
		t.Error("default format should equal text format")
	}
}

func TestFormatSnapshot_KeysSorted(t *testing.T) {
	out, _ := FormatSnapshot(sampleSnapshot(), "text")
	idxApp := strings.Index(out, "APP_ENV")
	idxPort := strings.Index(out, "PORT")
	if idxApp > idxPort {
		t.Error("expected keys to be sorted alphabetically")
	}
}
