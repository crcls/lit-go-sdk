package config

import "time"

type Config struct {
	MinimumNodeCount uint8
	Network          string
	RequestTimeout   time.Duration
	Version          string
}

func New(network string) *Config {
	return &Config{
		MinimumNodeCount: MINIMUM_NODE_COUNT,
		Network:          network,
		RequestTimeout:   REQUEST_TIMEOUT,
		Version:          VERSION,
	}
}
