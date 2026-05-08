package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envdiff/internal/envfile"
)

var (
	format  string
	maskOn  bool
	dryRun  bool
	syncFull bool
)

var rootCmd = &cobra.Command{
	Use:   "envdiff",
	Short: "Diff and sync .env files across environments",
}

var diffCmd = &cobra.Command{
	Use:   "diff <source> <target>",
	Short: "Show differences between two .env files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		src, err := envfile.Parse(args[0])
		if err != nil {
			return fmt.Errorf("parsing source: %w", err)
		}
		tgt, err := envfile.Parse(args[1])
		if err != nil {
			return fmt.Errorf("parsing target: %w", err)
		}

		results := envfile.Diff(src, tgt)

		if maskOn {
			masker := envfile.NewMasker()
			results = masker.MaskDiffResults(results)
		}

		return envfile.FormatDiff(os.Stdout, results, envfile.OutputFormat(format))
	},
}

var syncCmd = &cobra.Command{
	Use:   "sync <source> <target>",
	Short: "Sync keys from source .env into target .env",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		src, err := envfile.Parse(args[0])
		if err != nil {
			return fmt.Errorf("parsing source: %w", err)
		}
		tgt, err := envfile.Parse(args[1])
		if err != nil {
			return fmt.Errorf("parsing target: %w", err)
		}

		mode := envfile.SyncAddMissing
		if syncFull {
			mode = envfile.SyncFull
		}

		result, err := envfile.Sync(src, tgt, args[1], envfile.SyncOptions{
			Mode:   mode,
			DryRun: dryRun,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Added: %d, Updated: %d, Removed: %d\n",
			len(result.Added), len(result.Updated), len(result.Removed))
		if dryRun {
			fmt.Println("(dry run — no changes written)")
		}
		return nil
	},
}

func init() {
	diffCmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, table, dotenv")
	diffCmd.Flags().BoolVar(&maskOn, "mask", false, "Mask sensitive values")

	syncCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview changes without writing")
	syncCmd.Flags().BoolVar(&syncFull, "full", false, "Remove keys not present in source")

	rootCmd.AddCommand(diffCmd, syncCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
