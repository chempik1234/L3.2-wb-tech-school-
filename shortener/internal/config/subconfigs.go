package config

// CacheConfig - config for semantic cache settings
type CacheConfig struct {
	MinRequestsBeforeCaching int `env:"MIN_REQUESTS_BEFORE_CACHING" env-default:"5"`
	LruCapacity              int `env:"LRU_CAPACITY" env-default:"20"`
}
