package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Config represents the TOML configuration for bqdo.
//
// Example TOML:
//
//	directory = "sql/"
//	project_id = "my-project"
//	dataset = "analytics"
//	location = "US"
//	[vars]
//	start_date = "2025-01-01"
//	env = "prod"
//
// All fields are optional and can be overridden via CLI flags.
// Required validation is deferred to command execution.
type Config struct {
	Directory string            `toml:"directory"`
	ProjectID string            `toml:"project_id"`
	Dataset   string            `toml:"dataset"`
	Location  string            `toml:"location"`
	Vars      map[string]string `toml:"vars"`
}

// Load reads and parses a TOML config file.
// If the file does not exist, returns (zero Config, nil).
func Load(path string) (Config, error) {
	var cfg Config
	if path == "" {
		return cfg, nil
	}
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("stat config: %w", err)
	}
	if info.IsDir() {
		return cfg, fmt.Errorf("config path %q is a directory, expected file", path)
	}
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return cfg, fmt.Errorf("decode toml: %w", err)
	}
	return cfg, nil
}
