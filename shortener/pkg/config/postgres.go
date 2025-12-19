package config

// PostgresConfig is the postgres connections config struct
type PostgresConfig struct {
	MasterDSN                    string   `env:"MASTER_DSN"`
	SlaveDSNs                    []string `env:"SLAVE_DSNS" envSeparator:" "`
	MaxOpenConnections           int      `env:"MAX_OPEN_CONNECTIONS" envDefault:"3"`
	MaxIdleConnections           int      `env:"MAX_IDLE_CONNECTIONS" envDefault:"5"`
	ConnectionMaxLifetimeSeconds int      `env:"CONNECTION_MAX_LIFETIME_SECONDS" envDefault:"0"`
}
