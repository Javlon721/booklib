package postgres

import "time"

type Config struct {
	DB_Name  string
	Password string
	User     string
	Host     string

	MaxConns        int32
	MinConns        int32
	MaxConnIdleTime time.Duration
	MaxConnLifetime time.Duration
}
