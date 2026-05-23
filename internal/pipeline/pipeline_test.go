package pipeline_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/watchlog/internal/filter"
	"github.com/user/watchlog/internal/output"
	"github.com/user/watchlog/internal/parser"
	"github.com/user/watchlog/internal/pipeline"
)

func TestPipeline_WritesMatchingEntry(t *testing.T) {
	p := parser.DefaultParser{}
	f := filter.New(nil) // no rules — matches everything
	var buf bytes.Buffer
	c := output.NewColorizer(false)
	pl := pipeline.New(p, f, c, &buf)

	lines := make(chan string, 1)
	lines <- `{"time":"2024-01-01T00:00:00Z","level":"info","msg":"hello"}`
	close(lines)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := pl.Run(ctx, lines); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "hello") {
		t.Errorf("expected output to contain 'hello', got: %q", buf.String())
	}
}

func TestPipeline_SkipsInvalidJSON(t *testing.T) {
	p := parser.DefaultParser{}
	f := filter.New(nil)
	var buf bytes.Buffer
	c := output.NewColorizer(false)
	pl := pipeline.New(p, f, c, &buf)

	lines := make(chan string, 2)
	lines <- `not json at all`
	lines <- `{"time":"2024-01-01T00:00:00Z","level":"info","msg":"valid"}`
	close(lines)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := pl.Run(ctx, lines); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "valid") {
		t.Errorf("expected 'valid' in output, got: %q", buf.String())
	}
}

func TestPipeline_RespectsContextCancellation(t *testing.T) {
	p := parser.DefaultParser{}
	f := filter.New(nil)
	var buf bytes.Buffer
	c := output.NewColorizer(false)
	pl := pipeline.New(p, f, c, &buf)

	lines := make(chan string) // never closed, never written
	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan error, 1)
	go func() { done <- pl.Run(ctx, lines) }()

	cancel()
	select {
	case err := <-done:
		if err != context.Canceled {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("pipeline did not stop after context cancellation")
	}
}
