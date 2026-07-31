package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	WordsAddr  string     `yaml:"words_address" env:"WORDS_ADDRESS" env-default:"localhost:81"`
	LogLevel   string     `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	HTTPServer HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address string        `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-default:"words:8080"`
	Timeout time.Duration `yaml:"timeout" env:"HTTP_SERVER_TIMEOUT" env-default:"5s"`
}

func MustLoad(configPath string) *Config {
	var cfg Config

	if configPath != "" {
		if _, err := os.Stat(configPath); err != nil {
			log.Fatalf("error opening config file: %s", err)
		}

		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			log.Fatalf("cannot read config %q: %s", configPath, err)
		}
	} else {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("cannot read env: %s", err)
		}
	}

	return &cfg
}
