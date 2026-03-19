package main

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	conf := struct {
		//Database
		Postgres struct {
			DB_Name  string
			Password string
			User     string
			Host     string
			Port     string

			MaxConns        int32         `envconfig:"MAX_CONNS"`
			MinConns        int32         `envconfig:"MIN_CONNS"`
			MaxConnIdleTime time.Duration `envconfig:"MAX_CONN_IDLE_TIME"`
			MaxConnLifetime time.Duration `envconfig:"MAX_CONN_LIFETIME"`
		}
	}{}

	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	err = envconfig.Process("", &conf)

	if err != nil {
		panic(err)
	}

	fmt.Println(conf)
}
