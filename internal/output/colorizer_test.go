package output_test

import (
	"bytes"
	"strings"
	"testing"

	"watchlog/internal/output"
	"watchlog/internal/parser"
)

func makeEntry(ts, level, msg string, fields map[string]interface{}) *parser.Entry {
	if fields == nil {
		fields = map[string]interface{}{}
	}
	return &parser.Entry{Timestamp: ts, Level: level, Message: msg, Fields: fields}
}

func TestColorizer_ContainsTimestamp(t *testing.T) {
	var buf bytes.Buffer
	c := output.NewColorizer(&buf, nil, true)
	c.Write(makeEntry("2024-01-01T00:00:00Z", "info", "hello", nil))
	if !strings.Contains(buf.String(), "2024-01-01T00:00:00Z") {
		t.Errorf("expected timestamp in output, got: %s", buf.String())
	}
}

func TestColorizer_ContainsLevel(t *testing.T) {
	var buf bytes.Buffer
	c := output.NewColorizer(&buf, nil, true)
	c.Write(makeEntry("2024-01-01T00:00:00Z", "error", "boom", nil))
	out := buf.String()
	if !strings.Contains(out, "ERROR") {
		t.Errorf("expected ERROR in output, got: %s", out)
	}
}

func TestColorizer_ContainsMessage(t *testing.T) {
	var buf bytes.Buffer
	c := output.NewColorizer(&buf, nil, true)
	c.Write(makeEntry("2024-01-01T00:00:00Z", "debug", "test message", nil))
	if !strings.Contains(buf.String(), "test message") {
		t.Errorf("expected message in output, got: %s", buf.String())
	}
}

func TestColorizer_ExtraFieldsIncluded(t *testing.T) {
	var buf bytes.Buffer
	c := output.NewColorizer(&buf, []string{"request_id"}, true)
	fields := map[string]interface{}{"request_id": "abc-123", "ignored": "nope"}
	c.Write(makeEntry("2024-01-01T00:00:00Z", "info", "req", fields))
	out := buf.String()
	if !strings.Contains(out, "request_id=abc-123") {
		t.Errorf("expected request_id in output, got: %s", out)
	}
	if strings.Contains(out, "ignored") {
		t.Errorf("unexpected field 'ignored' in output: %s", out)
	}
}

func TestColorizer_MissingExtraFieldSkipped(t *testing.T) {
	var buf bytes.Buffer
	c := output.NewColorizer(&buf, []string{"trace_id"}, true)
	c.Write(makeEntry("2024-01-01T00:00:00Z", "warn", "missing field", nil))
	if strings.Contains(buf.String(), "trace_id") {
		t.Errorf("should not include absent field in output")
	}
}
