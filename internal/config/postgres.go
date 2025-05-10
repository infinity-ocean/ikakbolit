package config

type Postgres struct {
	DSN string `env:"POSTGRES_DSN,notEmpty" json:"-"` // Hide in logs
}
