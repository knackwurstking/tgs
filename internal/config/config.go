package config

type Config struct {
	Token string `json:"token" yaml:"token"`
}

func New() *Config {
	return &Config{}
}
