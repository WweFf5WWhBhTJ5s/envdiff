package envfile

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single key-value pair in an .env file.
type Entry struct {
	Key     string
	Value   string
	Comment string
	Line    int
}

// EnvFile holds all parsed entries from a .env file.
type EnvFile struct {
	Path    string
	Entries []Entry
	Index   map[string]*Entry
}

// Parse reads and parses an .env file from the given path.
func Parse(path string) (*EnvFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening env file: %w", err)
	}
	defer f.Close()

	env := &EnvFile{
		Path:  path,
		Index: make(map[string]*Entry),
	}

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		key, value, comment, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		entry := Entry{
			Key:     key,
			Value:   value,
			Comment: comment,
			Line:    lineNum,
		}
		env.Entries = append(env.Entries, entry)
		env.Index[key] = &env.Entries[len(env.Entries)-1]
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning env file: %w", err)
	}

	return env, nil
}

// parseLine splits a line into key, value, and optional inline comment.
func parseLine(line string) (key, value, comment string, err error) {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid format: %q", line)
	}

	key = strings.TrimSpace(parts[0])
	raw := strings.TrimSpace(parts[1])

	if idx := strings.Index(raw, " #"); idx != -1 {
		value = strings.TrimSpace(raw[:idx])
		comment = strings.TrimSpace(raw[idx+2:])
	} else {
		value = raw
	}

	value = strings.Trim(value, `"`)
	value = strings.Trim(value, `'`)

	return key, value, comment, nil
}
