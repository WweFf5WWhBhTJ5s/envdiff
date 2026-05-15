package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envdiff/internal/envfile"
)

func init() {
	var (
		label   string
		output  string
		format  string
		maskOn  bool
		compare string
	)

	cmd := &cobra.Command{
		Use:   "snapshot <file>",
		Short: "Take or compare a snapshot of an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSnapshot(args[0], label, output, format, maskOn, compare)
		},
	}

	cmd.Flags().StringVar(&label, "label", "", "Label for this snapshot")
	cmd.Flags().StringVar(&output, "out", "", "Save snapshot to this JSON file")
	cmd.Flags().StringVar(&format, "format", "text", "Output format: text|table|json")
	cmd.Flags().BoolVar(&maskOn, "mask", false, "Mask sensitive values in snapshot")
	cmd.Flags().StringVar(&compare, "compare", "", "Compare against a previously saved snapshot file")

	rootCmd.AddCommand(cmd)
}

func runSnapshot(src, label, output, format string, maskOn bool, comparePath string) error {
	ef, err := envfile.Parse(src)
	if err != nil {
		return fmt.Errorf("parse %s: %w", src, err)
	}

	opts := envfile.SnapshotOptions{
		Label:  label,
		Source: src,
	}
	if maskOn {
		opts.Masker = envfile.NewMasker()
	}

	snap := envfile.TakeSnapshot(ef, opts)

	if output != "" {
		if err := envfile.SaveSnapshot(snap, output); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "snapshot saved to %s\n", output)
	}

	if comparePath != "" {
		base, err := envfile.LoadSnapshot(comparePath)
		if err != nil {
			return fmt.Errorf("load snapshot %s: %w", comparePath, err)
		}
		results := envfile.DiffSnapshot(base, snap)
		out, err := envfile.FormatDiff(results, format)
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	}

	out, err := envfile.FormatSnapshot(snap, format)
	if err != nil {
		return err
	}
	fmt.Print(out)
	return nil
}
