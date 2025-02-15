package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Env string // Env is a custom type representing the environment

const (
	Env_Test Env = "test" // for the test environment
	Env_DEV  Env = "dev"  // for the development environment.
)

// The Config struct holds all the configuration settings for the application.
// Each field in the struct corresponds to an environment variable, specified using the env tag.
// Example: ApiServerPort string env:"APISERVER_PORT" means the ApiServerPort field will be populated
// from the APISERVER_PORT environment variable

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

// This method generates a PostgreSQL database connection URL based on the configuration.

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

// This function creates and returns a Config object
// It uses the env.ParseAs function to automatically read environment variables and populate the Config struct.
func NewConfig() (*Config, error) {
	cfg, err := env.ParseAs[Config]() //generic function to populate config fields with the env variables
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}
	return &cfg, nil

}
