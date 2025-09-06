package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"townsag/url_shortener/api/util"
)


// TODO: measure for what sizes of struct does passing the struct by reference become
//		 faster than passing the struct by value
func getConfiguration() (*pgxpool.Config, error) {
	var portEnv string = util.GetEnvWithDefault("POSTGRES_PORT", "5432")
	// declaring port inside of the if condition expression would create a
	// new version of port scoped to the if statement. This version of port
	// would shadow but not replace the function scoped version of port
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		port = 5432
	}

	host := util.GetEnvWithDefault("POSTGRES_HOST", "localhost")
	dbName := util.GetEnvWithDefault("POSTGRES_DB", "postgres")
	user := util.GetEnvWithDefault("POSTGRES_USER", "admin")
	password := util.GetEnvWithDefault("POSTGRES_PASSWORD", "password")
	poolMaxCons := util.GetEnvWithDefault("POOL_MAX_CONS", "25")

	return pgxpool.ParseConfig(fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s pool_max_cons=%s",
		host, port, user, password, dbName, poolMaxCons,
	))
}

func createDBConnectionPool(ctx context.Context, config *pgxpool.Config) (*pgxpool.Pool, error) {
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create a database connection pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping the new connection pool: %w", err)
	}
	return pool, nil
}

type redisConfig struct {
	address string
}

func getRedisConfiguration() *redisConfig {
	return &redisConfig{
		address: fmt.Sprintf(
			"%s:%s",
			util.GetEnvWithDefault("REDIS_HOST", "localhost"),
			util.GetEnvWithDefault("REDIS_PORT", "6379"),
		),
	}
}

func createRedisConnection(ctx context.Context, config *redisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: config.address,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("unable to reach redis server: %w", err)
	}
	return rdb, nil
}
