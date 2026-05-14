package envfile

import (
	"fmt"
	"strings"
)

// LintSeverity represents the severity level of a lint issue.
type LintSeverity string

const (
	LintWarning LintSeverity = "warning"
	LintError   LintSeverity = "error"
)

// LintIssue represents a single lint finding.
type LintIssue struct {
	Line     int
	Key      string
	Message  string
	Severity LintSeverity
}

func (i LintIssue) String() string {
	if i.Line > 0 {
		return fmt.Sprintf("%s (line %d): %s", i.Severity, i.Line, i.Message)
	}
	return fmt.Sprintf("%s: %s", i.Severity, i.Message)
}

// LintResult holds all issues found during linting.
type LintResult struct {
	Issues []LintIssue
}

func (r *LintResult) HasErrors() bool {
	for _, issue := range r.Issues {
		if issue.Severity == LintError {
			return true
		}
	}
	return false
}

// Lint checks an EnvFile for style and correctness issues beyond basic validation.
func Lint(ef EnvFile) LintResult {
	result := LintResult{}

	for i, entry := range ef.Entries {
		lineNum := i + 1

		// Warn on empty values
		if entry.Value == "" {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Key:      entry.Key,
				Message:  fmt.Sprintf("key %q has an empty value", entry.Key),
				Severity: LintWarning,
			})
		}

		// Warn on keys that are not UPPER_SNAKE_CASE
		if entry.Key != strings.ToUpper(entry.Key) {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Key:      entry.Key,
				Message:  fmt.Sprintf("key %q is not uppercase; consider %q", entry.Key, strings.ToUpper(entry.Key)),
				Severity: LintWarning,
			})
		}

		// Error on values with unquoted whitespace
		if strings.ContainsAny(entry.Value, " \t") && !entry.Quoted {
			result.Issues = append(result.Issues, LintIssue{
				Line:     lineNum,
				Key:      entry.Key,
				Message:  fmt.Sprintf("key %q has unquoted whitespace in value", entry.Key),
				Severity: LintError,
			})
		}
	}

	return result
}
