// Package filter implements rule-based filtering for structured log entries.
//
// A Filter is composed of one or more Rules. Each Rule targets a specific
// field name and matches its value against a regular expression. A log entry
// must satisfy ALL rules (logical AND) to pass the filter.
//
// Example usage:
//
//	rules := []*filter.Rule{
//		filter.NewRule("level", "error|warn"),
//		filter.NewRule("service", "^api"),
//	}
//	f, err := filter.New(rules)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if f.Match(entry) {
//		// entry passes all filter rules
//	}
package filter
