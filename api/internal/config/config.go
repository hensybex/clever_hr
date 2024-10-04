// config.go
package config

import (
	"fmt"
	"os"
)

type Config struct {
	DatabaseURL string
	SearxngURL  string
}

func LoadConfig() (*Config, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	sslMode := os.Getenv("POSTGRES_SSL_MODE")

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode)

	searxngURL := os.Getenv("SEARXNG_URL")

	return &Config{
		DatabaseURL: databaseURL,
		SearxngURL:  searxngURL,
	}, nil
}
