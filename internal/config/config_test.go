package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/watchlog/internal/config"
)

func writeTempConfig(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal config: %v", err)
	}
	p := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(p, b, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return p
}

func TestLoad_ValidConfig(t *testing.T) {
	p := writeTempConfig(t, map[string]any{
		"paths": []string{"/var/log/app.log"},
	})
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LevelField != "level" {
		t.Errorf("expected default level_field=level, got %q", cfg.LevelField)
	}
	if cfg.MessageField != "message" {
		t.Errorf("expected default message_field=message, got %q", cfg.MessageField)
	}
	if cfg.TimestampField != "time" {
		t.Errorf("expected default timestamp_field=time, got %q", cfg.TimestampField)
	}
}

func TestLoad_CustomFields(t *testing.T) {
	p := writeTempConfig(t, map[string]any{
		"paths":           []string{"/tmp/test.log"},
		"level_field":     "severity",
		"message_field":   "msg",
		"timestamp_field": "ts",
	})
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.LevelField != "severity" {
		t.Errorf("got level_field=%q", cfg.LevelField)
	}
}

func TestLoad_MissingPaths(t *testing.T) {
	p := writeTempConfig(t, map[string]any{})
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for missing paths")
	}
}

func TestLoad_InvalidFilterRule(t *testing.T) {
	p := writeTempConfig(t, map[string]any{
		"paths":   []string{"/tmp/test.log"},
		"filters": []map[string]any{{"field": "", "pattern": "error"}},
	})
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid filter rule")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDefaults_DoNotOverrideExisting(t *testing.T) {
	cfg := &config.Config{
		Paths:          []string{"/tmp/x.log"},
		LevelField:     "lvl",
		MessageField:   "text",
		TimestampField: "@timestamp",
	}
	cfg.Defaults()
	if cfg.LevelField != "lvl" {
		t.Errorf("Defaults() overwrote LevelField")
	}
}
