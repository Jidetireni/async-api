package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Jidetireni/async-api.git/config"
	_ "github.com/lib/pq"
)

func DbInit(c *config.Config) (*sql.DB, error) {
	url := c.DbUrl()
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	return db, nil
}
