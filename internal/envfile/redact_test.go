package envfile

import (
	"strings"
	"testing"
)

func makeRedactEnv(pairs ...string) EnvFile {
	ef := EnvFile{Path: ".env"}
	for i := 0; i+1 < len(pairs); i += 2 {
		ef.Entries = append(ef.Entries, EnvEntry{Key: pairs[i], Value: pairs[i+1]})
	}
	return ef
}

func TestRedactor_NonSensitivePassthrough(t *testing.T) {
	r := NewRedactor(DefaultRedactOptions())
	got := r.RedactValue("APP_NAME", "myapp")
	if got != "myapp" {
		t.Errorf("expected passthrough, got %q", got)
	}
}

func TestRedactor_MaskMode(t *testing.T) {
	opts := RedactOptions{Mode: RedactMask, MaskChar: "*", MaskLen: 6}
	r := NewRedactor(opts)
	got := r.RedactValue("SECRET_KEY", "supersecret")
	if got != "******" {
		t.Errorf("expected ******, got %q", got)
	}
}

func TestRedactor_HashMode(t *testing.T) {
	opts := RedactOptions{Mode: RedactHash}
	r := NewRedactor(opts)
	got := r.RedactValue("API_SECRET", "abc123")
	if !strings.HasPrefix(got, "sha256:") {
		t.Errorf("expected sha256 prefix, got %q", got)
	}
}

func TestRedactor_BlankMode(t *testing.T) {
	opts := RedactOptions{Mode: RedactBlank}
	r := NewRedactor(opts)
	got := r.RedactValue("DB_PASSWORD", "hunter2")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestRedactor_LengthMode(t *testing.T) {
	opts := RedactOptions{Mode: RedactLength}
	r := NewRedactor(opts)
	got := r.RedactValue("AUTH_TOKEN", "abcdefgh")
	if got != "<8 chars>" {
		t.Errorf("expected '<8 chars>', got %q", got)
	}
}

func TestRedactor_RedactEnvFile(t *testing.T) {
	ef := makeRedactEnv(
		"APP_ENV", "production",
		"DB_PASSWORD", "s3cr3t",
		"PORT", "8080",
		"API_KEY", "key-xyz",
	)
	r := NewRedactor(DefaultRedactOptions())
	out := r.RedactEnvFile(ef)

	for _, e := range out.Entries {
		switch e.Key {
		case "APP_ENV", "PORT":
			if e.Value == "******" {
				t.Errorf("key %s should not be redacted", e.Key)
			}
		case "DB_PASSWORD", "API_KEY":
			if e.Value != "******" {
				t.Errorf("key %s should be redacted, got %q", e.Key, e.Value)
			}
		}
	}
}

func TestRedactor_CustomMasker(t *testing.T) {
	m := NewMaskerWithPatterns([]string{"CUSTOM_SECRET"})
	r := NewRedactorWithMasker(m, RedactOptions{Mode: RedactMask, MaskChar: "#", MaskLen: 4})
	got := r.RedactValue("CUSTOM_SECRET", "value")
	if got != "####" {
		t.Errorf("expected ####, got %q", got)
	}
	// Standard sensitive key should NOT be masked (custom masker only has CUSTOM_SECRET)
	got2 := r.RedactValue("API_KEY", "value")
	if got2 == "####" {
		t.Errorf("API_KEY should not be masked by custom masker")
	}
}
