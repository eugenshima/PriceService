// Package config contains the configurations
package config

import (
	"github.com/caarlos0/env/v9"
)

// Config struct contains our configuration variables
type Config struct {
	RedisConnectionString string `env:"REDIS_CONNECTION_STRING"`
}

// NewConfig creates a new Config
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
