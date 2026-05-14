package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envdiff/internal/envfile"
)

var lintFormat string

var lintCmd = &cobra.Command{
	Use:   "lint <file>",
	Short: "Lint a .env file for style and correctness issues",
	Args:  cobra.ExactArgs(1),
	RunE:  runLint,
}

func init() {
	lintCmd.Flags().StringVarP(&lintFormat, "format", "f", "text", "Output format: text, table, json")
	rootCmd.AddCommand(lintCmd)
}

func runLint(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	ef, err := envfile.Parse(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	// Also run structural validation before linting
	validationIssues := envfile.Validate(ef)
	if len(validationIssues) > 0 {
		fmt.Fprintln(os.Stderr, "Validation errors found before linting:")
		for _, vi := range validationIssues {
			fmt.Fprintf(os.Stderr, "  - %s\n", vi)
		}
		return fmt.Errorf("fix validation errors before linting")
	}

	result := envfile.Lint(ef)
	fmt.Print(envfile.FormatLintResult(result, lintFormat))

	if result.HasErrors() {
		os.Exit(1)
	}
	return nil
}
