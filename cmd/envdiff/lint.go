package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/user/envdiff/internal/linter"
	"github.com/spf13/cobra"
)

func newLintCmd() *cobra.Command {
	var strict bool

	cmd := &cobra.Command{
		Use:   "lint <file>",
		Short: "Lint an env file for common issues",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLint(args[0], strict)
		},
	}

	cmd.Flags().BoolVar(&strict, "strict", false, "exit non-zero even for warnings")
	return cmd
}

func runLint(path string, strict bool) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	violations := linter.Lint(lines)

	if len(violations) == 0 {
		fmt.Fprintln(os.Stdout, "no lint violations found")
		return nil
	}

	var sb strings.Builder
	for _, v := range violations {
		sb.WriteString(v.String())
		sb.WriteByte('\n')
	}
	fmt.Fprint(os.Stderr, sb.String())

	if strict || hasSevereViolations(violations) {
		return fmt.Errorf("%d lint violation(s) found", len(violations))
	}
	return nil
}

func hasSevereViolations(vs []linter.Violation) bool {
	for _, v := range vs {
		if v.Rule == linter.RuleDuplicateKey || v.Rule == linter.RuleKeyNamingConvention {
			return true
		}
	}
	return false
}
