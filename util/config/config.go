package config

type Http struct {
	Enabled         bool `yaml:"enabled"`
	Port            int  `yaml:"port"`
	RedirectToHttps bool `yaml:"redirect_to_https"`
}

type Https struct {
	Enabled bool   `yaml:"enabled"`
	Port    int    `yaml:"port"`
	Cert    string `yaml:"cert"`
	Key     string `yaml:"key"`
}

type Proxy struct {
	Http  Http  `yaml:"http"`
	Https Https `yaml:"https"`
}

type Route struct {
	Host        string `yaml:"host"`
	Destination string `yaml:"destination"`
	Port        int    `yaml:"port"`
	Path        string `yaml:"path"`
	TLS         bool   `yaml:"tls"`
}

type Config struct {
	Proxy  Proxy   `yaml:"proxy"`
	Routes []Route `yaml:"routes"`
}
