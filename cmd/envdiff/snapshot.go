package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/envdiff/internal/parser"
	"github.com/user/envdiff/internal/snapshotter"
)

func newSnapshotCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Save and compare environment snapshots",
	}
	cmd.AddCommand(newSnapshotSaveCmd(), newSnapshotDiffCmd())
	return cmd
}

func newSnapshotSaveCmd() *cobra.Command {
	var label string
	cmd := &cobra.Command{
		Use:   "save <env-file> <output-snapshot>",
		Short: "Save a snapshot of an env file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := parser.ParseFile(args[0])
			if err != nil {
				return fmt.Errorf("parse: %w", err)
			}
			if label == "" {
				label = time.Now().UTC().Format(time.RFC3339)
			}
			if err := snapshotter.Save(args[1], label, env); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Snapshot saved to %s (label: %s)\n", args[1], label)
			return nil
		},
	}
	cmd.Flags().StringVarP(&label, "label", "l", "", "Label for the snapshot (default: current timestamp)")
	return cmd
}

func newSnapshotDiffCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "diff <base-snapshot> <head-snapshot>",
		Short: "Compare two snapshots",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			base, err := snapshotter.Load(args[0])
			if err != nil {
				return err
			}
			head, err := snapshotter.Load(args[1])
			if err != nil {
				return err
			}
			added, removed, changed := snapshotter.Compare(base, head)
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Base: %s (%s)\n", base.Label, base.Timestamp.Format(time.RFC3339))
			fmt.Fprintf(out, "Head: %s (%s)\n\n", head.Label, head.Timestamp.Format(time.RFC3339))
			printKeys(out, "+ Added", added)
			printKeys(out, "- Removed", removed)
			printKeys(out, "~ Changed", changed)
			if len(added)+len(removed)+len(changed) == 0 {
				fmt.Fprintln(out, "No differences found.")
			}
			return nil
		},
	}
}

func printKeys(out *os.File, header string, keys []string) {
	if len(keys) == 0 {
		return
	}
	sort.Strings(keys)
	fmt.Fprintf(out, "%s:\n", header)
	for _, k := range keys {
		fmt.Fprintf(out, "  %s\n", k)
	}
}
