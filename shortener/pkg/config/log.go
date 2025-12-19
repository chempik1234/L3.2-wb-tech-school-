package config

// LogConfig is the config struct for logging
//
// available log levels: "trace", "debug", "info", "warn", "error", "fatal", "panic"
type LogConfig struct {
	LogLevel string `env:"LEVEL" envDefault:"info"`
}
