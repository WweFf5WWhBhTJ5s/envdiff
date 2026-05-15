package envfile

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// RedactMode controls how sensitive values are redacted in output.
type RedactMode string

const (
	RedactMask   RedactMode = "mask"   // Replace with ***
	RedactHash   RedactMode = "hash"   // Replace with sha256 prefix
	RedactBlank  RedactMode = "blank"  // Replace with empty string
	RedactLength RedactMode = "length" // Replace with length hint, e.g. <12 chars>
)

// RedactOptions configures the redaction behaviour.
type RedactOptions struct {
	Mode     RedactMode
	MaskChar string // used for RedactMask, default "*"
	MaskLen  int    // fixed mask length, 0 = use actual length
}

// DefaultRedactOptions returns sensible defaults.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Mode:     RedactMask,
		MaskChar: "*",
		MaskLen:  6,
	}
}

// Redactor applies redaction to env file entries.
type Redactor struct {
	masker *Masker
	opts   RedactOptions
}

// NewRedactor creates a Redactor using the default masker and options.
func NewRedactor(opts RedactOptions) *Redactor {
	return &Redactor{masker: NewMasker(), opts: opts}
}

// NewRedactorWithMasker creates a Redactor with a custom Masker.
func NewRedactorWithMasker(m *Masker, opts RedactOptions) *Redactor {
	return &Redactor{masker: m, opts: opts}
}

// RedactValue redacts a single value if the key is considered sensitive.
func (r *Redactor) RedactValue(key, value string) string {
	if !r.masker.IsSensitive(key) {
		return value
	}
	return r.applyMode(value)
}

// RedactEnvFile returns a copy of the EnvFile with sensitive values redacted.
func (r *Redactor) RedactEnvFile(ef EnvFile) EnvFile {
	out := EnvFile{Path: ef.Path, Entries: make([]EnvEntry, len(ef.Entries))}
	for i, e := range ef.Entries {
		out.Entries[i] = EnvEntry{
			Key:   e.Key,
			Value: r.RedactValue(e.Key, e.Value),
		}
	}
	return out
}

func (r *Redactor) applyMode(value string) string {
	switch r.opts.Mode {
	case RedactHash:
		h := sha256.Sum256([]byte(value))
		return fmt.Sprintf("sha256:%x", h[:4])
	case RedactBlank:
		return ""
	case RedactLength:
		return fmt.Sprintf("<%d chars>", len(value))
	default: // RedactMask
		ch := r.opts.MaskChar
		if ch == "" {
			ch = "*"
		}
		l := r.opts.MaskLen
		if l <= 0 {
			l = len(value)
		}
		return strings.Repeat(ch, l)
	}
}
