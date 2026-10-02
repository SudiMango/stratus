package config

type Proxy struct {
	Port int `yaml:"port"`
}

type Route struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Path string `yaml:"path"`
	TLS  bool   `yaml:"tls"`
}

type Config struct {
	Proxy  Proxy   `yaml:"proxy"`
	Routes []Route `yaml:"routes"`
}
