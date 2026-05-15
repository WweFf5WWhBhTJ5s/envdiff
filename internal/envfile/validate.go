package envfile

import (
	"fmt"
	"strings"
	"unicode"
)

// ValidationError represents a single validation issue found in an env file.
type ValidationError struct {
	Key     string
	Line    int
	Message string
}

func (e ValidationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("line %d: key %q: %s", e.Line, e.Key, e.Message)
	}
	return fmt.Sprintf("key %q: %s", e.Key, e.Message)
}

// ValidationResult holds all errors found during validation.
type ValidationResult struct {
	Errors []ValidationError
}

func (r *ValidationResult) HasErrors() bool {
	return len(r.Errors) > 0
}

func (r *ValidationResult) Error() string {
	msgs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		msgs[i] = e.Error()
	}
	return strings.Join(msgs, "\n")
}

// Validate checks an EnvFile for common issues such as empty keys,
// invalid key characters, and empty values for required keys.
func Validate(ef EnvFile, requiredKeys []string) ValidationResult {
	result := ValidationResult{}

	seen := make(map[string]bool)
	for i, entry := range ef.Entries {
		lineNum := i + 1

		if entry.Key == "" {
			result.Errors = append(result.Errors, ValidationError{
				Line:    lineNum,
				Message: "empty key is not allowed",
			})
			continue
		}

		if err := validateKeyName(entry.Key); err != nil {
			result.Errors = append(result.Errors, ValidationError{
				Key:     entry.Key,
				Line:    lineNum,
				Message: err.Error(),
			})
		}

		if seen[entry.Key] {
			result.Errors = append(result.Errors, ValidationError{
				Key:     entry.Key,
				Line:    lineNum,
				Message: "duplicate key",
			})
		}
		seen[entry.Key] = true
	}

	for _, req := range requiredKeys {
		val, ok := ef.Values[req]
		if !ok || strings.TrimSpace(val) == "" {
			result.Errors = append(result.Errors, ValidationError{
				Key:     req,
				Message: "required key is missing or empty",
			})
		}
	}

	return result
}

// FilterErrors returns only the ValidationErrors that match the given key.
// This is useful when callers need to inspect errors for a specific key
// without iterating over the full error list.
func (r *ValidationResult) FilterErrors(key string) []ValidationError {
	var matched []ValidationError
	for _, e := range r.Errors {
		if e.Key == key {
			matched = append(matched, e)
		}
	}
	return matched
}

func validateKeyName(key string) error {
	for i, ch := range key {
		if i == 0 && unicode.IsDigit(ch) {
			return fmt.Errorf("key must not start with a digit")
		}
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '_' {
			return fmt.Errorf("key contains invalid character %q", ch)
		}
	}
	return nil
}
