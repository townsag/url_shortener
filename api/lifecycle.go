package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"

	"townsag/url_shortener/api/util"
)

type dbConfig struct {
	host     string
	port     int
	dbName   string
	user     string
	password string
}

// TODO: measure for what sizes of struct does passing the struct by reference become
//
//	faster than passing the struct by value
func getConfiguration() *dbConfig {
	var portEnv string = util.GetEnvWithDefault("POSTGRES_PORT", "5432")
	// declaring port inside of the if condition expression would create a
	// new version of port scoped to the if statement. This version of port
	// would shadow but not replace the function scoped version of port
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		port = 5432
	}

	config := dbConfig{
		host:     util.GetEnvWithDefault("POSTGRES_HOST", "localhost"),
		port:     port,
		dbName:   util.GetEnvWithDefault("POSTGRES_DB", "postgres"),
		user:     util.GetEnvWithDefault("POSTGRES_USER", "admin"),
		password: util.GetEnvWithDefault("POSTGRES_PASSWORD", "password"),
	}
	return &config
}

func createDBConnection(ctx context.Context, config *dbConfig) (*pgx.Conn, error) {
	conn, err := pgx.Connect(
		ctx,
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s",
			config.host, config.port, config.user, config.password, config.dbName,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close(ctx)
		return nil, fmt.Errorf("failed to ping database after creating connection: %w", err)
	}

	return conn, nil
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
