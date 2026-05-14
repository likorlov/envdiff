package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "envdiff",
		Short: "Diff and reconcile environment variable files across deployment stages",
	}

	root.AddCommand(
		newDiffCmd(),
		newValidateCmd(),
		newReconcileCmd(),
		newExportCmd(),
		newMergeCmd(),
		newLintCmd(),
		newInterpolateCmd(),
		newProfileCmd(),
		newSnapshotCmd(),
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newDiffCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "diff <base> <head>",
		Short: "Show differences between two env files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDiff(cmd, args, format)
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text, dotenv, json")
	return cmd
}

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <env-file>",
		Short: "Validate an env file against common rules",
		Args:  cobra.ExactArgs(1),
		RunE:  runValidate,
	}
}

func newReconcileCmd() *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "reconcile <source> <target>",
		Short: "Reconcile target env file to match source",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReconcile(cmd, args, dryRun)
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Print steps without applying")
	return cmd
}
