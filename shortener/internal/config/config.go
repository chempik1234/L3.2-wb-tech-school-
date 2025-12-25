package config

import (
	"fmt"
	config2 "github.com/chempik1234/L3.2-wb-tech-school-/shortener/pkg/config"
	"github.com/wb-go/wbf/config"
)

// AppConfig is THE whole config struct
type AppConfig struct {
	HTTPServerConfig config2.HTTPServerConfig `env-prefix:"SHORTENER_HTTP_SERVER_"`
	LogConfig        config2.LogConfig        `env-prefix:"SHORTENER_LOG_"`

	CacheConfig CacheConfig `env-prefix:"SHORTENER_CACHE_"`

	PostgresConfig config2.PostgresConfig `env-prefix:"SHORTENER_POSTGRES_"`
	RedisConfig    config2.RedisConfig    `env-prefix:"SHORTENER_REDIS_"`

	PostgresRetryConfig config2.RetryStrategyConfig `env-prefix:"SHORTENER_RETRY_POSTGRES_"`
	RedisRetryConfig    config2.RetryStrategyConfig `env-prefix:"SHORTENER_RETRY_REDIS_"`

	MaxLinkLen            int `env:"SHORTENER_MAX_LINK_LEN"`
	BatchingPeriodSeconds int `env:"SHORTENER_BATCHING_PERIOD_SECONDS"`
}

// NewAppConfig creates a new struct of "THE config"
//
// has some defaults
//
// message to code reviewer - please destroy the wbf/config package with nuke and use `cleanenv` with env/yaml tags
func NewAppConfig(configFilePath, envFilePath string) (*AppConfig, error) {
	var appConfig *AppConfig

	cfg := config.New()

	//region defaults
	cfg.SetDefault("shortener.http_server.port", 8080)
	cfg.SetDefault("shortener.log.level", "info")

	cfg.SetDefault("shortener.postgres.max_open_connections", 2)
	cfg.SetDefault("shortener.postgres.max_idle_connections", 2)
	cfg.SetDefault("shortener.postgres.connection_max_lifetime_seconds", 0)

	cfg.SetDefault("shortener.cache_config.min_requests_before_caching", 3)

	cfg.SetDefault("shortener.redis.db", 0)
	cfg.SetDefault("shortener.redis.ttl_seconds", 20)

	// retry: attempts
	cfg.SetDefault("shortener.retry_redis.attempts", 2)
	cfg.SetDefault("shortener.retry_postgres.attempts", 2)
	cfg.SetDefault("shortener.retry_kafka.attempts", 2)

	// retry: delay
	cfg.SetDefault("shortener.retry_redis.delay_milliseconds", 300)
	cfg.SetDefault("shortener.retry_postgres.delay_milliseconds", 300)
	cfg.SetDefault("shortener.retry_kafka.delay_milliseconds", 300)

	// retry: backoff
	cfg.SetDefault("shortener.retry_redis.backoff", 1.2)
	cfg.SetDefault("shortener.retry_postgres.backoff", 1.5)
	cfg.SetDefault("shortener.retry_kafka.backoff", 1.5)

	cfg.SetDefault("shortener.max_link_len", 6)
	cfg.SetDefault("shortener.batching_period_seconds", 10)
	//endregion

	// region flags

	// does it work as "set as default"? if I specify 8080 as default flag value, then it's always non-empty
	// how about that? If it's empty, will it override my .env?
	// _ = cfg.DefineFlag("q", "queue", "RABBITMQ_QUEUE", "", "RabbitMQ queue name")

	//endregion

	var err error
	if len(configFilePath) > 0 {
		err = cfg.LoadConfigFiles(envFilePath)
	}
	if len(envFilePath) > 0 {
		err = cfg.LoadEnvFiles(envFilePath)
	}
	cfg.EnableEnv("")

	if err != nil {
		return appConfig, fmt.Errorf("failed to load config: %w", err)
	}

	appConfig = &AppConfig{
		HTTPServerConfig: config2.HTTPServerConfig{
			Port: cfg.GetInt("shortener.http_server.port"),
		},
		LogConfig: config2.LogConfig{
			LogLevel: cfg.GetString("shortener.log.level"),
		},
		CacheConfig: CacheConfig{
			MinRequestsBeforeCaching: cfg.GetInt("shortener.cache_config.min_requests_before_caching"),
		},
		PostgresConfig: config2.PostgresConfig{
			MasterDSN:                    cfg.GetString("shortener.postgres.master_dsn"),
			SlaveDSNs:                    cfg.GetStringSlice("shortener.postgres.slave_dsns"),
			MaxOpenConnections:           cfg.GetInt("shortener.postgres.max_open_connections"),
			MaxIdleConnections:           cfg.GetInt("shortener.postgres.max_idle_connections"),
			ConnectionMaxLifetimeSeconds: cfg.GetInt("shortener.postgres.connection_max_lifetime_seconds"),
		},
		RedisConfig: config2.RedisConfig{
			Addr:     cfg.GetString("shortener.redis.addr"),
			Password: cfg.GetString("shortener.redis.password"),
			DB:       cfg.GetInt("shortener.redis.db"),
			// #i_cant_commit_to_wbf_so_destroy_wbf TTLSeconds: 0,
		},
		PostgresRetryConfig: config2.RetryStrategyConfig{
			Attempts:          cfg.GetInt("shortener.retry_postgres.attempts"),
			DelayMilliseconds: cfg.GetInt("shortener.retry_postgres.delay_milliseconds"),
			Backoff:           cfg.GetFloat64("shortener.retry_postgres.backoff"),
		},
		RedisRetryConfig: config2.RetryStrategyConfig{
			Attempts:          cfg.GetInt("shortener.retry_redis.attempts"),
			DelayMilliseconds: cfg.GetInt("shortener.retry_redis.delay_milliseconds"),
			Backoff:           cfg.GetFloat64("shortener.retry_redis.backoff"),
		},
		MaxLinkLen:            cfg.GetInt("shortener.max_link_len"),
		BatchingPeriodSeconds: cfg.GetInt("shortener.batching_period_seconds"),
	}

	return appConfig, nil
}
