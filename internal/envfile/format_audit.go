package envfile

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"
	"bytes"
)

// FormatAuditLog renders an AuditLog in the requested format.
// Supported formats: "text", "table", "json".
func FormatAuditLog(log *AuditLog, format string) string {
	switch strings.ToLower(format) {
	case "table":
		return formatAuditTable(log)
	case "json":
		return formatAuditJSON(log)
	default:
		return formatAuditText(log)
	}
}

func formatAuditText(log *AuditLog) string {
	if len(log.Entries) == 0 {
		return "No audit entries.\n"
	}
	var sb strings.Builder
	for _, e := range log.Entries {
		sb.WriteString(fmt.Sprintf("[%s] %s %s",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Action,
			e.Key,
		))
		switch e.Action {
		case AuditChanged:
			sb.WriteString(fmt.Sprintf(" (%q -> %q)", e.OldValue, e.NewValue))
		case AuditAdded:
			sb.WriteString(fmt.Sprintf(" (value: %q)", e.NewValue))
		case AuditRemoved:
			sb.WriteString(fmt.Sprintf(" (was: %q)", e.OldValue))
		}
		sb.WriteString(fmt.Sprintf(" [source: %s]\n", e.Source))
	}
	return sb.String()
}

func formatAuditTable(log *AuditLog) string {
	if len(log.Entries) == 0 {
		return "No audit entries.\n"
	}
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIMESTAMP\tACTION\tKEY\tOLD\tNEW\tSOURCE")
	for _, e := range log.Entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			e.Timestamp.Format("2006-01-02T15:04:05Z"),
			e.Action, e.Key, e.OldValue, e.NewValue, e.Source,
		)
	}
	w.Flush()
	return buf.String()
}

func formatAuditJSON(log *AuditLog) string {
	type jsonEntry struct {
		Timestamp string `json:"timestamp"`
		Action    string `json:"action"`
		Key       string `json:"key"`
		OldValue  string `json:"old_value,omitempty"`
		NewValue  string `json:"new_value,omitempty"`
		Source    string `json:"source"`
	}
	out := make([]jsonEntry, 0, len(log.Entries))
	for _, e := range log.Entries {
		out = append(out, jsonEntry{
			Timestamp: e.Timestamp.Format("2006-01-02T15:04:05Z"),
			Action:    string(e.Action),
			Key:       e.Key,
			OldValue:  e.OldValue,
			NewValue:  e.NewValue,
			Source:    e.Source,
		})
	}
	b, _ := json.MarshalIndent(out, "", "  ")
	return string(b) + "\n"
}
