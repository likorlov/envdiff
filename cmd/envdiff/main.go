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

	root.AddCommand(newDiffCmd())
	root.AddCommand(newValidateCmd())
	root.AddCommand(newReconcileCmd())
	root.AddCommand(newExportCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newDiffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff <base> <target>",
		Short: "Show differences between two env files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("diff command not yet wired")
		},
	}
}

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate an env file against naming conventions",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate(args[0])
		},
	}
}

func newReconcileCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "reconcile <base> <target>",
		Short: "Print reconciliation steps to align base with target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("reconcile command not yet wired")
		},
	},
}
