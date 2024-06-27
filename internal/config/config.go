package config

type Config struct {
	Addr_Port string // `env:"addr_port"`
	//DBconfig  string //`env:"DB"`
}

func NewConfig() *Config {
	return &Config{
		Addr_Port: "127.0.0.1:10000",
	}
}
