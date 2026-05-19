package parser

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultParser_ParseValidJSON(t *testing.T) {
	p := DefaultParser()
	line := `{"time":"2024-01-15T10:30:00Z","level":"info","message":"server started","port":8080}`

	entry, err := p.Parse("app", line)
	require.NoError(t, err)

	assert.Equal(t, "app", entry.Source)
	assert.Equal(t, "info", entry.Level)
	assert.Equal(t, "server started", entry.Message)
	assert.Equal(t, line, entry.Raw)
	assert.Equal(t, float64(8080), entry.Fields["port"])

	expected, _ := time.Parse(time.RFC3339, "2024-01-15T10:30:00Z")
	assert.True(t, entry.Timestamp.Equal(expected))
}

func TestDefaultParser_ParseAlternativeKeys(t *testing.T) {
	p := DefaultParser()
	line := `{"ts":1705312200,"severity":"error","msg":"connection refused"}`

	entry, err := p.Parse("svc", line)
	require.NoError(t, err)

	assert.Equal(t, "error", entry.Level)
	assert.Equal(t, "connection refused", entry.Message)
	assert.Equal(t, int64(1705312200), entry.Timestamp.Unix())
}

func TestDefaultParser_ParseInvalidJSON(t *testing.T) {
	p := DefaultParser()
	_, err := p.Parse("app", "not json at all")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not valid JSON")
}

func TestDefaultParser_ParseEmptyLine(t *testing.T) {
	p := DefaultParser()
	_, err := p.Parse("app", "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty line")
}

func TestDefaultParser_MissingOptionalFields(t *testing.T) {
	p := DefaultParser()
	line := `{"custom_key":"custom_value"}`

	entry, err := p.Parse("app", line)
	require.NoError(t, err)

	assert.Empty(t, entry.Level)
	assert.Empty(t, entry.Message)
	assert.True(t, entry.Timestamp.IsZero())
	assert.Equal(t, "custom_value", entry.Fields["custom_key"])
}
