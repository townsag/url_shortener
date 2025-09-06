package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testPool *pgxpool.Pool
	pgContainer *postgres.PostgresContainer
	setupOncePG sync.Once
	cleanupOncePG sync.Once
	dbr *redis.Client
	redisContainer *tcredis.RedisContainer
	setupOnceRedis sync.Once
	cleanupOnceRedis sync.Once
)

/*
Goals:
- one postgres container that exists for the lifetime of the tests for the
  handlers package
- container is created at the beginning of testing
- container is destroyed at the end of testing
TODO:
- look into creating database snapshots that can be restored before each test
	- https://golang.testcontainers.org/modules/postgres/#using-snapshots
*/

func setupPostgresContainer() (*pgxpool.Pool, error) {
	var err error = nil
	setupOncePG.Do(
		func() {
			ctx := context.Background()
			fmt.Println("creating postgres container")
			pgContainer, err = postgres.Run(
				ctx,
				"postgres:17-alpine",
				postgres.WithInitScripts(filepath.Join("..", "sql", "schema.sql")),
				postgres.WithDatabase("testing"),
				postgres.WithUsername("testing"),
				postgres.WithPassword("testing"),
				postgres.BasicWaitStrategies(),
			)
			if err != nil {
				err = fmt.Errorf("failed to start testing postgres container: %w", err)
				return
			}
			// create a connection to the postgres test container
			fmt.Println("creating connection to postgres container")
			var dbURL string
			dbURL, err = pgContainer.ConnectionString(ctx)
			if err!= nil {
				err = fmt.Errorf("unable to connect to postgres container %w", err)
				return
			}
			testPool, err = pgxpool.New(ctx, dbURL)
			if err != nil {
				err = fmt.Errorf("unable to connect to postgres container: %w", err)
				return
			}
 			// TODO: create a snapshot of the empty database
			// fmt.Println("creating snapshot of postgres database")
			// err = pgContainer.Snapshot(ctx)
			// if err != nil {
			// 	err = fmt.Errorf("unable to make a snapshot of the pg database: %w", err)
			// 	return 
			// }
		},
	)

	return testPool, err
}

func cleanupPostgresContainer() error {
	var err error = nil
	cleanupOncePG.Do(
		func() {
			if testPool != nil{
				testPool.Close()
			}
			if pgContainer != nil {
				fmt.Println("closing postgres testcontainer")
				err = testcontainers.TerminateContainer(pgContainer)
				if err != nil {
					err = fmt.Errorf("unable to cleanup postgres container: %w", err)
				}
			}
		},
	)
	return err
}

// TODO:
// func cleanupPostgresData() error {
// 	ctx := context.Background()
// 	return pgContainer.Restore(ctx)
// }

func setupRedisContainer() (*redis.Client, error) {
	var err error = nil
	setupOnceRedis.Do(
		func() {
			// create a redis container and make a connection to it
			ctx := context.Background()
			fmt.Println("creating redis container")
			redisContainer, err = tcredis.Run(
				ctx,
				"redis:latest",
				tcredis.WithConfigFile(filepath.Join("..", "..", "redis", "redis.conf")),
			)
			if err != nil {
				err = fmt.Errorf("failed to start redis test container: %w", err)
				return
			}
			// create the redis client
			var uri string
			uri, err = redisContainer.ConnectionString(ctx)
			if err != nil {
				err = fmt.Errorf("failed to get a connection string from the redis container: %w", err)
				return
			}
			var opt *redis.Options
			opt, err = redis.ParseURL(uri)
			if err != nil {
				err = fmt.Errorf("failed to create client options from connection string: %w", err)
			}
			dbr = redis.NewClient(opt)
			// TODO: ping the client here
			// return the redis client and any errors that are encountered
		},
	)
	return dbr, err
}

func cleanupRedisContainer() error {
	var err error = nil
	cleanupOnceRedis.Do(
		func() {
			if dbr != nil {
				dbr.Close()
			}
			if redisContainer != nil {
				fmt.Println("cleaning up redis container")
				err = testcontainers.TerminateContainer(redisContainer)
				if err != nil {
					err = fmt.Errorf("unable to clean up redis container: %w", err)
				}
			}
		},
	)
	return err 
}

/*
CHECKPOINT:
  - you were in the middle of adding test automation for the redis client in the healthy route and the redirect to
	long url route
*/