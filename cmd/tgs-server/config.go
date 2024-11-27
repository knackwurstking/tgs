package main

type Config struct {
	Token           string          `json:"token" yaml:"token"`
	IPCommandConfig IPCommandConfig `json:"ip" yaml:"ip"`
}

func NewConfig() *Config {
	return &Config{}
}

type IPCommandConfig struct {
	Targets *TargetsConfig `json:"targets" yaml:"targets"`
}

type TargetsConfig struct {
	Users []UserTarget `json:"users" yaml:"users"`
	Chats []ChatTarget `json:"chats" yaml:"chats"`
}

type UserTarget struct {
	ID int64 `json:"id" yaml:"id"`
}

type ChatTarget struct {
	ID   int64  `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}
