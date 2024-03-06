package db

import (
	"context"
	"fmt"

	"./config.go"

	"github.com/jackc/pgx/v5" // Update version number
)

var pool *pgx.ConnPool

func init() {
	// Load configuration from a separate file

	cfg, err := config.LoadConfig("config.yaml") // Replace with your config file path
	if err != nil {
		panic(fmt.Errorf("error loading config: %w", err))
	}

	poolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			Database: cfg.Database.DBName,
		},
		MaxConnections: 10, // Adjust as needed
	}

	var err error
	pool, err = pgx.NewConnPool(context.Background(), poolConfig)
	if err != nil {
		panic(fmt.Errorf("error creating connection pool: %w", err))
	}
}

// GetConn returns a connection from the pool
func GetConn() *pgx.ConnPool {
	return pool
}
