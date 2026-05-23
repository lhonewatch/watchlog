// Package pipeline wires together source, parser, filter, and output
// into a single processing pipeline for log entries.
package pipeline

import (
	"context"
	"io"

	"github.com/user/watchlog/internal/filter"
	"github.com/user/watchlog/internal/output"
	"github.com/user/watchlog/internal/parser"
)

// Parser converts a raw log line into a LogEntry.
type Parser interface {
	Parse(line string) (*parser.LogEntry, error)
}

// Filter decides whether a LogEntry should be forwarded.
type Filter interface {
	Match(entry *parser.LogEntry) (bool, error)
}

// Colorizer formats a LogEntry as a human-readable, colorized string.
type Colorizer interface {
	Format(entry *parser.LogEntry) string
}

// Pipeline reads lines from a channel, parses, filters, and writes them.
type Pipeline struct {
	parser    Parser
	filter    Filter
	colorizer Colorizer
	writer    io.Writer
}

// New creates a Pipeline with the provided components.
func New(p Parser, f Filter, c Colorizer, w io.Writer) *Pipeline {
	return &Pipeline{
		parser:    p,
		filter:    f,
		colorizer: c,
		writer:    w,
	}
}

// Run consumes lines from the given channel until it is closed or ctx is done.
// Unparseable lines are silently skipped; filter errors cause the entry to be
// dropped. Run blocks until the context is cancelled or lines is closed.
func (pl *Pipeline) Run(ctx context.Context, lines <-chan string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case line, ok := <-lines:
			if !ok {
				return nil
			}
			entry, err := pl.parser.Parse(line)
			if err != nil || entry == nil {
				continue
			}
			matched, err := pl.filter.Match(entry)
			if err != nil || !matched {
				continue
			}
			fmt.Fprintln(pl.writer, pl.colorizer.Format(entry))
		}
	}
}
