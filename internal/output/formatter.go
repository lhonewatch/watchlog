package output

import (
	"fmt"
	"io"

	"github.com/user/watchlog/internal/parser"
)

// Formatter combines a Colorizer with an io.Writer so that a single
// call to Write renders and emits a log entry.
type Formatter struct {
	colorizer *Colorizer
	w         io.Writer
}

// NewFormatter returns a Formatter that writes colourised log lines to w.
// If noColor is true the underlying Colorizer will omit ANSI escape codes.
func NewFormatter(w io.Writer, noColor bool) *Formatter {
	return &Formatter{
		colorizer: NewColorizer(noColor),
		w:         w,
	}
}

// Write renders entry as a colourised line and writes it to the
// underlying writer followed by a newline.  It returns any write error.
func (f *Formatter) Write(entry parser.Entry) error {
	line := f.colorizer.Format(entry)
	_, err := fmt.Fprintln(f.w, line)
	return err
}
