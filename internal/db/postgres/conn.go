package postgres

import (
	"context"
	"net/url"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(cfg Config) (*pgxpool.Pool, error) {
	uri := url.URL{
		Scheme: "postgres",
		Host:   cfg.Host,
		User:   url.UserPassword(cfg.User, cfg.Password),
		Path:   cfg.DB_Name,
	}

	pgConfig, err := pgxpool.ParseConfig(uri.String())

	if err != nil {
		return nil, err
	}

	pgConfig.MaxConns = cfg.MaxConns
	pgConfig.MinConns = cfg.MinConns
	pgConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	pgConfig.MaxConnLifetime = cfg.MaxConnLifetime

	pool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)

	if err != nil {
		return nil, err
	}

	return pool, nil
}
