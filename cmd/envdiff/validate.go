package main

import (
	"fmt"
	"os"

	"github.com/yourorg/envdiff/internal/parser"
	"github.com/yourorg/envdiff/internal/validator"
)

// runValidate parses the given env file and validates all keys/values
// against the default rule set. Exits with code 1 if violations are found.
// If strict is true, an additional rule requiring non-empty values is applied.
func runValidate(filePath string, strict bool) error {
	env, err := parser.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse %q: %w", filePath, err)
	}

	rules := validator.DefaultRules
	if strict {
		rules = append(rules, validator.NoEmptyValues)
	}

	violations := validator.Validate(env, rules)

	if len(violations) == 0 {
		fmt.Printf("✓ %s passed validation (%d keys checked)\n", filePath, len(env))
		return nil
	}

	return reportViolations(filePath, violations)
}

// reportViolations prints all validation violations to stderr and exits with
// code 1. It returns nil to satisfy the error return type of the caller,
// since the process will have already exited.
func reportViolations(filePath string, violations []validator.Violation) error {
	fmt.Fprintf(os.Stderr, "✗ %s has %d validation issue(s):\n", filePath, len(violations))
	fmt.Fprintln(os.Stderr, validator.FormatViolations(violations))
	os.Exit(1)
	return nil
}
