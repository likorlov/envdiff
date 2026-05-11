package interpolator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// varPattern matches ${VAR_NAME} and $VAR_NAME style references.
var varPattern = regexp.MustCompile(`\$\{([A-Z_][A-Z0-9_]*)\}|\$([A-Z_][A-Z0-9_]*)`)

// Options controls interpolation behaviour.
type Options struct {
	// FallbackToOS allows falling back to the process environment when a
	// variable is not found in the provided env map.
	FallbackToOS bool
	// FailOnMissing returns an error when a referenced variable cannot be
	// resolved instead of leaving the placeholder in place.
	FailOnMissing bool
}

// Interpolate resolves variable references within the values of env.
// References in keys are not expanded. The original map is not modified.
func Interpolate(env map[string]string, opts Options) (map[string]string, error) {
	result := make(map[string]string, len(env))
	for k, v := range env {
		result[k] = v
	}

	for k, v := range result {
		expanded, err := expand(v, result, opts)
		if err != nil {
			return nil, fmt.Errorf("interpolating %q: %w", k, err)
		}
		result[k] = expanded
	}
	return result, nil
}

func expand(value string, env map[string]string, opts Options) (string, error) {
	var expandErr error
	result := varPattern.ReplaceAllStringFunc(value, func(match string) string {
		if expandErr != nil {
			return match
		}
		name := extractName(match)
		if resolved, ok := env[name]; ok {
			return resolved
		}
		if opts.FallbackToOS {
			if osVal, ok := os.LookupEnv(name); ok {
				return osVal
			}
		}
		if opts.FailOnMissing {
			expandErr = fmt.Errorf("variable %q is not defined", name)
			return match
		}
		return match
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}

func extractName(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}
