package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
)

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// Set default environment if not specified
	if cfg.Env == "" {
		cfg.Env = "development"
	}

	return cfg, nil
}
