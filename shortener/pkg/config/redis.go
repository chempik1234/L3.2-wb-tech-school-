package config

// RedisConfig is the redis connection config struct
type RedisConfig struct {
	Addr       string `env:"ADDR" envDefault:"localhost:6379"`
	Password   string `env:"PASSWORD" envDefault:""`
	DB         int    `env:"DB" envDefault:"0"`
	TTLSeconds int    `env:"TTL_SECONDS" envDefault:"0"`
}
