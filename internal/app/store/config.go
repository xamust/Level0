package store

type Config struct {
	DBUrl string `toml:"database_url"`
}

func NewConfig() *Config {
	return &Config{}
}
