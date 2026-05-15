package envfile

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var interpolationPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// InterpolateOptions controls how variable interpolation behaves.
type InterpolateOptions struct {
	// FallbackToEnv allows falling back to the OS environment if a variable
	// is not found in the provided EnvFile.
	FallbackToEnv bool
	// Strict causes an error if a referenced variable cannot be resolved.
	Strict bool
}

// Interpolate resolves variable references within an EnvFile's values.
// References of the form $VAR or ${VAR} are replaced with their resolved values.
// Variables are resolved first from the EnvFile itself, then optionally from
// the OS environment if FallbackToEnv is set.
func Interpolate(ef EnvFile, opts InterpolateOptions) (EnvFile, error) {
	resolved := make(map[string]string, len(ef.Entries))
	for _, entry := range ef.Entries {
		resolved[entry.Key] = entry.Value
	}

	result := EnvFile{
		Path:    ef.Path,
		Entries: make([]Entry, 0, len(ef.Entries)),
	}

	for _, entry := range ef.Entries {
		interpolated, err := interpolateValue(entry.Value, resolved, opts)
		if err != nil {
			return EnvFile{}, fmt.Errorf("interpolating key %q: %w", entry.Key, err)
		}
		result.Entries = append(result.Entries, Entry{
			Key:   entry.Key,
			Value: interpolated,
		})
	}

	return result, nil
}

func interpolateValue(value string, vars map[string]string, opts InterpolateOptions) (string, error) {
	var resolveErr error

	result := interpolationPattern.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}

		varName := extractVarName(match)
		if v, ok := vars[varName]; ok {
			return v
		}
		if opts.FallbackToEnv {
			if v, ok := os.LookupEnv(varName); ok {
				return v
			}
		}
		if opts.Strict {
			resolveErr = fmt.Errorf("undefined variable %q", varName)
			return match
		}
		return ""
	})

	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}

func extractVarName(match string) string {
	if strings.HasPrefix(match, "${") {
		return match[2 : len(match)-1]
	}
	return match[1:]
}
