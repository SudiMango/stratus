package config_test

import (
	"strings"
	"testing"

	"github.com/SudiMango/stratus/internal/config"
)

func TestConfigValidateAggregatesProxyAndRouteErrors(t *testing.T) {
	cfg := config.Config{
		Proxy: config.Proxy{},
		Routes: config.Routes{
			{
				Host:        "",
				Destination: "",
				Port:        0,
				Path:        "invalid",
			},
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Fatal("Validate() error = nil, want proxy and route errors")
	}

	message := strings.ToLower(err.Error())
	for _, expected := range []string{"proxy", "routes", "host", "destination", "path"} {
		if !strings.Contains(message, expected) {
			t.Errorf("Validate() error = %q, want %q", err, expected)
		}
	}
}
