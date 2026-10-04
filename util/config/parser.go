package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"go.yaml.in/yaml/v4"
)

func GetConfig(file_path string) (*Config, error) {
	data, err := os.ReadFile(file_path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Load(data, &cfg, yaml.WithKnownFields())
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ValidateConfig(cfg *Config) error {
	if cfg.Proxy.Http.Enabled == false && cfg.Proxy.Https.Enabled == false {
		return errors.New("Http and Https are both disabled. Please enable at least 1.")
	}

	if cfg.Proxy.Https.Enabled == false && cfg.Proxy.Http.RedirectToHttps == true {
		return errors.New("Cannot redirect to https when https is not enabled.")
	}

	return nil
}

func BuildURL(route Route) string {
	var target strings.Builder

	prefix := "http://"
	if route.TLS == true {
		prefix = "https://"
	}

	target.WriteString(prefix)
	target.WriteString(route.Destination)
	target.WriteString(":")
	target.WriteString(strconv.Itoa(route.Port))

	return target.String()
}
