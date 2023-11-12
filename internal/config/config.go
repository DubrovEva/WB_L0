package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env"`
	HTTPServer HTTPServer `yaml:"http_server"`
	DB         DB         `yaml:"db"`
	Nats       Nats       `yaml:"nats"`
}

type HTTPServer struct {
	URL string `yaml:"address" env-default:"localhost:8080"`
}

type DB struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int64  `yaml:"port"`
}

type Nats struct {
	URL       string `yaml:"url"`
	Channel   string `yaml:"channel"`
	SubClient string `yaml:"subscriber"`
	PubClient string `yaml:"publisher"`
	Cluster   string `yaml:"cluster"`
}

func MustLoad() (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig("./config/local.yml", &config); err != nil {
		return nil, err
	}
	return &config, nil
}
