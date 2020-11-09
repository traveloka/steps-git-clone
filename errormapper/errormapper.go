package errormapper

import (
	"regexp"

	"github.com/bitrise-io/bitrise-init/step"
)

const (
	// UnknownParam ...
	UnknownParam = "::unknown::"
	// DetailedErrorRecKey ...
	DetailedErrorRecKey = "DetailedError"
)

// DetailedError ...
type DetailedError struct {
	Title       string
	Description string
}

// NewDetailedErrorRecommendation ...
func NewDetailedErrorRecommendation(detailedError DetailedError) step.Recommendation {
	return step.Recommendation{
		DetailedErrorRecKey: detailedError,
	}
}

// DetailedErrorBuilder ...
type DetailedErrorBuilder = func(...string) DetailedError

// GetParamAt ...
func GetParamAt(index int, params []string) string {
	res := UnknownParam
	if index >= 0 && len(params) > index {
		res = params[index]
	}
	return res
}

// PatternToDetailedErrorBuilder ...
type PatternToDetailedErrorBuilder map[string]DetailedErrorBuilder

// PatternErrorMatcher ...
type PatternErrorMatcher struct {
	DefaultBuilder   DetailedErrorBuilder
	PatternToBuilder PatternToDetailedErrorBuilder
}

// Run ...
func (m *PatternErrorMatcher) Run(msg string) step.Recommendation {
	for pattern, builder := range m.PatternToBuilder {
		re := regexp.MustCompile(pattern)
		if re.MatchString(msg) {
			// [search_string, match1, match2, ...]
			matches := re.FindStringSubmatch((msg))
			// Drop the first item, which is always the search_string itself
			// [search_string] -> []
			// [search_string, match1, ...] -> [match1, ...]
			params := matches[1:]
			detail := builder(params...)
			return NewDetailedErrorRecommendation(detail)
		}
	}

	detail := m.DefaultBuilder(msg)
	return NewDetailedErrorRecommendation(detail)
}