package envfile

import (
	"regexp"
	"strings"
)

// DefaultSensitivePatterns contains common patterns for sensitive keys.
var DefaultSensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|pwd)`),
	regexp.MustCompile(`(?i)(secret|token|api_key|apikey)`),
	regexp.MustCompile(`(?i)(private_key|privatekey)`),
	regexp.MustCompile(`(?i)(auth|credential|credentials)`),
}

const maskValue = "***"

// Masker holds configuration for masking sensitive values.
type Masker struct {
	patterns []*regexp.Regexp
}

// NewMasker creates a Masker with the default sensitive patterns.
func NewMasker() *Masker {
	return &Masker{patterns: DefaultSensitivePatterns}
}

// NewMaskerWithPatterns creates a Masker with custom patterns.
func NewMaskerWithPatterns(patterns []*regexp.Regexp) *Masker {
	return &Masker{patterns: patterns}
}

// IsSensitive reports whether the given key matches any sensitive pattern.
func (m *Masker) IsSensitive(key string) bool {
	for _, p := range m.patterns {
		if p.MatchString(key) {
			return true
		}
	}
	return false
}

// MaskValue returns the masked placeholder if the key is sensitive,
// otherwise it returns the original value unchanged.
func (m *Masker) MaskValue(key, value string) string {
	if m.IsSensitive(key) {
		return maskValue
	}
	return value
}

// MaskEnvFile returns a new EnvFile with sensitive values replaced by the
// mask placeholder. The original EnvFile is not modified.
func (m *Masker) MaskEnvFile(ef *EnvFile) *EnvFile {
	masked := &EnvFile{
		Entries: make([]Entry, len(ef.Entries)),
	}
	for i, e := range ef.Entries {
		masked.Entries[i] = Entry{
			Key:   e.Key,
			Value: m.MaskValue(e.Key, e.Value),
		}
	}
	return masked
}

// MaskDiffResults returns a copy of the diff results with sensitive values masked.
func (m *Masker) MaskDiffResults(results []DiffResult) []DiffResult {
	out := make([]DiffResult, len(results))
	for i, r := range results {
		out[i] = DiffResult{
			Key:    r.Key,
			Status: r.Status,
			OldVal: m.MaskValue(r.Key, r.OldVal),
			NewVal: m.MaskValue(r.Key, r.NewVal),
		}
	}
	return out
}

// MaskString masks a sensitive value string for display.
func MaskString(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.Repeat("*", len(maskValue))
}
