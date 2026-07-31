package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Port string `yaml:"port" env:"FILESERVER_PORT" envDefault:"1234"`
}

func New(configPath string) (*Config, error) {
	var cfg Config

	if configPath != "" {
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			return nil, err
		}
	} else {
		if err := cleanenv.ReadEnv(&cfg); err != nil {
			return nil, err
		}
	}
	return &cfg, nil
}
