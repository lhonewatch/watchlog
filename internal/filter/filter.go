// Package filter provides log entry filtering based on field matching rules.
package filter

import (
	"regexp"
	"strings"
)

// Rule represents a single filter condition applied to a log field.
type Rule struct {
	Field   string
	Pattern string
	regexp  *regexp.Regexp
}

// Filter holds a collection of rules and applies them to log entries.
type Filter struct {
	rules []*Rule
}

// New creates a Filter with the provided rules. Returns an error if any
// pattern fails to compile as a regular expression.
func New(rules []*Rule) (*Filter, error) {
	for _, r := range rules {
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, err
		}
		r.regexp = re
	}
	return &Filter{rules: rules}, nil
}

// NewRule constructs a Rule for the given field and pattern.
func NewRule(field, pattern string) *Rule {
	return &Rule{Field: field, Pattern: pattern}
}

// Match reports whether the log entry (represented as a map of fields)
// satisfies ALL rules in the filter. An empty filter matches everything.
func (f *Filter) Match(entry map[string]string) bool {
	for _, r := range f.rules {
		val, ok := entry[r.Field]
		if !ok {
			return false
		}
		if !r.regexp.MatchString(val) {
			return false
		}
	}
	return true
}

// MatchLevel is a convenience helper that checks whether the entry's level
// field (case-insensitive) matches one of the accepted levels.
func MatchLevel(entry map[string]string, levels []string) bool {
	if len(levels) == 0 {
		return true
	}
	entryLevel := strings.ToLower(entry["level"])
	for _, l := range levels {
		if strings.ToLower(l) == entryLevel {
			return true
		}
	}
	return false
}
