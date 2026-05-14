package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"envdiff/internal/parser"
	"envdiff/internal/profiler"
)

func newProfileCmd() *cobra.Command {
	var required []string
	var optional []string
	var name string
	var jsonOut bool

	cmd := &cobra.Command{
		Use:   "profile <file>",
		Short: "Check an env file against a named profile of required/optional keys",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runProfile(args[0], name, required, optional, jsonOut)
		},
	}

	cmd.Flags().StringVar(&name, "name", "default", "profile name")
	cmd.Flags().StringSliceVar(&required, "require", nil, "comma-separated required keys")
	cmd.Flags().StringSliceVar(&optional, "optional", nil, "comma-separated optional keys")
	cmd.Flags().BoolVar(&jsonOut, "json", false, "output results as JSON")

	return cmd
}

func runProfile(file, name string, required, optional []string, jsonOut bool) error {
	env, err := parser.ParseFile(file)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", file, err)
	}

	p := profiler.Profile{
		Name:     name,
		Required: required,
		Optional: optional,
	}

	violations := profiler.Check(env, p)
	coverage := profiler.Coverage(env, p)

	if jsonOut {
		type result struct {
			Profile    string               `json:"profile"`
			Coverage   float64              `json:"optional_coverage_pct"`
			Violations []profiler.Violation `json:"violations"`
		}
		r := result{
			Profile:    name,
			Coverage:   coverage,
			Violations: violations,
		}
		if r.Violations == nil {
			r.Violations = []profiler.Violation{}
		}
		return json.NewEncoder(os.Stdout).Encode(r)
	}

	fmt.Fprintf(os.Stdout, "Profile : %s\n", name)
	fmt.Fprintf(os.Stdout, "File    : %s\n", file)
	fmt.Fprintf(os.Stdout, "Coverage: %.1f%% of optional keys present\n", coverage)

	if len(violations) == 0 {
		fmt.Fprintln(os.Stdout, "Status  : OK")
		return nil
	}

	fmt.Fprintln(os.Stdout, "Status  : FAIL")
	fmt.Fprintln(os.Stdout, "Violations:")
	fmt.Fprintln(os.Stdout, profiler.FormatViolations(violations))

	keys := make([]string, 0, len(violations))
	for _, v := range violations {
		keys = append(keys, v.Key)
	}
	return fmt.Errorf("profile %q failed: missing or empty required keys: %s", name, strings.Join(keys, ", "))
}
