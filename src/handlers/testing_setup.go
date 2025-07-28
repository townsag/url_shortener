package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/jackc/pgx/v5"
)

var (
	conn *pgx.Conn
	pgContainer *postgres.PostgresContainer
	setupOnce sync.Once
	cleanupOnce sync.Once
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

func setupPostgresContainer() (*pgx.Conn, error) {
	var err error = nil
	setupOnce.Do(
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
			conn, err = pgx.Connect(ctx, dbURL)
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

	return conn, err
}

func cleanupPostgresContainer() error {
	var err error = nil
	cleanupOnce.Do(
		func() {
			ctx := context.Background()
			if conn != nil{
				_ = conn.Close(ctx)
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