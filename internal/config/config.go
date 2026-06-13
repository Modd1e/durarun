package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const defaultDotEnvPath = ".env"

type Environment string

const (
	EnvironmentDevelopment Environment = "DEV"
	EnvironmentProduction  Environment = "PROD"
	EnvironmentTest        Environment = "TEST"
)

type Config struct {
	Environment Environment `env:"ENV,required,notEmpty"`
	Postgres    PostgresConfig
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST,required,notEmpty"`
	Port     uint16 `env:"POSTGRES_PORT" envDefault:"5432"`
	User     string `env:"POSTGRES_USER,required,notEmpty"`
	Password string `env:"POSTGRES_PASS,required,notEmpty"`
	Database string `env:"POSTGRES_DATABASE,required,notEmpty"`
}

func Load() (Config, error) {
	return load(defaultDotEnvPath, os.Environ())
}

func (e Environment) IsProduction() bool {
	return e == EnvironmentProduction
}

func (e Environment) String() string {
	return string(e)
}

func load(dotEnvPath string, processEnvironment []string) (Config, error) {
	environment, err := readEnvironment(dotEnvPath, processEnvironment)
	if err != nil {
		return Config{}, fmt.Errorf("load configuration sources: %w", err)
	}

	cfg, err := env.ParseAsWithOptions[Config](env.Options{
		Environment: environment,
	})
	if err != nil {
		return Config{}, fmt.Errorf("parse configuration: %w", err)
	}

	cfg.Environment = normalizeEnvironment(cfg.Environment)

	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("validate configuration: %w", err)
	}

	return cfg, nil
}

func normalizeEnvironment(environment Environment) Environment {
	switch strings.ToUpper(strings.TrimSpace(environment.String())) {
	case "DEV", "DEVELOPMENT":
		return EnvironmentDevelopment
	case "PROD", "PRODUCTION":
		return EnvironmentProduction
	case "TEST", "TESTING":
		return EnvironmentTest
	default:
		return Environment(strings.ToUpper(strings.TrimSpace(environment.String())))
	}
}

func readEnvironment(dotEnvPath string, processEnvironment []string) (map[string]string, error) {
	values, err := godotenv.Read(dotEnvPath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("read %q: %w", dotEnvPath, err)
		}

		values = make(map[string]string)
	}

	for _, item := range processEnvironment {
		key, value, ok := strings.Cut(item, "=")
		if ok {
			values[key] = value
		}
	}

	return values, nil
}

func (c Config) validate() error {
	var validationErrors []error

	switch c.Environment {
	case EnvironmentDevelopment, EnvironmentProduction, EnvironmentTest:
	default:
		validationErrors = append(
			validationErrors,
			fmt.Errorf(
				"ENV must be one of %q, %q, or %q, got %q",
				EnvironmentDevelopment,
				EnvironmentProduction,
				EnvironmentTest,
				c.Environment,
			),
		)
	}

	if strings.TrimSpace(c.Postgres.Host) == "" {
		validationErrors = append(validationErrors, errors.New("POSTGRES_HOST must not be blank"))
	}
	if c.Postgres.Port == 0 {
		validationErrors = append(validationErrors, errors.New("POSTGRES_PORT must be between 1 and 65535"))
	}
	if strings.TrimSpace(c.Postgres.User) == "" {
		validationErrors = append(validationErrors, errors.New("POSTGRES_USER must not be blank"))
	}
	if c.Postgres.Password == "" {
		validationErrors = append(validationErrors, errors.New("POSTGRES_PASS must not be empty"))
	}
	if strings.TrimSpace(c.Postgres.Database) == "" {
		validationErrors = append(validationErrors, errors.New("POSTGRES_DATABASE must not be blank"))
	}

	return errors.Join(validationErrors...)
}
