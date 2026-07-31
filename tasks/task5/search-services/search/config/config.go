package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel      string `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	Address       string `yaml:"search_address" env:"SEARCH_ADDRESS" env-default:"localhost:80"`
	UpdateAddress string `yaml:"update_address" env:"UPDATE_ADDRESS" env-default:"update:81"`
	WordsAddress  string `yaml:"words_address" env:"WORDS_ADDRESS" env-default:"words:82"`
	Concurrency   int    `yaml:"concurrency" env:"SEARCH_CONCURRENCY" env-default:"1"`
}

func MustLoad(configPath string) Config {
	var cfg Config

	if configPath != "" {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			log.Fatalf("cannot read config %q: %s", configPath, err)
		}
	} else {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			log.Fatalf("cannot read env: %s", err)
		}
	}

	return cfg
}
