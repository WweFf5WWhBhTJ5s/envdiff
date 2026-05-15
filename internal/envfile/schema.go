package envfile

import (
	"fmt"
	"regexp"
	"strings"
)

// SchemaRule defines a validation rule for a specific key.
type SchemaRule struct {
	Key      string
	Required bool
	Pattern  *regexp.Regexp
	Allowed  []string
}

// Schema holds a set of rules to validate an EnvFile against.
type Schema struct {
	Rules []SchemaRule
}

// SchemaViolation represents a single schema validation failure.
type SchemaViolation struct {
	Key     string
	Message string
}

// SchemaResult holds the outcome of a schema validation.
type SchemaResult struct {
	Violations []SchemaViolation
	Valid      bool
}

// ValidateSchema checks an EnvFile against the provided Schema.
func ValidateSchema(ef EnvFile, schema Schema) SchemaResult {
	result := SchemaResult{Valid: true}

	keyMap := make(map[string]string, len(ef.Entries))
	for _, entry := range ef.Entries {
		keyMap[entry.Key] = entry.Value
	}

	for _, rule := range schema.Rules {
		val, exists := keyMap[rule.Key]

		if rule.Required && !exists {
			result.Violations = append(result.Violations, SchemaViolation{
				Key:     rule.Key,
				Message: "required key is missing",
			})
			result.Valid = false
			continue
		}

		if !exists {
			continue
		}

		if rule.Pattern != nil && !rule.Pattern.MatchString(val) {
			result.Violations = append(result.Violations, SchemaViolation{
				Key:     rule.Key,
				Message: fmt.Sprintf("value %q does not match required pattern %s", val, rule.Pattern.String()),
			})
			result.Valid = false
		}

		if len(rule.Allowed) > 0 {
			found := false
			for _, a := range rule.Allowed {
				if strings.EqualFold(val, a) {
					found = true
					break
				}
			}
			if !found {
				result.Violations = append(result.Violations, SchemaViolation{
					Key:     rule.Key,
					Message: fmt.Sprintf("value %q is not in allowed list [%s]", val, strings.Join(rule.Allowed, ", ")),
				})
				result.Valid = false
			}
		}
	}

	return result
}
