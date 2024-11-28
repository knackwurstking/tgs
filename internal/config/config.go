package config

func New() *Config {
	return NewConfig()
}

type Config struct {
	Token string           `json:"token" yaml:"token"`
	IP    *CommandConfigIP `json:"ip" yaml:"ip"`
}

func NewConfig() *Config {
	return &Config{
		IP: NewCommandConfigIP(),
	}
}
