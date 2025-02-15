package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Jidetireni/async-api.git/config"
	_ "github.com/lib/pq"
)

// This function initializes and returns a connection to a PostgreSQL database.
// It takes a Config object (c *config.Config) as input, which contains the database configuration (e.g., host, port, username, etc.).

func DbInit(c *config.Config) (*sql.DB, error) {
	url := c.DbUrl()                     // to get the database connection URL
	db, err := sql.Open("postgres", url) // to create a connection to the PostgreSQL database.
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection %v", err)
	}

	// A context with a 5-second timeout is created to avoid hanging indefinitely if the database is unreachable.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx) // To ensure the connection is valid, the function pings the database
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	return db, nil
}
