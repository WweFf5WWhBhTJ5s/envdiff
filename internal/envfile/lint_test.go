package envfile

import (
	"testing"
)

func makeLintEnv(entries []EnvEntry) EnvFile {
	return EnvFile{Entries: entries}
}

func TestLint_CleanFile(t *testing.T) {
	ef := makeLintEnv([]EnvEntry{
		{Key: "DATABASE_URL", Value: "postgres://localhost/db", Quoted: false},
		{Key: "API_KEY", Value: "abc123", Quoted: false},
	})
	result := Lint(ef)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(result.Issues), result.Issues)
	}
}

func TestLint_EmptyValue(t *testing.T) {
	ef := makeLintEnv([]EnvEntry{
		{Key: "EMPTY_KEY", Value: "", Quoted: false},
	})
	result := Lint(ef)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != LintWarning {
		t.Errorf("expected warning, got %s", result.Issues[0].Severity)
	}
	if result.HasErrors() {
		t.Error("expected no errors")
	}
}

func TestLint_LowercaseKey(t *testing.T) {
	ef := makeLintEnv([]EnvEntry{
		{Key: "my_key", Value: "value", Quoted: false},
	})
	result := Lint(ef)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != LintWarning {
		t.Errorf("expected warning severity")
	}
}

func TestLint_UnquotedWhitespace(t *testing.T) {
	ef := makeLintEnv([]EnvEntry{
		{Key: "SPACED_VALUE", Value: "hello world", Quoted: false},
	})
	result := Lint(ef)
	if len(result.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(result.Issues))
	}
	if result.Issues[0].Severity != LintError {
		t.Errorf("expected error severity, got %s", result.Issues[0].Severity)
	}
	if !result.HasErrors() {
		t.Error("expected HasErrors to return true")
	}
}

func TestLint_QuotedWhitespaceIsOK(t *testing.T) {
	ef := makeLintEnv([]EnvEntry{
		{Key: "SPACED_VALUE", Value: "hello world", Quoted: true},
	})
	result := Lint(ef)
	if len(result.Issues) != 0 {
		t.Errorf("expected no issues for quoted whitespace, got %d", len(result.Issues))
	}
}

func TestLint_MultipleIssues(t *testing.T) {
	ef := makeLintEnv([]EnvEntry{
		{Key: "lower_key", Value: "", Quoted: false},
	})
	result := Lint(ef)
	// expects both empty value warning AND lowercase warning
	if len(result.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(result.Issues))
	}
}
