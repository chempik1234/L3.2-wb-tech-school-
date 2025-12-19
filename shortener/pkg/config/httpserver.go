package config

// HTTPServerConfig is the config struct for servers (only HTTP_PORT)
type HTTPServerConfig struct {
	Port int `env:"PORT" envDefault:"8080"`
}
