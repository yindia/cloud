package config

import "fmt"

// Config holds the application configuration
type Config struct {
	ServerPort  string `envconfig:"SERVER_PORT" default:"8080"`
	WorkerCount int    `envconfig:"WORKER_COUNT" default:"100"`
	Database    DatabaseConfig
}

// DatabaseConfig holds the database connection configuration
type DatabaseConfig struct {
	Username     string `envconfig:"DB_USERNAME"`
	Password     string `envconfig:"DB_PASSWORD"`
	Host         string `envconfig:"DB_HOST"`
	Port         string `envconfig:"DB_PORT"`
	Database     string `envconfig:"DB_DATABASE"`
	SSLMode      string `envconfig:"DB_SSL_MODE" default:"require"`
	PoolMaxConns int    `envconfig:"DB_POOL_MAX_CONNS" default:"1"`
}

// ToMigrationUri returns a string for the migration package with the correct prefix
func (d DatabaseConfig) ToMigrationUri() string {
	return fmt.Sprintf("pgx5://%s:%s@%s:%s/%s?sslmode=%s",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
		d.SSLMode,
	)
}

// ToDbConnectionUri returns a connection URI for the pgx package
func (d DatabaseConfig) ToDbConnectionUri() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.Username,
		d.Password,
		d.Host,
		d.Port,
		d.Database,
		d.SSLMode,
	)
}
