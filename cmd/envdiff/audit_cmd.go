package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"envdiff/internal/envfile"
)

func init() {
	var auditFormat string
	var maskSecrets bool
	var source string

	auditCmd := &cobra.Command{
		Use:   "audit <base> <target>",
		Short: "Generate an audit log of changes between two .env files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAudit(args[0], args[1], auditFormat, source, maskSecrets)
		},
	}

	auditCmd.Flags().StringVarP(&auditFormat, "format", "f", "text", "Output format: text, table, json")
	auditCmd.Flags().BoolVar(&maskSecrets, "mask", true, "Mask sensitive values in audit output")
	auditCmd.Flags().StringVar(&source, "source", "", "Label identifying the source of the change (e.g. 'ci', 'deploy-v1.2')")

	rootCmd.AddCommand(auditCmd)
}

func runAudit(basePath, targetPath, format, source string, maskSecrets bool) error {
	baseFile, err := envfile.Parse(basePath)
	if err != nil {
		return fmt.Errorf("parsing base file: %w", err)
	}

	targetFile, err := envfile.Parse(targetPath)
	if err != nil {
		return fmt.Errorf("parsing target file: %w", err)
	}

	results := envfile.Diff(baseFile, targetFile)

	if source == "" {
		source = fmt.Sprintf("%s -> %s", basePath, targetPath)
	}

	var masker *envfile.Masker
	if maskSecrets {
		masker = envfile.NewMasker()
	}

	log := envfile.Audit(results, source, masker)
	out := envfile.FormatAuditLog(log, format)
	fmt.Fprint(os.Stdout, out)
	return nil
}
