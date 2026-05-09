package envfile

import (
	"testing"
)

func makeValidEnvFile(entries []Entry, values map[string]string) EnvFile {
	return EnvFile{Entries: entries, Values: values}
}

func TestValidate_ValidFile(t *testing.T) {
	ef := makeValidEnvFile(
		[]Entry{{Key: "APP_NAME", Value: "envdiff"}, {Key: "PORT", Value: "8080"}},
		map[string]string{"APP_NAME": "envdiff", "PORT": "8080"},
	)
	result := Validate(ef, []string{"APP_NAME", "PORT"})
	if result.HasErrors() {
		t.Fatalf("expected no errors, got: %s", result.Error())
	}
}

func TestValidate_DuplicateKey(t *testing.T) {
	ef := makeValidEnvFile(
		[]Entry{{Key: "FOO", Value: "bar"}, {Key: "FOO", Value: "baz"}},
		map[string]string{"FOO": "baz"},
	)
	result := Validate(ef, nil)
	if !result.HasErrors() {
		t.Fatal("expected duplicate key error")
	}
	if result.Errors[0].Key != "FOO" || result.Errors[0].Message != "duplicate key" {
		t.Errorf("unexpected error: %+v", result.Errors[0])
	}
}

func TestValidate_InvalidKeyCharacter(t *testing.T) {
	ef := makeValidEnvFile(
		[]Entry{{Key: "BAD-KEY", Value: "val"}},
		map[string]string{"BAD-KEY": "val"},
	)
	result := Validate(ef, nil)
	if !result.HasErrors() {
		t.Fatal("expected invalid character error")
	}
}

func TestValidate_KeyStartsWithDigit(t *testing.T) {
	ef := makeValidEnvFile(
		[]Entry{{Key: "1BAD", Value: "val"}},
		map[string]string{"1BAD": "val"},
	)
	result := Validate(ef, nil)
	if !result.HasErrors() {
		t.Fatal("expected digit-start error")
	}
	if result.Errors[0].Message != "key must not start with a digit" {
		t.Errorf("unexpected message: %s", result.Errors[0].Message)
	}
}

func TestValidate_MissingRequiredKey(t *testing.T) {
	ef := makeValidEnvFile(
		[]Entry{{Key: "APP_NAME", Value: "envdiff"}},
		map[string]string{"APP_NAME": "envdiff"},
	)
	result := Validate(ef, []string{"APP_NAME", "SECRET_KEY"})
	if !result.HasErrors() {
		t.Fatal("expected missing required key error")
	}
	if result.Errors[0].Key != "SECRET_KEY" {
		t.Errorf("expected SECRET_KEY error, got: %s", result.Errors[0].Key)
	}
}

func TestValidate_EmptyRequiredValue(t *testing.T) {
	ef := makeValidEnvFile(
		[]Entry{{Key: "TOKEN", Value: ""}},
		map[string]string{"TOKEN": ""},
	)
	result := Validate(ef, []string{"TOKEN"})
	if !result.HasErrors() {
		t.Fatal("expected empty required value error")
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	ef := makeValidEnvFile(
		[]Entry{{Key: "BAD KEY", Value: "x"}, {Key: "FOO", Value: ""}, {Key: "FOO", Value: "y"}},
		map[string]string{"BAD KEY": "x", "FOO": "y"},
	)
	result := Validate(ef, []string{"MISSING"})
	if len(result.Errors) < 3 {
		t.Errorf("expected at least 3 errors, got %d: %s", len(result.Errors), result.Error())
	}
}
