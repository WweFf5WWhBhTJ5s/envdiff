package envfile

import (
	"fmt"
	"time"
)

// AuditAction represents the type of change recorded in an audit entry.
type AuditAction string

const (
	AuditAdded   AuditAction = "added"
	AuditRemoved AuditAction = "removed"
	AuditChanged AuditAction = "changed"
)

// AuditEntry records a single change event for an env key.
type AuditEntry struct {
	Timestamp time.Time
	Key       string
	Action    AuditAction
	OldValue  string
	NewValue  string
	Source    string
}

// AuditLog holds a collection of audit entries.
type AuditLog struct {
	Entries []AuditEntry
}

// Audit generates an AuditLog from a slice of DiffResult entries.
// The source label identifies which file pair produced the diff.
func Audit(results []DiffResult, source string, masker *Masker) *AuditLog {
	log := &AuditLog{}
	now := time.Now().UTC()

	for _, r := range results {
		if r.Status == StatusUnchanged {
			continue
		}

		oldVal := r.OldValue
		newVal := r.NewValue

		if masker != nil && masker.IsSensitive(r.Key) {
			oldVal = masker.MaskValue(oldVal)
			newVal = masker.MaskValue(newVal)
		}

		entry := AuditEntry{
			Timestamp: now,
			Key:       r.Key,
			Action:    diffStatusToAction(r.Status),
			OldValue:  oldVal,
			NewValue:  newVal,
			Source:    source,
		}
		log.Entries = append(log.Entries, entry)
	}

	return log
}

func diffStatusToAction(s DiffStatus) AuditAction {
	switch s {
	case StatusAdded:
		return AuditAdded
	case StatusRemoved:
		return AuditRemoved
	case StatusChanged:
		return AuditChanged
	default:
		return AuditAction(fmt.Sprintf("unknown(%s)", s))
	}
}
