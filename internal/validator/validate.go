package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for environment variable keys or values.
type Rule struct {
	Name    string
	Pattern *regexp.Regexp
	Message string
}

// Violation represents a single validation failure.
type Violation struct {
	Key     string
	Value   string
	Rule    string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("key %q: %s", v.Key, v.Message)
}

var (
	// KeyFormat requires keys to be uppercase alphanumeric with underscores.
	KeyFormat = Rule{
		Name:    "key_format",
		Pattern: regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`),
		Message: "key must be uppercase alphanumeric with underscores",
	}

	// NoEmptyValues disallows blank values.
	NoEmptyValues = Rule{
		Name:    "no_empty_values",
		Pattern: regexp.MustCompile(`^.+$`),
		Message: "value must not be empty",
	}

	// NoWhitespaceKeys disallows whitespace in keys.
	NoWhitespaceKeys = Rule{
		Name:    "no_whitespace_keys",
		Pattern: regexp.MustCompile(`^\S+$`),
		Message: "key must not contain whitespace",
	}
)

// DefaultRules is the standard set of validation rules.
var DefaultRules = []Rule{KeyFormat, NoWhitespaceKeys}

// Validate checks all entries in the env map against the provided rules.
// It returns a slice of Violations (empty if all pass).
func Validate(env map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for key, value := range env {
		for _, rule := range rules {
			var subject string
			if strings.HasPrefix(rule.Name, "key_") || strings.HasPrefix(rule.Name, "no_whitespace") {
				subject = key
			} else {
				subject = value
			}
			if !rule.Pattern.MatchString(subject) {
				violations = append(violations, Violation{
					Key:     key,
					Value:   value,
					Rule:    rule.Name,
					Message: rule.Message,
				})
			}
		}
	}

	return violations
}

// FormatViolations returns a human-readable summary of violations.
func FormatViolations(violations []Violation) string {
	if len(violations) == 0 {
		return "no validation issues found"
	}
	var sb strings.Builder
	for _, v := range violations {
		fmt.Fprintf(&sb, "[%s] %s\n", v.Rule, v.Error())
	}
	return strings.TrimRight(sb.String(), "\n")
}
