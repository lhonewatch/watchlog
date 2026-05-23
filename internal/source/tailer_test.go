package source_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/watchlog/internal/source"
)

func TestTailer_NewWithNoPathsReturnsTailer(t *testing.T) {
	// Tailer with no paths should succeed (reads stdin).
	tlr, err := source.New([]string{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tlr == nil {
		t.Fatal("expected non-nil tailer")
	}
	tlr.Close()
}

func TestTailer_LinesChannelIsNotNil(t *testing.T) {
	tlr, err := source.New([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer tlr.Close()

	if tlr.Lines() == nil {
		t.Fatal("Lines() channel should not be nil")
	}
}

func TestTailer_TailsExistingFile(t *testing.T) {
	tmp := t.TempDir()
	logFile := filepath.Join(tmp, "test.log")

	f, err := os.Create(logFile)
	if err != nil {
		t.Fatalf("could not create temp file: %v", err)
	}

	tlr, err := source.New([]string{logFile})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer tlr.Close()

	// Write a line after the tailer starts.
	time.Sleep(50 * time.Millisecond)
	_, _ = f.WriteString(`{"level":"info","msg":"hello"}` + "\n")
	_ = f.Close()

	select {
	case line := <-tlr.Lines():
		if line.Err != nil {
			t.Fatalf("unexpected line error: %v", line.Err)
		}
		if line.Source != logFile {
			t.Errorf("expected source %q, got %q", logFile, line.Source)
		}
		if line.Text == "" {
			t.Error("expected non-empty line text")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for line from tailed file")
	}
}

func TestTailer_CloseIsIdempotent(t *testing.T) {
	tlr, err := source.New([]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Calling Close multiple times should not panic.
	tlr.Close()
	tlr.Close()
}
