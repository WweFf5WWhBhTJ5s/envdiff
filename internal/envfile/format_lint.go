package envfile

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"bytes"
)

// FormatLintResult formats lint issues for display.
func FormatLintResult(result LintResult, format string) string {
	switch format {
	case "table":
		return formatLintTable(result)
	case "json":
		return formatLintJSON(result)
	default:
		return formatLintText(result)
	}
}

func formatLintText(result LintResult) string {
	if len(result.Issues) == 0 {
		return "No lint issues found.\n"
	}
	var sb strings.Builder
	for _, issue := range result.Issues {
		icon := "⚠"
		if issue.Severity == LintError {
			icon = "✖"
		}
		sb.WriteString(fmt.Sprintf("%s %s\n", icon, issue.String()))
	}
	return sb.String()
}

func formatLintTable(result LintResult) string {
	if len(result.Issues) == 0 {
		return "No lint issues found.\n"
	}
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "SEVERITY\tLINE\tKEY\tMESSAGE")
	fmt.Fprintln(w, "--------\t----\t---\t-------")
	for _, issue := range result.Issues {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\n",
			strings.ToUpper(string(issue.Severity)),
			issue.Line,
			issue.Key,
			issue.Message,
		)
	}
	w.Flush()
	return buf.String()
}

func formatLintJSON(result LintResult) string {
	if len(result.Issues) == 0 {
		return `{"issues":[]}` + "\n"
	}
	var sb strings.Builder
	sb.WriteString(`{"issues":[`)
	for i, issue := range result.Issues {
		if i > 0 {
			sb.WriteString(",")
		}
		sb.WriteString(fmt.Sprintf(
			`{"severity":%q,"line":%d,"key":%q,"message":%q}`,
			issue.Severity, issue.Line, issue.Key, issue.Message,
		))
	}
	sb.WriteString(`]}`)
	sb.WriteString("\n")
	return sb.String()
}
