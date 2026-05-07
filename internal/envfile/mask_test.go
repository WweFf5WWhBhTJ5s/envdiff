package envfile

import (
	"testing"
)

func TestMasker_IsSensitive(t *testing.T) {
	m := NewMasker()

	sensitiveKeys := []string{
		"DB_PASSWORD", "API_SECRET", "AUTH_TOKEN",
		"PRIVATE_KEY", "AWS_API_KEY", "USER_CREDENTIALS",
		"password", "secret", "api_key",
	}
	for _, key := range sensitiveKeys {
		if !m.IsSensitive(key) {
			t.Errorf("expected key %q to be sensitive", key)
		}
	}

	safeKeys := []string{"HOST", "PORT", "APP_ENV", "DEBUG", "LOG_LEVEL"}
	for _, key := range safeKeys {
		if m.IsSensitive(key) {
			t.Errorf("expected key %q to be safe", key)
		}
	}
}

func TestMasker_MaskValue(t *testing.T) {
	m := NewMasker()

	if got := m.MaskValue("DB_PASSWORD", "supersecret"); got != "***" {
		t.Errorf("expected '***', got %q", got)
	}
	if got := m.MaskValue("HOST", "localhost"); got != "localhost" {
		t.Errorf("expected 'localhost', got %q", got)
	}
}

func TestMasker_MaskEnvFile(t *testing.T) {
	ef := &EnvFile{
		Entries: []Entry{
			{Key: "HOST", Value: "localhost"},
			{Key: "DB_PASSWORD", Value: "s3cr3t"},
			{Key: "PORT", Value: "5432"},
			{Key: "API_KEY", Value: "abc123"},
		},
	}

	m := NewMasker()
	masked := m.MaskEnvFile(ef)

	expected := map[string]string{
		"HOST":        "localhost",
		"DB_PASSWORD": "***",
		"PORT":        "5432",
		"API_KEY":     "***",
	}
	for _, e := range masked.Entries {
		if want, ok := expected[e.Key]; ok {
			if e.Value != want {
				t.Errorf("key %q: expected %q, got %q", e.Key, want, e.Value)
			}
		}
	}

	// Ensure original is not modified
	for _, e := range ef.Entries {
		if e.Key == "DB_PASSWORD" && e.Value != "s3cr3t" {
			t.Error("original EnvFile was modified")
		}
	}
}

func TestMasker_MaskDiffResults(t *testing.T) {
	m := NewMasker()

	results := []DiffResult{
		{Key: "HOST", Status: StatusChanged, OldVal: "old", NewVal: "new"},
		{Key: "DB_PASSWORD", Status: StatusChanged, OldVal: "oldpass", NewVal: "newpass"},
		{Key: "API_TOKEN", Status: StatusAdded, OldVal: "", NewVal: "tok123"},
	}

	masked := m.MaskDiffResults(results)

	if masked[0].OldVal != "old" || masked[0].NewVal != "new" {
		t.Error("safe key values should not be masked")
	}
	if masked[1].OldVal != "***" || masked[1].NewVal != "***" {
		t.Error("DB_PASSWORD values should be masked")
	}
	if masked[2].NewVal != "***" {
		t.Error("API_TOKEN value should be masked")
	}
}
