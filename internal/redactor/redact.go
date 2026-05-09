package redactor

import (
	"regexp"
	"strings"
)

// DefaultSensitivePatterns are common patterns for sensitive environment variable keys.
var DefaultSensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(password|passwd|pwd)`),
	regexp.MustCompile(`(?i)(secret|token|api_key|apikey)`),
	regexp.MustCompile(`(?i)(private_key|privkey)`),
	regexp.MustCompile(`(?i)(auth|credential|cred)`),
	regexp.MustCompile(`(?i)(dsn|database_url|db_url)`),
}

const redactedPlaceholder = "***REDACTED***"

// Options configures the redaction behaviour.
type Options struct {
	// ExtraPatterns are additional regexp patterns to match sensitive keys.
	ExtraPatterns []*regexp.Regexp
	// Placeholder overrides the default redaction string.
	Placeholder string
}

// Redact returns a copy of env with sensitive values replaced by a placeholder.
// Keys are matched case-insensitively against DefaultSensitivePatterns and any
// ExtraPatterns supplied via opts.
func Redact(env map[string]string, opts Options) map[string]string {
	placeholder := redactedPlaceholder
	if opts.Placeholder != "" {
		placeholder = opts.Placeholder
	}

	patterns := append(DefaultSensitivePatterns, opts.ExtraPatterns...)

	result := make(map[string]string, len(env))
	for k, v := range env {
		if isSensitive(k, patterns) {
			result[k] = placeholder
		} else {
			result[k] = v
		}
	}
	return result
}

// IsSensitiveKey reports whether the given key matches any sensitive pattern.
func IsSensitiveKey(key string) bool {
	return isSensitive(key, DefaultSensitivePatterns)
}

func isSensitive(key string, patterns []*regexp.Regexp) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if p.MatchString(upper) {
			return true
		}
	}
	return false
}
