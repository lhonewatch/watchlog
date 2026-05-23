// Package source provides log source management for watchlog.
package source

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/nxadm/tail"
)

// Line represents a single line read from a log source.
type Line struct {
	Source string
	Text   string
	Err    error
}

// Tailer tails one or more files and emits lines on a shared channel.
type Tailer struct {
	lines  chan Line
	tails  []*tail.Tail
}

// New creates a Tailer that watches the given file paths.
// If paths is empty, the tailer reads from stdin.
func New(paths []string) (*Tailer, error) {
	t := &Tailer{
		lines: make(chan Line, 256),
	}

	if len(paths) == 0 {
		go t.readStdin()
		return t, nil
	}

	for _, p := range paths {
		cfg := tail.Config{
			Follow:    true,
			ReOpen:    true,
			MustExist: false,
			Poll:      false,
			Location:  &tail.SeekInfo{Offset: 0, Whence: io.SeekEnd},
			Logger:    tail.DiscardingLogger,
		}
		tf, err := tail.TailFile(p, cfg)
		if err != nil {
			t.Close()
			return nil, err
		}
		t.tails = append(t.tails, tf)
		go t.forward(p, tf)
	}

	return t, nil
}

// Lines returns the channel on which tailed lines are delivered.
func (t *Tailer) Lines() <-chan Line {
	return t.lines
}

// Close stops all active tail goroutines.
func (t *Tailer) Close() {
	for _, tf := range t.tails {
		_ = tf.Stop()
	}
}

func (t *Tailer) forward(source string, tf *tail.Tail) {
	for l := range tf.Lines {
		if l.Err != nil {
			t.lines <- Line{Source: source, Err: l.Err}
			continue
		}
		t.lines <- Line{Source: source, Text: l.Text}
	}
}

func (t *Tailer) readStdin() {
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(tmp)
		if n > 0 {
			if tmp[0] == '\n' {
				t.lines <- Line{Source: "stdin", Text: string(buf)}
				buf = buf[:0]
			} else {
				buf = append(buf, tmp[0])
			}
		}
		if err != nil {
			if len(buf) > 0 {
				t.lines <- Line{Source: "stdin", Text: string(buf)}
			}
			return
		}
		_ = context.Background()
		_ = time.Second // keep import used
	}
}
