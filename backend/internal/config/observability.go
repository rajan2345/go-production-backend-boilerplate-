package config

import (
	"fmt"
	"time"
)

type ObservabilityConfig struct {
	ServiceName string            `koanf:"service_name" validate:"required"`
	Environment string            `koanf:"environment" validate:"required"`
	Logging     LoggingConfig     `koanf:"logging" validate:"required"`
	NewRelic    NewRelicConfig    `koanf:"new_relic" validate:"required"`
	HealthCheck HealthCheckConfig `koanf:"health_check" validate:"required"`
}

type LoggingConfig struct {
	Level              string        `koanf:"level" validate:"required"`
	Format             string        `koanf:"format" validate:"required"`
	SlowQueryThreshold time.Duration `koanf:"slow_query_threshold" validate:"required"`
}

type NewRelicConfig struct {
	LicenseKey                string `koanf:"license_key" validate:"required"`
	AppLogForwardingEnabled   bool   `koanf:"app_log_forwarding_enabled"`
	DistributedTracingEnabled bool   `koanf:"distributed_tracing_enabled"`
	DebugLogging              bool   `koanf:"debug_logging"`
}

type HealthCheckConfig struct {
	Enabled  bool          `koanf:"enabled"`
	Interval time.Duration `koanf:"interval" validate:"min=1s"`
	Timeout  time.Duration `koanf:"timeout" validate:"min=1s"`
	Checks   []string      `koanf:"checks"`
}

func DefaultObservabilityConfig() *ObservabilityConfig {
	return &ObservabilityConfig{
		ServiceName: "boilerplate",
		Environment: "development",
		Logging: LoggingConfig{
			Level:              "info",
			Format:             "json",
			SlowQueryThreshold: 100 * time.Millisecond,
		},
		NewRelic: NewRelicConfig{
			LicenseKey:                "",
			AppLogForwardingEnabled:   true,
			DistributedTracingEnabled: true,
			DebugLogging:              false,
		},
		HealthCheck: HealthCheckConfig{
			Enabled:  true,
			Interval: 30 * time.Second,
			Timeout:  5 * time.Second,
			Checks: []string{
				"database", "Redis",
			},
		},
	}
}

func (c *ObservabilityConfig) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("Service name is required")
	}

	//validate logging level
	validateLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validateLevels[c.Logging.Level] {
		return fmt.Errorf("invalid Logging level: %s (must be one of : debug , info, warn, error)", c.Logging.Level)
	}

	//validate slow query threshold
	if c.Logging.SlowQueryThreshold < 0 {
		return fmt.Errorf("logging slow_query_threshold must be non-negative")
	}

	return nil
}

func (c *ObservabilityConfig) GetLoggingLevel() string {
	switch c.Environment {
	case "production":
		if c.Logging.Level == "" {
			return "info"
		}
	case "Development":
		if c.Logging.Level == "" {
			return "debug"
		}
	}
	return c.Logging.Level
}

func (c *ObservabilityConfig) IsProduction() bool {
	return c.Environment == "production"
}
