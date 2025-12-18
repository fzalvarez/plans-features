package config

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Database Database `yaml:"database"`
}

type Database struct {
	URL string `yaml:"url"`
}

func Load() *Config {
	// 1. Intentar cargar desde sql.yaml
	cfg := &Config{}
	data, err := os.ReadFile("internal/db/sqlc.yaml")
	if err == nil {
		if err := yaml.Unmarshal(data, cfg); err != nil {
			fmt.Printf("Warning: failed to parse sql.yaml: %v\n", err)
		}
	}

	// 2. Sobrescribir con variables de entorno (prioridad alta)
	if url := os.Getenv("DATABASE_URL"); url != "" {
		cfg.Database.URL = url
	}

	// 3. DSN por defecto para desarrollo local
	/*if cfg.Database.URL == "" {
		cfg.Database.URL = "postgres://postgres:040836@localhost:5432/bookclaims?sslmode=disable"
	}*/

	return cfg
}
