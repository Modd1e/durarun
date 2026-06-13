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
	DatabaseURL string      `env:"DATABASE_URL,required,notEmpty"`
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

	if strings.TrimSpace(c.DatabaseURL) == "" {
		validationErrors = append(validationErrors, errors.New("DatabaseURL must not be blank"))
	}

	return errors.Join(validationErrors...)
}
