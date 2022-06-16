package cashdata

type Config struct {
	CountOfCash int `toml:"count_of_cash"`
}

func NewConfig() *Config {
	return &Config{
		CountOfCash: 1,
	}
}
