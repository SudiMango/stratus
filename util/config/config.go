package config

import "time"

// Proxy configs

type Http struct {
	Enabled         bool `yaml:"enabled"`
	Port            int  `yaml:"port"`
	RedirectToHttps bool `yaml:"redirect_to_https"`
}

type Https struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type Certificate struct {
	Host string `yaml:"host"`
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

type Certificates []Certificate

type Proxy struct {
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	MaxHeaderBytes    int           `yaml:"max_header_bytes"`

	Http  Http  `yaml:"http"`
	Https Https `yaml:"https"`

	Certificates Certificates `yaml:"certificates"`
}

// Route configs

type Route struct {
	Host        string `yaml:"host"`
	Destination string `yaml:"destination"`
	Port        int    `yaml:"port"`
	Path        string `yaml:"path"`
	TLS         bool   `yaml:"tls"`
}

// All configs

type Routes []Route

type Config struct {
	Proxy  Proxy  `yaml:"proxy"`
	Routes Routes `yaml:"routes"`
}
