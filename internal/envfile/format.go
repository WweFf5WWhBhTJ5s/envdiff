package envfile

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

// OutputFormat controls how diff results are rendered.
type OutputFormat string

const (
	FormatText  OutputFormat = "text"
	FormatTable OutputFormat = "table"
	FormatDotenv OutputFormat = "dotenv"
)

// FormatDiff renders a slice of DiffResult to the given writer in the specified format.
func FormatDiff(w io.Writer, results []DiffResult, format OutputFormat) error {
	switch format {
	case FormatTable:
		return formatTable(w, results)
	case FormatDotenv:
		return formatDotenv(w, results)
	default:
		return formatText(w, results)
	}
}

func formatText(w io.Writer, results []DiffResult) error {
	for _, r := range results {
		var line string
		switch r.Status {
		case StatusAdded:
			line = fmt.Sprintf("+ %s=%s", r.Key, r.NewValue)
		case StatusRemoved:
			line = fmt.Sprintf("- %s=%s", r.Key, r.OldValue)
		case StatusChanged:
			line = fmt.Sprintf("~ %s: %s -> %s", r.Key, r.OldValue, r.NewValue)
		case StatusUnchanged:
			line = fmt.Sprintf("  %s=%s", r.Key, r.NewValue)
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func formatTable(w io.Writer, results []DiffResult) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "STATUS\tKEY\tOLD VALUE\tNEW VALUE")
	fmt.Fprintln(tw, strings.Repeat("-", 6)+"\t"+strings.Repeat("-", 20)+"\t"+strings.Repeat("-", 20)+"\t"+strings.Repeat("-", 20))
	for _, r := range results {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.Status, r.Key, r.OldValue, r.NewValue)
	}
	return tw.Flush()
}

func formatDotenv(w io.Writer, results []DiffResult) error {
	for _, r := range results {
		if r.Status == StatusRemoved {
			continue
		}
		val := r.NewValue
		if r.Status == StatusUnchanged {
			val = r.OldValue
			if r.NewValue != "" {
				val = r.NewValue
			}
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", r.Key, val); err != nil {
			return err
		}
	}
	return nil
}
