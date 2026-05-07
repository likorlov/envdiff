package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/envdiff/internal/merger"
	"github.com/yourorg/envdiff/internal/parser"
)

func newMergeCmd() *cobra.Command {
	var strategyFlag string

	cmd := &cobra.Command{
		Use:   "merge <base> <incoming>",
		Short: "Merge two env files, resolving conflicts by strategy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMerge(args[0], args[1], strategyFlag)
		},
	}

	cmd.Flags().StringVarP(&strategyFlag, "strategy", "s", "ours",
		"Conflict resolution strategy: ours | theirs | union")
	return cmd
}

func runMerge(basePath, incomingPath, strategyFlag string) error {
	base, err := parser.ParseFile(basePath)
	if err != nil {
		return fmt.Errorf("reading base file: %w", err)
	}
	incoming, err := parser.ParseFile(incomingPath)
	if err != nil {
		return fmt.Errorf("reading incoming file: %w", err)
	}

	var strategy merger.Strategy
	switch strategyFlag {
	case "ours":
		strategy = merger.StrategyOurs
	case "theirs":
		strategy = merger.StrategyTheirs
	case "union":
		strategy = merger.StrategyUnion
	default:
		return fmt.Errorf("unknown strategy %q: choose ours, theirs, or union", strategyFlag)
	}

	res := merger.Merge(base, incoming, strategy)

	if len(res.Conflicts) > 0 {
		fmt.Fprintln(os.Stderr, "# Conflicts resolved:")
		for _, c := range res.Conflicts {
			fmt.Fprintf(os.Stderr, "#   %s: base=%q theirs=%q -> resolved=%q\n",
				c.Key, c.BaseVal, c.TheirVal, c.Resolved)
		}
	}

	for _, k := range sortedKeys(res.Env) {
		fmt.Printf("%s=%s\n", k, res.Env[k])
	}
	return nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
