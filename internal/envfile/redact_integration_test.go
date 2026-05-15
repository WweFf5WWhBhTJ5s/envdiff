package envfile_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/envfile"
)

func writeRedactTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp env: %v", err)
	}
	return p
}

func TestRedact_ParseAndRedact(t *testing.T) {
	path := writeRedactTempEnv(t, "APP_ENV=production\nDB_PASSWORD=hunter2\nAPI_KEY=abc123\nPORT=3000\n")

	ef, err := envfile.Parse(path)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	r := envfile.NewRedactor(envfile.DefaultRedactOptions())
	out := r.RedactEnvFile(ef)

	values := make(map[string]string)
	for _, e := range out.Entries {
		values[e.Key] = e.Value
	}

	if values["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should be unchanged, got %q", values["APP_ENV"])
	}
	if values["PORT"] != "3000" {
		t.Errorf("PORT should be unchanged, got %q", values["PORT"])
	}
	if values["DB_PASSWORD"] != "******" {
		t.Errorf("DB_PASSWORD should be masked, got %q", values["DB_PASSWORD"])
	}
	if values["API_KEY"] != "******" {
		t.Errorf("API_KEY should be masked, got %q", values["API_KEY"])
	}
}

func TestRedact_HashIsDeterministic(t *testing.T) {
	opts := envfile.RedactOptions{Mode: envfile.RedactHash}
	r := envfile.NewRedactor(opts)

	v1 := r.RedactValue("SECRET_KEY", "same-value")
	v2 := r.RedactValue("SECRET_KEY", "same-value")
	if v1 != v2 {
		t.Errorf("hash should be deterministic: %q != %q", v1, v2)
	}
	if !strings.HasPrefix(v1, "sha256:") {
		t.Errorf("expected sha256 prefix, got %q", v1)
	}
}

func TestRedact_DifferentValuesProduceDifferentHashes(t *testing.T) {
	opts := envfile.RedactOptions{Mode: envfile.RedactHash}
	r := envfile.NewRedactor(opts)

	h1 := r.RedactValue("API_KEY", "value-one")
	h2 := r.RedactValue("API_KEY", "value-two")
	if h1 == h2 {
		t.Errorf("different values should produce different hashes")
	}
}
