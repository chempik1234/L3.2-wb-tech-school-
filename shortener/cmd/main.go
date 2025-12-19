package main

import (
	"context"
	"fmt"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/adapters/storage"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/config"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/models"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/service"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/internal/transport"
	"github.com/chempik1234/L3.2-wb-tech-school-/shortener/pkg/adapters_wbf/cache"
	"github.com/chempik1234/super-danis-library-golang/pkg/postgres"
	"github.com/chempik1234/super-danis-library-golang/pkg/server"
	"github.com/chempik1234/super-danis-library-golang/pkg/server/httpserver"
	"github.com/chempik1234/super-danis-library-golang/pkg/services"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	log.Println("starting shortener service (main.go:27)")

	//region load config from env
	cfg, err := config.NewAppConfig("", "") // /app/config.yaml
	if err != nil {
		log.Fatal(fmt.Errorf("error loading config: %w", err))
	}
	//endregion

	//region init zlog.Logger with given LogLevel
	zlog.InitConsole()
	err = zlog.SetLevel(cfg.LogConfig.LogLevel)
	if err != nil {
		zlog.Logger.Fatal().Err(fmt.Errorf("error setting log level to '%s': %w", cfg.LogConfig.LogLevel, err))
	}
	//endregion

	//region retry (define first for later postgres, rabbitmq, redis connections)
	postgresRetryStrategy := cfg.PostgresRetryConfig.ToStrategy()
	redisRetryStrategy := cfg.RedisRetryConfig.ToStrategy()

	zlog.Logger.Info().Msg("retry policies created")
	//endregion

	//region postgres
	var postgresDB *dbpg.DB

	// connect to postgres with retry
	err = retry.Do(
		func() error {
			var postgresConnErr error

			postgresDB, postgresConnErr = dbpg.New(

				cfg.PostgresConfig.MasterDSN,
				cfg.PostgresConfig.SlaveDSNs,

				&dbpg.Options{
					MaxOpenConns:    cfg.PostgresConfig.MaxOpenConnections,
					MaxIdleConns:    cfg.PostgresConfig.MaxIdleConnections,
					ConnMaxLifetime: time.Duration(cfg.PostgresConfig.ConnectionMaxLifetimeSeconds) * time.Second,
				})

			return postgresConnErr
		},
		postgresRetryStrategy)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("couldn't create postgres balancer")
	}

	zlog.Logger.Info().Msg("postgres balancer created")

	fmt.Println("postgres stub", postgresDB)

	migrationsPath := "file:///app/db/migrations"

	err = postgres.MigrateUp(cfg.PostgresConfig.MasterDSN, migrationsPath)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("couldn't migrate postgres on master DSN")
	}
	for i, dsn := range cfg.PostgresConfig.SlaveDSNs {
		if len(dsn) == 0 {
			continue
		}
		err = postgres.MigrateUp(dsn, migrationsPath)
		if err != nil {
			zlog.Logger.Fatal().Err(err).Int("dsn_index", i).Msg("couldn't migrate postgres on slave DSN")
		}
	}
	//endregion

	//region redis
	redisClient := redis.New(
		cfg.RedisConfig.Addr,
		cfg.RedisConfig.Password,
		cfg.RedisConfig.DB,
	)
	zlog.Logger.Info().Msg("redis created")
	//endregion

	//region services
	shortenerStorageRepository := storage.NewShortenerStorageInMemoryRepo()
	cacheService := services.NewCachePopularService[string, models.Link](
		cfg.CacheConfig.MinRequestsBeforeCaching,
		cfg.CacheConfig.LruCapacity,
		cache.NewRedisWBFCache[string, models.Link](redisClient, redisRetryStrategy),
	)
	shortenerService := service.NewShortenerService(shortenerStorageRepository, cacheService)
	//endregion

	ctx, stopCtx := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	//region run in background
	wg.Add(1)
	go func(wg *sync.WaitGroup, ctx2 context.Context) {
		defer wg.Done()
		//TODO: some background service run
	}(wg, ctx)
	//endregion

	//region Start HTTP
	httpHandler := transport.NewShortenerHandler(shortenerService)
	appRouter := transport.AssembleRouter(httpHandler)

	// this VVV is work of art, but with [*http.Server]
	appServer := server.NewGracefulServer[*http.Server](
		httpserver.NewGracefulServerImplementationHTTP(appRouter),
	)

	zlog.Logger.Info().Int("http_port", cfg.HTTPServerConfig.Port).Msg("server starting :http_port")

	err = appServer.GracefulRun(ctx, cfg.HTTPServerConfig.Port)
	//endregion

	//region shutdown
	if err != nil {
		zlog.Logger.Error().Msg(fmt.Errorf("http server error: %w", err).Error())
	}

	zlog.Logger.Info().Msg("server gracefully stopped")

	stopCtx()
	wg.Wait()
	zlog.Logger.Info().Msg("background operations gracefully stopped")
	//endregion
}
