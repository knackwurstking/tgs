package config

func New() *Config {
	return NewConfig()
}

type Config struct {
	Token           string           `json:"token" yaml:"token"`
	IPCommandConfig *IPCommandConfig `json:"ip" yaml:"ip"`
}

func NewConfig() *Config {
	return &Config{
		IPCommandConfig: NewIPCommandConfig(),
	}
}
