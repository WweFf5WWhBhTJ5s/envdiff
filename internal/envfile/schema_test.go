package envfile

import (
	"regexp"
	"testing"
)

func makeSchemaEnv(pairs map[string]string) EnvFile {
	ef := EnvFile{}
	for k, v := range pairs {
		ef.Entries = append(ef.Entries, EnvEntry{Key: k, Value: v})
	}
	return ef
}

func TestValidateSchema_AllPresent(t *testing.T) {
	ef := makeSchemaEnv(map[string]string{"APP_ENV": "production", "PORT": "8080"})
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "APP_ENV", Required: true},
			{Key: "PORT", Required: true},
		},
	}
	res := ValidateSchema(ef, schema)
	if !res.Valid {
		t.Fatalf("expected valid, got violations: %+v", res.Violations)
	}
}

func TestValidateSchema_MissingRequired(t *testing.T) {
	ef := makeSchemaEnv(map[string]string{"PORT": "8080"})
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "APP_ENV", Required: true},
		},
	}
	res := ValidateSchema(ef, schema)
	if res.Valid {
		t.Fatal("expected invalid due to missing required key")
	}
	if len(res.Violations) != 1 || res.Violations[0].Key != "APP_ENV" {
		t.Errorf("unexpected violations: %+v", res.Violations)
	}
}

func TestValidateSchema_PatternMismatch(t *testing.T) {
	ef := makeSchemaEnv(map[string]string{"PORT": "abc"})
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "PORT", Required: true, Pattern: regexp.MustCompile(`^\d+$`)},
		},
	}
	res := ValidateSchema(ef, schema)
	if res.Valid {
		t.Fatal("expected invalid due to pattern mismatch")
	}
}

func TestValidateSchema_AllowedValues(t *testing.T) {
	ef := makeSchemaEnv(map[string]string{"LOG_LEVEL": "verbose"})
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "LOG_LEVEL", Required: true, Allowed: []string{"debug", "info", "warn", "error"}},
		},
	}
	res := ValidateSchema(ef, schema)
	if res.Valid {
		t.Fatal("expected invalid due to disallowed value")
	}
}

func TestValidateSchema_OptionalMissingKey(t *testing.T) {
	ef := makeSchemaEnv(map[string]string{"APP_ENV": "staging"})
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "OPTIONAL_KEY", Required: false, Pattern: regexp.MustCompile(`^[a-z]+$`)},
		},
	}
	res := ValidateSchema(ef, schema)
	if !res.Valid {
		t.Fatalf("expected valid when optional key is absent, got: %+v", res.Violations)
	}
}

func TestValidateSchema_MultipleViolations(t *testing.T) {
	ef := makeSchemaEnv(map[string]string{"LOG_LEVEL": "trace"})
	schema := Schema{
		Rules: []SchemaRule{
			{Key: "APP_ENV", Required: true},
			{Key: "LOG_LEVEL", Required: true, Allowed: []string{"debug", "info"}},
		},
	}
	res := ValidateSchema(ef, schema)
	if res.Valid {
		t.Fatal("expected invalid")
	}
	if len(res.Violations) != 2 {
		t.Errorf("expected 2 violations, got %d: %+v", len(res.Violations), res.Violations)
	}
}
