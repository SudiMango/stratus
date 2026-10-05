package config_test

import (
	"strings"
	"testing"
	"time"

	"github.com/SudiMango/stratus/util/config"
)

func TestGetConfigParsesDurationStrings(t *testing.T) {
	cfg, err := config.GetConfig(configFixture("valid.yaml"))
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}

	if cfg.Proxy.ReadTimeout != 30*time.Second {
		t.Errorf("ReadTimeout = %v, want 30s", cfg.Proxy.ReadTimeout)
	}
	if cfg.Proxy.ReadHeaderTimeout != 500*time.Millisecond {
		t.Errorf("ReadHeaderTimeout = %v, want 500ms", cfg.Proxy.ReadHeaderTimeout)
	}
	if cfg.Proxy.WriteTimeout != 90*time.Second {
		t.Errorf("WriteTimeout = %v, want 1m30s", cfg.Proxy.WriteTimeout)
	}
	if cfg.Proxy.IdleTimeout != 2*time.Minute {
		t.Errorf("IdleTimeout = %v, want 2m", cfg.Proxy.IdleTimeout)
	}
}

func TestGetConfigTreatsOmittedAndNullDurationsAsZero(t *testing.T) {
	tests := map[string]string{
		"omitted": "durations-omitted.yaml",
		"null":    "durations-null.yaml",
	}

	for name, fixture := range tests {
		t.Run(name, func(t *testing.T) {
			cfg, err := config.GetConfig(configFixture(fixture))
			if err != nil {
				t.Fatalf("GetConfig() error = %v", err)
			}

			if cfg.Proxy.ReadTimeout != 0 ||
				cfg.Proxy.ReadHeaderTimeout != 0 ||
				cfg.Proxy.WriteTimeout != 0 ||
				cfg.Proxy.IdleTimeout != 0 {
				t.Errorf("timeouts = %+v, want all zero", cfg.Proxy)
			}
		})
	}
}

func TestGetConfigRejectsNumericDuration(t *testing.T) {
	_, err := config.GetConfig(configFixture("invalid-numeric-duration.yaml"))
	if err == nil {
		t.Fatal("GetConfig() error = nil, want numeric duration to be rejected")
	}
}

func TestGetConfigRejectsUnknownFields(t *testing.T) {
	_, err := config.GetConfig(configFixture("unknown-field.yaml"))
	if err == nil {
		t.Fatal("GetConfig() error = nil, want unknown field to be rejected")
	}
	if !strings.Contains(strings.ToLower(err.Error()), "unexpected_setting") {
		t.Errorf("GetConfig() error = %q, want field name", err)
	}
}
