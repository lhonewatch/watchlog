// Package output provides colorized, human-readable formatting for
// structured log entries produced by the parser package.
//
// # Colorizer
//
// NewColorizer returns a Colorizer that formats a parser.Entry into a
// single terminal line with ANSI colour codes.  Colour selection is
// driven by the log level so that INFO, WARN, ERROR and DEBUG each
// receive a distinct colour.  Any extra fields present in the entry are
// appended as key=value pairs after the core timestamp / level / message
// triple.
//
// # Formatter
//
// NewFormatter wraps a Colorizer and an io.Writer.  Call
// Formatter.Write to render an entry and emit it to the writer in one
// step.  This is the type consumed by pipeline.Pipeline.
package output
