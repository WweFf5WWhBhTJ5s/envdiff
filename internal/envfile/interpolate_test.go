package envfile

import (
	"os"
	"testing"
)

func makeInterpolateEnv(pairs ...string) EnvFile {
	ef := EnvFile{}
	for i := 0; i+1 < len(pairs); i += 2 {
		ef.Entries = append(ef.Entries, Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return ef
}

func TestInterpolate_NoReferences(t *testing.T) {
	ef := makeInterpolateEnv("FOO", "bar", "BAZ", "qux")
	out, err := Interpolate(ef, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Entries[0].Value != "bar" || out.Entries[1].Value != "qux" {
		t.Errorf("expected unchanged values, got %+v", out.Entries)
	}
}

func TestInterpolate_BraceStyle(t *testing.T) {
	ef := makeInterpolateEnv("BASE", "/opt/app", "LOG_DIR", "${BASE}/logs")
	out, err := Interpolate(ef, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.Entries[1].Value; got != "/opt/app/logs" {
		t.Errorf("expected /opt/app/logs, got %q", got)
	}
}

func TestInterpolate_DollarStyle(t *testing.T) {
	ef := makeInterpolateEnv("HOST", "localhost", "URL", "http://$HOST:8080")
	out, err := Interpolate(ef, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.Entries[1].Value; got != "http://localhost:8080" {
		t.Errorf("expected http://localhost:8080, got %q", got)
	}
}

func TestInterpolate_FallbackToEnv(t *testing.T) {
	os.Setenv("_TEST_ENVDIFF_HOST", "envhost")
	defer os.Unsetenv("_TEST_ENVDIFF_HOST")

	ef := makeInterpolateEnv("URL", "http://${_TEST_ENVDIFF_HOST}:9000")
	out, err := Interpolate(ef, InterpolateOptions{FallbackToEnv: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.Entries[0].Value; got != "http://envhost:9000" {
		t.Errorf("expected http://envhost:9000, got %q", got)
	}
}

func TestInterpolate_UndefinedSilent(t *testing.T) {
	ef := makeInterpolateEnv("VAL", "prefix_${UNDEFINED}_suffix")
	out, err := Interpolate(ef, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out.Entries[0].Value; got != "prefix__suffix" {
		t.Errorf("expected empty substitution, got %q", got)
	}
}

func TestInterpolate_StrictErrorsOnUndefined(t *testing.T) {
	ef := makeInterpolateEnv("VAL", "${MISSING_VAR}")
	_, err := Interpolate(ef, InterpolateOptions{Strict: true})
	if err == nil {
		t.Fatal("expected error for undefined variable in strict mode")
	}
}

func TestInterpolate_PreservesPath(t *testing.T) {
	ef := makeInterpolateEnv("A", "hello")
	ef.Path = "/some/path/.env"
	out, err := Interpolate(ef, InterpolateOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Path != ef.Path {
		t.Errorf("path not preserved: got %q", out.Path)
	}
}
