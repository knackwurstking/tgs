package config

func New() *Config {
	return NewConfig()
}

type Config struct {
	Token   string                `json:"token" yaml:"token"`
	IP      *CommandConfigIP      `json:"ip,omitempty" yaml:"ip,omitempty"`
	Stats   *CommandConfigStats   `json:"stats,omitempty" yaml:"stats,omitempty"`
	Journal *CommandConfigJournal `json:"journal,omitempty" yaml:"journal,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		IP: NewCommandConfigIP(),
	}
}
