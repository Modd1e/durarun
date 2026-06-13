package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadParsesTypedConfiguration(t *testing.T) {
	cfg, err := load(filepath.Join(t.TempDir(), ".env"), validEnvironment())
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Environment != EnvironmentDevelopment {
		t.Fatalf("Environment = %q, want %q", cfg.Environment, EnvironmentDevelopment)
	}
	if cfg.Postgres.Port != 5432 {
		t.Fatalf("Postgres.Port = %d, want 5432", cfg.Postgres.Port)
	}
	if cfg.Postgres.Password != "secret" {
		t.Fatalf("Postgres.Password = %q, want %q", cfg.Postgres.Password, "secret")
	}
}

func TestLoadUsesDotEnvAsFallback(t *testing.T) {
	dotEnvPath := filepath.Join(t.TempDir(), ".env")
	dotEnv := strings.Join([]string{
		"ENV=DEV",
		"POSTGRES_HOST=dotenv-db",
		"POSTGRES_PORT=5433",
		"POSTGRES_USER=dotenv-user",
		"POSTGRES_PASS=dotenv-password",
		"POSTGRES_DATABASE=dotenv-database",
	}, "\n")

	if err := os.WriteFile(dotEnvPath, []byte(dotEnv), 0o600); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	cfg, err := load(dotEnvPath, []string{
		"POSTGRES_HOST=process-db",
		"POSTGRES_PASS=process-password",
	})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.Postgres.Host != "process-db" {
		t.Fatalf("Postgres.Host = %q, want process environment value", cfg.Postgres.Host)
	}
	if cfg.Postgres.Password != "process-password" {
		t.Fatalf("Postgres.Password = %q, want process environment value", cfg.Postgres.Password)
	}
	if cfg.Postgres.User != "dotenv-user" {
		t.Fatalf("Postgres.User = %q, want .env fallback value", cfg.Postgres.User)
	}
	if cfg.Postgres.Port != 5433 {
		t.Fatalf("Postgres.Port = %d, want 5433", cfg.Postgres.Port)
	}
}

func TestLoadNormalizesEnvironment(t *testing.T) {
	environment := validEnvironment()
	environment[0] = "ENV= production "

	cfg, err := load(filepath.Join(t.TempDir(), ".env"), environment)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if !cfg.Environment.IsProduction() {
		t.Fatalf("Environment = %q, want production", cfg.Environment)
	}
}

func TestLoadRejectsMissingRequiredVariable(t *testing.T) {
	environment := validEnvironment()
	environment = environment[:len(environment)-1]

	_, err := load(filepath.Join(t.TempDir(), ".env"), environment)
	if err == nil {
		t.Fatal("load config succeeded, want an error")
	}
	if !strings.Contains(err.Error(), "POSTGRES_DATABASE") {
		t.Fatalf("error = %q, want missing variable name", err)
	}
}

func TestLoadRejectsInvalidPort(t *testing.T) {
	environment := validEnvironment()
	environment[2] = "POSTGRES_PORT=70000"

	_, err := load(filepath.Join(t.TempDir(), ".env"), environment)
	if err == nil {
		t.Fatal("load config succeeded, want an error")
	}
	if !strings.Contains(err.Error(), `field "Port"`) {
		t.Fatalf("error = %q, want invalid field name", err)
	}
}

func TestLoadRejectsUnknownEnvironment(t *testing.T) {
	environment := validEnvironment()
	environment[0] = "ENV=STAGING"

	_, err := load(filepath.Join(t.TempDir(), ".env"), environment)
	if err == nil {
		t.Fatal("load config succeeded, want an error")
	}
	if !strings.Contains(err.Error(), `ENV must be one of`) {
		t.Fatalf("error = %q, want environment validation error", err)
	}
}

func TestLoadRejectsMalformedDotEnv(t *testing.T) {
	dotEnvPath := filepath.Join(t.TempDir(), ".env")
	if err := os.WriteFile(dotEnvPath, []byte("INVALID LINE\n"), 0o600); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	_, err := load(dotEnvPath, validEnvironment())
	if err == nil {
		t.Fatal("load config succeeded, want an error")
	}
	if !strings.Contains(err.Error(), "load configuration sources") {
		t.Fatalf("error = %q, want source loading error", err)
	}
}

func validEnvironment() []string {
	return []string{
		"ENV=DEV",
		"POSTGRES_HOST=localhost",
		"POSTGRES_PORT=5432",
		"POSTGRES_USER=postgres",
		"POSTGRES_PASS=secret",
		"POSTGRES_DATABASE=durarun",
	}
}
