package config

func New() *Config {
	return NewConfig()
}

type Config struct {
	Token string              `json:"token" yaml:"token"`
	IP    *CommandConfigIP    `json:"ip" yaml:"ip"`
	Stats *CommandConfigStats `json:"stats" yaml:"stats"`
}

func NewConfig() *Config {
	return &Config{
		IP: NewCommandConfigIP(),
	}
}
