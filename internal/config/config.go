// Package config provides configuration loading and validation for watchlog.
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// FilterRule represents a single filter rule in the configuration.
type FilterRule struct {
	Field   string `json:"field"`
	Pattern string `json:"pattern"`
}

// Config holds the full watchlog runtime configuration.
type Config struct {
	// Paths is the list of log file paths to tail.
	Paths []string `json:"paths"`

	// LevelField is the JSON key used to extract the log level.
	LevelField string `json:"level_field,omitempty"`

	// MessageField is the JSON key used to extract the log message.
	MessageField string `json:"message_field,omitempty"`

	// TimestampField is the JSON key used to extract the timestamp.
	TimestampField string `json:"timestamp_field,omitempty"`

	// FilterRules is a list of field/pattern rules used to filter log entries.
	FilterRules []FilterRule `json:"filters,omitempty"`

	// MinLevel filters out entries below this log level.
	MinLevel string `json:"min_level,omitempty"`

	// NoColor disables colorized output when true.
	NoColor bool `json:"no_color,omitempty"`
}

// Defaults applies sensible default values to unset fields.
func (c *Config) Defaults() {
	if c.LevelField == "" {
		c.LevelField = "level"
	}
	if c.MessageField == "" {
		c.MessageField = "message"
	}
	if c.TimestampField == "" {
		c.TimestampField = "time"
	}
}

// Validate returns an error if the configuration is invalid.
func (c *Config) Validate() error {
	if len(c.Paths) == 0 {
		return fmt.Errorf("config: at least one path must be specified")
	}
	for i, r := range c.FilterRules {
		if r.Field == "" {
			return fmt.Errorf("config: filter rule %d missing field name", i)
		}
		if r.Pattern == "" {
			return fmt.Errorf("config: filter rule %d missing pattern", i)
		}
	}
	return nil
}

// Load reads a JSON config file from path and returns a validated Config.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: decode %q: %w", path, err)
	}

	cfg.Defaults()

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
