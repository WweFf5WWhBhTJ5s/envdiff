package envfile

import (
	"encoding/json"
	"strings"
	"testing"
)

func makeConvertEnv() EnvFile {
	return EnvFile{
		Entries: []Entry{
			{Key: "APP_NAME", Value: "envdiff"},
			{Key: "DEBUG", Value: "false"},
			{Key: "DB_URL", Value: "postgres://localhost/db"},
		},
	}
}

func TestConvert_JSON(t *testing.T) {
	ef := makeConvertEnv()
	out, err := Convert(ef, ExportJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if m["APP_NAME"] != "envdiff" {
		t.Errorf("expected APP_NAME=envdiff, got %q", m["APP_NAME"])
	}
	if m["DEBUG"] != "false" {
		t.Errorf("expected DEBUG=false, got %q", m["DEBUG"])
	}
}

func TestConvert_Shell(t *testing.T) {
	ef := makeConvertEnv()
	out, err := Convert(ef, ExportShell)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export APP_NAME='envdiff'") {
		t.Errorf("expected shell export line, got:\n%s", out)
	}
	if !strings.Contains(out, "export DEBUG='false'") {
		t.Errorf("expected shell export for DEBUG, got:\n%s", out)
	}
}

func TestConvert_Shell_QuotesSingleQuote(t *testing.T) {
	ef := EnvFile{Entries: []Entry{{Key: "MSG", Value: "it's alive"}}}
	out, err := Convert(ef, ExportShell)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `export MSG='it'\''s alive'`) {
		t.Errorf("single quote not escaped properly, got:\n%s", out)
	}
}

func TestConvert_Dotenv(t *testing.T) {
	ef := makeConvertEnv()
	out, err := Convert(ef, ExportDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "APP_NAME=envdiff") {
		t.Errorf("expected dotenv line, got:\n%s", out)
	}
}

func TestConvert_Dotenv_QuotesSpaces(t *testing.T) {
	ef := EnvFile{Entries: []Entry{{Key: "GREETING", Value: "hello world"}}}
	out, err := Convert(ef, ExportDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `GREETING="hello world"`) {
		t.Errorf("expected quoted value for space, got:\n%s", out)
	}
}

func TestConvert_UnknownFormat(t *testing.T) {
	ef := makeConvertEnv()
	_, err := Convert(ef, "xml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
}

func TestConvertMap(t *testing.T) {
	m := map[string]string{"Z_KEY": "z", "A_KEY": "a", "M_KEY": "m"}
	ef := ConvertMap(m)
	if len(ef.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(ef.Entries))
	}
	if ef.Entries[0].Key != "A_KEY" {
		t.Errorf("expected first key A_KEY, got %s", ef.Entries[0].Key)
	}
	if ef.Entries[2].Key != "Z_KEY" {
		t.Errorf("expected last key Z_KEY, got %s", ef.Entries[2].Key)
	}
}
