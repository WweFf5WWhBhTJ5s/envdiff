package envfile

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ExportFormat represents a supported export format.
type ExportFormat string

const (
	ExportJSON   ExportFormat = "json"
	ExportShell  ExportFormat = "shell"
	ExportDotenv ExportFormat = "dotenv"
)

// Convert serializes an EnvFile into the given export format.
func Convert(ef EnvFile, format ExportFormat) (string, error) {
	switch format {
	case ExportJSON:
		return convertJSON(ef)
	case ExportShell:
		return convertShell(ef), nil
	case ExportDotenv:
		return convertDotenv(ef), nil
	default:
		return "", fmt.Errorf("unsupported export format: %q", format)
	}
}

func convertJSON(ef EnvFile) (string, error) {
	m := make(map[string]string, len(ef.Entries))
	for _, e := range ef.Entries {
		m[e.Key] = e.Value
	}
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(b), nil
}

func convertShell(ef EnvFile) string {
	var sb strings.Builder
	for _, e := range ef.Entries {
		// Quote the value to handle spaces and special chars safely.
		quoted := strings.ReplaceAll(e.Value, `'`, `'\''`)
		fmt.Fprintf(&sb, "export %s='%s'\n", e.Key, quoted)
	}
	return sb.String()
}

func convertDotenv(ef EnvFile) string {
	var sb strings.Builder
	for _, e := range ef.Entries {
		if strings.ContainsAny(e.Value, " \t\n#") {
			fmt.Fprintf(&sb, "%s=%q\n", e.Key, e.Value)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
		}
	}
	return sb.String()
}

// ConvertMap converts a raw map to an EnvFile, sorting keys for determinism.
func ConvertMap(m map[string]string) EnvFile {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ef := EnvFile{}
	for _, k := range keys {
		ef.Entries = append(ef.Entries, Entry{Key: k, Value: m[k]})
	}
	return ef
}
