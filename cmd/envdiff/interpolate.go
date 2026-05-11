package main

import (
	"fmt"
	"os"

	"github.com/user/envdiff/internal/exporter"
	"github.com/user/envdiff/internal/interpolator"
	"github.com/user/envdiff/internal/parser"
	"github.com/spf13/cobra"
)

func newInterpolateCmd() *cobra.Command {
	var (
		outputFormat  string
		fallbackToOS  bool
		failOnMissing bool
	)

	cmd := &cobra.Command{
		Use:   "interpolate <file>",
		Short: "Expand variable references within an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInterpolate(args[0], outputFormat, fallbackToOS, failOnMissing)
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "format", "f", "dotenv", "Output format: dotenv, json, shell")
	cmd.Flags().BoolVar(&fallbackToOS, "fallback-os", false, "Fall back to process environment for unresolved variables")
	cmd.Flags().BoolVar(&failOnMissing, "fail-missing", false, "Exit with error when a referenced variable cannot be resolved")

	return cmd
}

func runInterpolate(file, format string, fallbackToOS, failOnMissing bool) error {
	env, err := parser.ParseFile(file)
	if err != nil {
		return fmt.Errorf("parsing %q: %w", file, err)
	}

	opts := interpolator.Options{
		FallbackToOS:  fallbackToOS,
		FailOnMissing: failOnMissing,
	}

	resolved, err := interpolator.Interpolate(env, opts)
	if err != nil {
		return fmt.Errorf("interpolation failed: %w", err)
	}

	out, err := exporter.Export(resolved, format)
	if err != nil {
		return fmt.Errorf("formatting output: %w", err)
	}

	fmt.Fprint(os.Stdout, out)
	return nil
}
