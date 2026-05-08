package linter

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule represents a linting rule with a name and description.
type Rule string

const (
	RuleDuplicateKey    Rule = "duplicate-key"
	RuleQuotingStyle    Rule = "quoting-style"
	RuleTrailingSpace   Rule = "trailing-space"
	RuleNoValue         Rule = "no-value"
	RuleKeyNamingConvention Rule = "key-naming-convention"
)

// Violation describes a single linting issue.
type Violation struct {
	Line    int
	Key     string
	Rule    Rule
	Message string
}

func (v Violation) String() string {
	return fmt.Sprintf("line %d [%s] %s: %s", v.Line, v.Rule, v.Key, v.Message)
}

var validKeyRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Lint analyses raw lines of an env file and returns any violations found.
func Lint(lines []string) []Violation {
	var violations []Violation
	seen := map[string]int{}

	for i, raw := range lines {
		lineNum := i + 1
		trimmed := strings.TrimRight(raw, " \t")

		if trimmed != raw {
			violations = append(violations, Violation{
				Line:    lineNum,
				Rule:    RuleTrailingSpace,
				Message: "trailing whitespace detected",
			})
		}

		if strings.HasPrefix(strings.TrimSpace(raw), "#") || strings.TrimSpace(raw) == "" {
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) < 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := parts[1]

		if prev, ok := seen[key]; ok {
			violations = append(violations, Violation{
				Line:    lineNum,
				Key:     key,
				Rule:    RuleDuplicateKey,
				Message: fmt.Sprintf("duplicate of key first seen on line %d", prev),
			})
		} else {
			seen[key] = lineNum
		}

		if !validKeyRe.MatchString(key) {
			violations = append(violations, Violation{
				Line:    lineNum,
				Key:     key,
				Rule:    RuleKeyNamingConvention,
				Message: "key must match ^[A-Z][A-Z0-9_]*$",
			})
		}

		if strings.TrimSpace(val) == "" {
			violations = append(violations, Violation{
				Line:    lineNum,
				Key:     key,
				Rule:    RuleNoValue,
				Message: "key has no value",
			})
		}

		if (strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) ||
			(strings.HasPrefix(val, `"`) && strings.HasSuffix(val, `"`)) {
			if !strings.HasPrefix(val, `"`) {
				violations = append(violations, Violation{
					Line:    lineNum,
					Key:     key,
					Rule:    RuleQuotingStyle,
					Message: "prefer double quotes over single quotes",
				})
			}
		}
	}

	return violations
}

// FormatViolations returns a human-readable summary of lint violations.
func FormatViolations(violations []Violation) string {
	if len(violations) == 0 {
		return "no lint violations found"
	}
	var sb strings.Builder
	for _, v := range violations {
		sb.WriteString(v.String())
		sb.WriteByte('\n')
	}
	return sb.String()
}
