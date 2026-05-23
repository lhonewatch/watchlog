package output

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/watchlog/internal/parser"
)

func makeFormatterEntry(level, msg string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		Level:     level,
		Message:   msg,
		Fields:    map[string]interface{}{},
	}
}

func TestFormatter_WriteEmitsNewline(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, true)
	entry := makeFormatterEntry("INFO", "hello")

	if err := f.Write(entry); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.HasSuffix(output, "\n") {
		t.Errorf("expected output to end with newline, got %q", output)
	}
}

func TestFormatter_WriteContainsMessage(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, true)
	entry := makeFormatterEntry("WARN", "disk almost full")

	_ = f.Write(entry)

	if !strings.Contains(buf.String(), "disk almost full") {
		t.Errorf("expected message in output, got %q", buf.String())
	}
}

func TestFormatter_WriteContainsLevel(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, true)
	entry := makeFormatterEntry("ERROR", "boom")

	_ = f.Write(entry)

	if !strings.Contains(buf.String(), "ERROR") {
		t.Errorf("expected level in output, got %q", buf.String())
	}
}

func TestFormatter_WriteExtraFields(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, true)
	entry := makeFormatterEntry("DEBUG", "trace")
	entry.Fields["request_id"] = "abc-123"

	_ = f.Write(entry)

	if !strings.Contains(buf.String(), "request_id") {
		t.Errorf("expected extra field key in output, got %q", buf.String())
	}
}

func TestFormatter_NoColorDoesNotContainEscapes(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(&buf, true)
	entry := makeFormatterEntry("INFO", "plain")

	_ = f.Write(entry)

	if strings.Contains(buf.String(), "\x1b[") {
		t.Errorf("expected no ANSI escape codes in no-colour mode, got %q", buf.String())
	}
}
