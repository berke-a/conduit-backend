package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pool *pgxpool.Pool

func init() {
	cfg, err := LoadConfig("config.yaml") // TODO: Fix config file path
	if err != nil {
		panic(fmt.Errorf("error loading config: %w", err))
	}

	fmt.Printf("Loaded config: %+v\n", cfg)

	poolConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		panic(fmt.Errorf("error parsing pool config: %w", err))
	}

	poolConfig.ConnConfig.Host = cfg.Host
	poolConfig.ConnConfig.Port = uint16(cfg.Port)
	poolConfig.ConnConfig.User = cfg.User
	poolConfig.ConnConfig.Password = cfg.Password
	poolConfig.ConnConfig.Database = cfg.DBName
	poolConfig.MaxConns = 10 // TODO: Adjust later

	pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		panic(fmt.Errorf("error creating connection pool: %w", err))
	}
}

func GetConn() *pgxpool.Pool {
	return pool
}
