package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"

	"github.com/user/envdiff/internal/envfile"
)

type schemaFileSpec struct {
	Rules []struct {
		Key      string   `json:"key"`
		Required bool     `json:"required"`
		Pattern  string   `json:"pattern,omitempty"`
		Allowed  []string `json:"allowed,omitempty"`
	} `json:"rules"`
}

func init() {
	var schemaFile string

	schemaCmd := &cobra.Command{
		Use:   "schema [env-file]",
		Short: "Validate an .env file against a JSON schema definition",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSchema(args[0], schemaFile)
		},
	}

	schemaCmd.Flags().StringVarP(&schemaFile, "schema", "s", ".envschema.json", "path to JSON schema file")
	rootCmd.AddCommand(schemaCmd)
}

func runSchema(envPath, schemaPath string) error {
	ef, err := envfile.Parse(envPath)
	if err != nil {
		return fmt.Errorf("parsing env file: %w", err)
	}

	raw, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("reading schema file: %w", err)
	}

	var spec schemaFileSpec
	if err := json.Unmarshal(raw, &spec); err != nil {
		return fmt.Errorf("parsing schema JSON: %w", err)
	}

	schema := envfile.Schema{}
	for _, r := range spec.Rules {
		rule := envfile.SchemaRule{
			Key:     r.Key,
			Required: r.Required,
			Allowed: r.Allowed,
		}
		if r.Pattern != "" {
			compiled, err := regexp.Compile(r.Pattern)
			if err != nil {
				return fmt.Errorf("invalid pattern for key %q: %w", r.Key, err)
			}
			rule.Pattern = compiled
		}
		schema.Rules = append(schema.Rules, rule)
	}

	result := envfile.ValidateSchema(ef, schema)
	if result.Valid {
		fmt.Println("schema validation passed")
		return nil
	}

	fmt.Fprintf(os.Stderr, "schema validation failed (%d violation(s)):\n", len(result.Violations))
	for _, v := range result.Violations {
		fmt.Fprintf(os.Stderr, "  [%s] %s\n", v.Key, v.Message)
	}
	os.Exit(1)
	return nil
}
