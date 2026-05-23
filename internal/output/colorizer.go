// Package output provides colorized, formatted rendering of log entries.
package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"

	"watchlog/internal/parser"
)

// levelColors maps log level strings to color functions.
var levelColors = map[string]*color.Color{
	"debug": color.New(color.FgCyan),
	"info":  color.New(color.FgGreen),
	"warn":  color.New(color.FgYellow),
	"warning": color.New(color.FgYellow),
	"error": color.New(color.FgRed),
	"fatal": color.New(color.FgRed, color.Bold),
}

var (
	timestampColor = color.New(color.FgHiBlack)
	keyColor       = color.New(color.FgBlue)
	messageColor   = color.New(color.FgWhite)
	defaultLevel   = color.New(color.FgWhite)
)

// Colorizer writes formatted, colorized log entries to a writer.
type Colorizer struct {
	out        io.Writer
	extraFields []string
	noColor    bool
}

// NewColorizer creates a Colorizer that writes to out.
// extraFields lists additional JSON fields to display after the message.
func NewColorizer(out io.Writer, extraFields []string, noColor bool) *Colorizer {
	if noColor {
		color.NoColor = true
	}
	return &Colorizer{out: out, extraFields: extraFields, noColor: noColor}
}

// Write formats and writes a single log entry.
func (c *Colorizer) Write(entry *parser.Entry) {
	level := strings.ToLower(entry.Level)
	lvlColor, ok := levelColors[level]
	if !ok {
		lvlColor = defaultLevel
	}

	ts := timestampColor.Sprint(entry.Timestamp)
	lvl := lvlColor.Sprintf("%-5s", strings.ToUpper(level))
	msg := messageColor.Sprint(entry.Message)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s [%s] %s", ts, lvl, msg))

	for _, field := range c.extraFields {
		if val, ok := entry.Fields[field]; ok {
			sb.WriteString(fmt.Sprintf(" %s=%v", keyColor.Sprint(field), val))
		}
	}

	fmt.Fprintln(c.out, sb.String())
}
