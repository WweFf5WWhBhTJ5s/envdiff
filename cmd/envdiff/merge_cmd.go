package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envdiff/internal/envfile"
)

func runMerge(args []string) error {
	fs := flag.NewFlagSet("merge", flag.ContinueOnError)
	strategyFlag := fs.String("strategy", "ours", "conflict resolution strategy: ours|theirs|error")
	outputFlag := fs.String("output", "", "output file path (default: stdout)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	remaining := fs.Args()
	if len(remaining) != 2 {
		return fmt.Errorf("usage: envdiff merge [flags] <base.env> <incoming.env>")
	}

	baseFile, err := envfile.Parse(remaining[0])
	if err != nil {
		return fmt.Errorf("parsing base file: %w", err)
	}

	incomingFile, err := envfile.Parse(remaining[1])
	if err != nil {
		return fmt.Errorf("parsing incoming file: %w", err)
	}

	var strategy envfile.MergeStrategy
	switch strings.ToLower(*strategyFlag) {
	case "ours":
		strategy = envfile.MergeStrategyOurs
	case "theirs":
		strategy = envfile.MergeStrategyTheirs
	case "error":
		strategy = envfile.MergeStrategyError
	default:
		return fmt.Errorf("unknown strategy %q: must be ours, theirs, or error", *strategyFlag)
	}

	result, err := envfile.Merge(baseFile, incomingFile, envfile.MergeOptions{Strategy: strategy})
	if err != nil {
		return fmt.Errorf("merge failed: %w", err)
	}

	if len(result.Conflicts) > 0 {
		fmt.Fprintf(os.Stderr, "merge conflicts resolved (%d):\n", len(result.Conflicts))
		for _, c := range result.Conflicts {
			fmt.Fprintf(os.Stderr, "  %s: %q vs %q → chose %q\n", c.Key, c.BaseVal, c.TheirVal, c.Chosen)
		}
	}

	var sb strings.Builder
	for _, entry := range result.File.Entries {
		if entry.Comment != "" {
			sb.WriteString("# " + entry.Comment + "\n")
		}
		sb.WriteString(entry.Key + "=" + entry.Value + "\n")
	}

	if *outputFlag == "" {
		fmt.Print(sb.String())
		return nil
	}

	if err := os.WriteFile(*outputFlag, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}
	fmt.Fprintf(os.Stderr, "merged file written to %s\n", *outputFlag)
	return nil
}
