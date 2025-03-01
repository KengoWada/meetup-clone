// Package db provides functions for establishing and managing
// database connections used throughout the application.
package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/DATA-DOG/go-txdb"
	"github.com/KengoWada/meetup-clone/internal/config"
	_ "github.com/lib/pq"
)

const PostgresDriver = "postgres"

// New creates a new database connection with the specified parameters.
// It connects to the database using the provided address, sets the maximum
// number of open and idle connections, and configures the maximum idle time.
func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime string, environment config.AppEnv) (db *sql.DB, err error) {
	if environment == config.AppEnvTest {
		db = sql.OpenDB(txdb.New(PostgresDriver, addr))
		if db == nil {
			return nil, errors.New("failed to connect to test database")
		}
	} else {
		db, err = sql.Open(PostgresDriver, addr)
		if err != nil {
			return nil, err
		}
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
