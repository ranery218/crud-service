package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Postgres struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DB       string `yaml:"db"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"postgres"`
	Session struct {
		TTL string `yaml:"ttl"`
	}
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, err
	}
	if v:= os.Getenv("POSTGRES_PASSWORD"); v != "" {
		cfg.Postgres.Password = v
	}
	return cfg, nil
}
