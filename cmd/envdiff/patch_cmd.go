package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envdiff/internal/envfile"
)

func init() {
	var (
		errorOnMissing bool
		errorOnUnknown bool
		dryRun         bool
	)

	cmd := &cobra.Command{
		Use:   "patch <file> KEY=VALUE...",
		Short: "Apply key=value overrides to an env file",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPatch(args[0], args[1:], errorOnMissing, errorOnUnknown, dryRun)
		},
	}

	cmd.Flags().BoolVar(&errorOnMissing, "error-on-missing", false, "Error if a key does not exist in the file")
	cmd.Flags().BoolVar(&errorOnUnknown, "error-on-unknown", false, "Error if a patch key is not present in the file")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print changes without writing")

	rootCmd.AddCommand(cmd)
}

func runPatch(path string, pairs []string, errorOnMissing, errorOnUnknown, dryRun bool) error {
	ef, err := envfile.Parse(path)
	if err != nil {
		return fmt.Errorf("patch: parse %s: %w", path, err)
	}

	patches := make(map[string]string, len(pairs))
	for _, p := range pairs {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("patch: invalid pair %q (expected KEY=VALUE)", p)
		}
		patches[parts[0]] = parts[1]
	}

	out, results, err := envfile.Patch(ef, patches, envfile.PatchOptions{
		ErrorOnMissing: errorOnMissing,
		ErrorOnUnknown: errorOnUnknown,
	})
	if err != nil {
		return err
	}

	for _, r := range results {
		switch r.Action {
		case "set":
			fmt.Fprintf(os.Stderr, "~ %s: %q -> %q\n", r.Key, r.OldVal, r.NewVal)
		case "added":
			fmt.Fprintf(os.Stderr, "+ %s=%q\n", r.Key, r.NewVal)
		}
	}

	if dryRun {
		fmt.Fprintln(os.Stderr, "(dry-run) no changes written")
		return nil
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("patch: open %s for writing: %w", path, err)
	}
	defer f.Close()

	for _, e := range out.Entries {
		if _, err := fmt.Fprintf(f, "%s=%s\n", e.Key, e.Value); err != nil {
			return fmt.Errorf("patch: write: %w", err)
		}
	}
	return nil
}
