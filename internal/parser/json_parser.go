package parser

import (
	"encoding/json"
	"fmt"
	"time"
)

// LogEntry represents a parsed JSON log line with extracted fields.
type LogEntry struct {
	Raw       string
	Fields    map[string]interface{}
	Timestamp time.Time
	Level     string
	Message   string
	Source    string
}

// Parser parses structured JSON log lines.
type Parser struct {
	TimestampKeys []string
	LevelKeys     []string
	MessageKeys   []string
}

// DefaultParser returns a Parser with common field name defaults.
func DefaultParser() *Parser {
	return &Parser{
		TimestampKeys: []string{"time", "timestamp", "ts", "@timestamp"},
		LevelKeys:     []string{"level", "severity", "lvl"},
		MessageKeys:   []string{"message", "msg", "text"},
	}
}

// Parse attempts to decode a raw log line as JSON and extract known fields.
func (p *Parser) Parse(source, line string) (*LogEntry, error) {
	if line == "" {
		return nil, fmt.Errorf("empty line")
	}

	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		return nil, fmt.Errorf("not valid JSON: %w", err)
	}

	entry := &LogEntry{
		Raw:    line,
		Fields: fields,
		Source: source,
	}

	entry.Timestamp = p.extractTime(fields)
	entry.Level = p.extractString(fields, p.LevelKeys)
	entry.Message = p.extractString(fields, p.MessageKeys)

	return entry, nil
}

func (p *Parser) extractTime(fields map[string]interface{}) time.Time {
	for _, key := range p.TimestampKeys {
		val, ok := fields[key]
		if !ok {
			continue
		}
		switch v := val.(type) {
		case string:
			for _, layout := range []string{time.RFC3339Nano, time.RFC3339, time.DateTime} {
				if t, err := time.Parse(layout, v); err == nil {
					return t
				}
			}
		case float64:
			return time.Unix(int64(v), 0)
		}
	}
	return time.Time{}
}

func (p *Parser) extractString(fields map[string]interface{}, keys []string) string {
	for _, key := range keys {
		if val, ok := fields[key]; ok {
			if s, ok := val.(string); ok {
				return s
			}
		}
	}
	return ""
}
