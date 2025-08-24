package ratelimit

import (
	"context"
	"time"
)

type Limiter interface {
	Allow(ctx context.Context, key string, rate int, window time.Duration) (bool, error)
	Reset(ctx context.Context, key string) error
	GetCount(ctx context.Context, key string) (int, error)
}

type Config struct {
	DefaultRate   int
	DefaultWindow time.Duration
	Rules         map[string]Rule
}

type Rule struct {
	Rate   int
	Window time.Duration
}

func NewConfig() *Config {
	return &Config{
		DefaultRate:   100,
		DefaultWindow: 2 * time.Minute,
		Rules: map[string]Rule{
			"summoner":   {Rate: 100, Window: 2 * time.Minute},
			"match":      {Rate: 100, Window: 2 * time.Minute},
			"league":     {Rate: 100, Window: 2 * time.Minute},
			"match-list": {Rate: 1000, Window: 10 * time.Second},
		},
	}
}

func (c *Config) GetRule(endpoint string) Rule {
	if rule, exists := c.Rules[endpoint]; exists {
		return rule
	}
	return Rule{Rate: c.DefaultRate, Window: c.DefaultWindow}
}
