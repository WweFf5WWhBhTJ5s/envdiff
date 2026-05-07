package envfile

// Entry represents a single key-value pair in an .env file.
type Entry struct {
	Key   string
	Value string
}

// EnvFile represents a parsed .env file.
type EnvFile struct {
	Entries []Entry
}

// Get returns the value for the given key and whether it was found.
func (ef *EnvFile) Get(key string) (string, bool) {
	for _, e := range ef.Entries {
		if e.Key == key {
			return e.Value, true
		}
	}
	return "", false
}

// Keys returns all keys in the EnvFile in order.
func (ef *EnvFile) Keys() []string {
	keys := make([]string, len(ef.Entries))
	for i, e := range ef.Entries {
		keys[i] = e.Key
	}
	return keys
}

// DiffStatus represents the type of change between two env files.
type DiffStatus string

const (
	StatusUnchanged DiffStatus = "unchanged"
	StatusChanged   DiffStatus = "changed"
	StatusAdded     DiffStatus = "added"
	StatusRemoved   DiffStatus = "removed"
)

// DiffResult holds the comparison result for a single key.
type DiffResult struct {
	Key    string
	Status DiffStatus
	OldVal string
	NewVal string
}
