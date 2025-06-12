package main

type Config struct {
	Token string `json:"token" yaml:"token"`
}

func NewConfig() *Config {
	return &Config{}
}
