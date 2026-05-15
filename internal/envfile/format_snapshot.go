package envfile

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"
)

// FormatSnapshot renders a Snapshot in the requested format.
func FormatSnapshot(snap Snapshot, format string) (string, error) {
	switch strings.ToLower(format) {
	case "table":
		return formatSnapshotTable(snap), nil
	case "json":
		return formatSnapshotJSON(snap)
	default:
		return formatSnapshotText(snap), nil
	}
}

func formatSnapshotText(snap Snapshot) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Snapshot: %s\n", snap.Label))
	sb.WriteString(fmt.Sprintf("# Source:   %s\n", snap.Source))
	sb.WriteString(fmt.Sprintf("# Time:     %s\n", snap.Timestamp.Format("2006-01-02T15:04:05Z")))
	sb.WriteString("\n")
	keys := sortedSnapshotKeys(snap)
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, snap.Entries[k]))
	}
	return sb.String()
}

func formatSnapshotTable(snap Snapshot) string {
	var sb strings.Builder
	w := tabwriter.NewWriter(&sb, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "KEY\tVALUE\n")
	fmt.Fprintf(w, "---\t-----\n")
	keys := sortedSnapshotKeys(snap)
	for _, k := range keys {
		fmt.Fprintf(w, "%s\t%s\n", k, snap.Entries[k])
	}
	w.Flush()
	return sb.String()
}

func formatSnapshotJSON(snap Snapshot) (string, error) {
	b, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return "", fmt.Errorf("format snapshot json: %w", err)
	}
	return string(b) + "\n", nil
}

func sortedSnapshotKeys(snap Snapshot) []string {
	keys := make([]string, 0, len(snap.Entries))
	for k := range snap.Entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
