// Package parser provides JSON log line parsing for watchlog.
//
// It decodes structured JSON log entries and extracts well-known fields
// such as timestamp, log level, and message from a variety of common
// field name conventions used by popular logging libraries (zap, logrus,
// zerolog, etc.).
//
// Basic usage:
//
//	p := parser.DefaultParser()
//	entry, err := p.Parse("my-service", rawLine)
//	if err != nil {
//		// handle non-JSON or empty lines
//	}
//	fmt.Println(entry.Level, entry.Message)
package parser
