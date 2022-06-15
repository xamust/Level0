package natsapp

import "github.com/nats-io/nats.go"

type Config struct {
	NatsAddr string `toml:"nats_addr"`
	NatsSubs string `toml:"nats_subs"`
}

func NewConfig() *Config {
	return &Config{
		NatsAddr: nats.DefaultURL,
		NatsSubs: "default",
	}
}
