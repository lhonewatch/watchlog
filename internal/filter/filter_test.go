package filter_test

import (
	"testing"

	"github.com/yourorg/watchlog/internal/filter"
)

func makeEntry(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestFilter_EmptyRulesMatchesAll(t *testing.T) {
	f, err := filter.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entry := makeEntry("level", "info", "msg", "hello")
	if !f.Match(entry) {
		t.Error("empty filter should match any entry")
	}
}

func TestFilter_MatchSingleRule(t *testing.T) {
	f, err := filter.New([]*filter.Rule{filter.NewRule("level", "error")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Match(makeEntry("level", "error")) {
		t.Error("expected match for level=error")
	}
	if f.Match(makeEntry("level", "info")) {
		t.Error("expected no match for level=info")
	}
}

func TestFilter_MissingFieldNoMatch(t *testing.T) {
	f, err := filter.New([]*filter.Rule{filter.NewRule("service", "api")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Match(makeEntry("level", "info")) {
		t.Error("missing field should not match")
	}
}

func TestFilter_InvalidPatternReturnsError(t *testing.T) {
	_, err := filter.New([]*filter.Rule{filter.NewRule("level", "[invalid")})
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestMatchLevel_CaseInsensitive(t *testing.T) {
	entry := makeEntry("level", "WARN")
	if !filter.MatchLevel(entry, []string{"warn", "error"}) {
		t.Error("expected case-insensitive level match")
	}
}

func TestMatchLevel_EmptyLevelsMatchAll(t *testing.T) {
	entry := makeEntry("level", "debug")
	if !filter.MatchLevel(entry, nil) {
		t.Error("empty levels list should match all")
	}
}
