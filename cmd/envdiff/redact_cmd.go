package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourorg/envdiff/internal/envfile"
)

func init() {
	var mode string
	var maskChar string
	var maskLen int
	var outputPath string

	cmd := &cobra.Command{
		Use:   "redact <file>",
		Short: "Redact sensitive values in an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRedact(args[0], mode, maskChar, maskLen, outputPath)
		},
	}

	cmd.Flags().StringVarP(&mode, "mode", "m", "mask", "Redaction mode: mask, hash, blank, length")
	cmd.Flags().StringVar(&maskChar, "mask-char", "*", "Character to use for mask mode")
	cmd.Flags().IntVar(&maskLen, "mask-len", 6, "Fixed mask length (0 = use actual value length)")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Write redacted output to file (default: stdout)")

	rootCmd.AddCommand(cmd)
}

func runRedact(filePath, mode, maskChar string, maskLen int, outputPath string) error {
	ef, err := envfile.Parse(filePath)
	if err != nil {
		return fmt.Errorf("parse error: %w", err)
	}

	opts := envfile.RedactOptions{
		Mode:     envfile.RedactMode(strings.ToLower(mode)),
		MaskChar: maskChar,
		MaskLen:  maskLen,
	}

	redactor := envfile.NewRedactor(opts)
	redacted := redactor.RedactEnvFile(ef)

	var sb strings.Builder
	for _, e := range redacted.Entries {
		if e.Value == "" {
			fmt.Fprintf(&sb, "%s=\n", e.Key)
		} else {
			fmt.Fprintf(&sb, "%s=%s\n", e.Key, e.Value)
		}
	}

	if outputPath != "" {
		if err := os.WriteFile(outputPath, []byte(sb.String()), 0o644); err != nil {
			return fmt.Errorf("write error: %w", err)
		}
		fmt.Fprintf(os.Stderr, "redacted output written to %s\n", outputPath)
		return nil
	}

	fmt.Print(sb.String())
	return nil
}
