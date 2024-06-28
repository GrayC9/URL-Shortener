package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr_Port string
	//DBconfig  string
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		return nil
	}
	return &Config{
		Addr_Port: getAddr("addr_port"),
	}
}

func getAddr(addr string) string {
	return os.Getenv(addr)
}
