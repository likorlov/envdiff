package main

import (
	"fmt"
	"os"

	"github.com/yourorg/envdiff/internal/exporter"
	"github.com/yourorg/envdiff/internal/parser"
	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "export <file>",
		Short: "Export an env file in a different format",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExport(args[0], exporter.Format(format))
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "dotenv", "Output format: dotenv, json, shell")
	return cmd
}

func runExport(filePath string, format exporter.Format) error {
	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", filePath, err)
	}

	out, err := exporter.Export(env, format)
	if err != nil {
		return fmt.Errorf("exporting: %w", err)
	}

	_, err = fmt.Fprint(os.Stdout, out)
	return err
}
