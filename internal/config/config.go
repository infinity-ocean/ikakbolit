package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PG_HOST string
	PG_DB   string
	PG_USER string
	PG_PASSWORD string
	PG_PORT string
	PG_SSL  string
}

func (c* Config) Parse() error {
	if err := godotenv.Load("../../.env"); err != nil {
		return err
	}
	c.PG_HOST = os.Getenv("POSTGRES_HOST")
	c.PG_DB   = os.Getenv("POSTGRES_DB")
	c.PG_USER = os.Getenv("POSTGRES_USER")
	c.PG_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	c.PG_PORT = os.Getenv("POSTGRES_PORT")
	c.PG_SSL  = os.Getenv("POSTGRES_SSL")

	missingFields := []string{}

	if c.PG_HOST == "" {
		missingFields = append(missingFields, "POSTGRES_HOST")
	}
	if c.PG_DB == "" {
		missingFields = append(missingFields, "POSTGRES_DB")
	}
	if c.PG_USER == "" {
		missingFields = append(missingFields, "POSTGRES_USER")
	}
	if c.PG_PASSWORD == "" {
		missingFields = append(missingFields, "POSTGRES_PASSWORD")
	}
	if c.PG_PORT == "" {
		missingFields = append(missingFields, "POSTGRES_PORT")
	}
	if c.PG_SSL == "" {
		missingFields = append(missingFields, "POSTGRES_SSL")
	}

	if len(missingFields) > 0 {
		panic("missing fields in config")
	}

	return nil
}