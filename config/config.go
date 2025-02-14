package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Env string

const (
	Env_Test Env = "test"
	Env_DEV  Env = "dev"
)

type Config struct {
	ApiServerPort string `env:"APISERVER_PORT"`
	ApiServerHost string `env:"APISERVER_HOST"`
	DbName        string `env:"DB_NAME"`
	DbHost        string `env:"DB_HOST"`
	DbPort        string `env:"DB_PORT"`
	DbPortTest    string `env:"DB_PORT_TEST"`
	DbUser        string `env:"DB_USER"`
	DbPassword    string `env:"DB_PASSWORD"`
	Env           Env    `env:"ENV" envDefault:"dev"`
	ProjectRoot   string `env:"PROJECT_ROOT"`
}

func (c *Config) DbUrl() string {
	port := c.DbPort
	if c.Env == Env_Test {
		port = c.DbPortTest
	}
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.DbUser,
		c.DbPassword,
		c.DbHost,
		port,
		c.DbName,
	)
}

func NewConfig() (*Config, error) {
	cfg, err := env.ParseAs[Config]() //generic function to populate config fields with the env variables
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}
	return &cfg, nil

}
