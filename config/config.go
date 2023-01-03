package config

import (
	"os"
	"time"
)

type Config struct {
	Debug            bool
	MinimumNodeCount uint8
	Network          string
	RequestTimeout   time.Duration
	Version          string
}

func New(network string) *Config {
	c := &Config{
		MinimumNodeCount: MINIMUM_NODE_COUNT,
		Network:          network,
		RequestTimeout:   REQUEST_TIMEOUT,
		Version:          VERSION,
	}

	if val, ok := os.LookupEnv("LIT_DEBUG"); ok && val == "true" {
		c.Debug = true
	}

	return c
}
