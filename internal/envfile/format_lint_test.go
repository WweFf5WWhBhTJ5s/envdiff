package envfile

import (
	"strings"
	"testing"
)

func sampleLintResult() LintResult {
	return LintResult{
		Issues: []LintIssue{
			{Line: 1, Key: "low_key", Message: `key "low_key" is not uppercase; consider "LOW_KEY"`, Severity: LintWarning},
			{Line: 3, Key: "BAD VALUE", Message: `key "BAD VALUE" has unquoted whitespace in value`, Severity: LintError},
		},
	}
}

func TestFormatLint_NoIssues(t *testing.T) {
	result := LintResult{}
	out := FormatLintResult(result, "text")
	if !strings.Contains(out, "No lint issues") {
		t.Errorf("expected no-issues message, got: %s", out)
	}
}

func TestFormatLint_Text(t *testing.T) {
	out := FormatLintResult(sampleLintResult(), "text")
	if !strings.Contains(out, "⚠") {
		t.Error("expected warning icon in text output")
	}
	if !strings.Contains(out, "✖") {
		t.Error("expected error icon in text output")
	}
	if !strings.Contains(out, "low_key") {
		t.Error("expected key name in output")
	}
}

func TestFormatLint_Table(t *testing.T) {
	out := FormatLintResult(sampleLintResult(), "table")
	if !strings.Contains(out, "SEVERITY") {
		t.Error("expected SEVERITY header in table output")
	}
	if !strings.Contains(out, "WARNING") {
		t.Error("expected WARNING in table output")
	}
	if !strings.Contains(out, "ERROR") {
		t.Error("expected ERROR in table output")
	}
}

func TestFormatLint_JSON(t *testing.T) {
	out := FormatLintResult(sampleLintResult(), "json")
	if !strings.HasPrefix(out, `{"issues":[`) {
		t.Errorf("expected JSON object, got: %s", out)
	}
	if !strings.Contains(out, `"severity"`) {
		t.Error("expected severity field in JSON")
	}
}

func TestFormatLint_DefaultIsText(t *testing.T) {
	outDefault := FormatLintResult(sampleLintResult(), "")
	outText := FormatLintResult(sampleLintResult(), "text")
	if outDefault != outText {
		t.Error("expected default format to match text format")
	}
}
