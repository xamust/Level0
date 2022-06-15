package service

import (
	"service/internal/app/natsapp"
	"service/internal/app/store"
)

type Config struct {
	BindAddr string `toml:"bind_addr"`
	LogLevel string `toml:"log_level"`
	Store    *store.Config
	NatsApp  *natsapp.Config
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "info",
		Store:    store.NewConfig(),
		NatsApp:  natsapp.NewConfig(),
	}
}
