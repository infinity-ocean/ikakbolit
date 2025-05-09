package config

type Postgres struct {
	DSN string `env:"PG_DSN,notEmpty" json:"-"` // Hide in logs
}
