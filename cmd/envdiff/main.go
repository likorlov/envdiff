package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "envdiff",
		Short: "Diff and reconcile environment variable files",
	}

	root.AddCommand(
		newDiffCmd(),
		newValidateCmd(),
		newReconcileCmd(),
		newExportCmd(),
		newMergeCmd(),
		newLintCmd(),
	)

	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func newDiffCmd() *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "diff <base> <compare>",
		Short: "Show differences between two env files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDiff(args[0], args[1], format)
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "text", "output format: text, dotenv, json")
	return cmd
}

func newValidateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate an env file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runValidate(args[0])
		},
	}
}

func newReconcileCmd() *cobra.Command {
	var apply bool
	cmd := &cobra.Command{
		Use:   "reconcile <base> <target>",
		Short: "Reconcile base env into target",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReconcile(args[0], args[1], apply)
		},
	}
	cmd.Flags().BoolVar(&apply, "apply", false, "write changes to target file")
	return cmd
}
