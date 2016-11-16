package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/namsral/flag"
	"log"
)

var config *Config

type Config struct {
	dbConnURL string
	listenAddr string

	DB *sqlx.DB
}

func configure() {
	parseFlags()

	db, err := sqlx.Connect("postgres", config.dbConnURL)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v\n", err)
	}
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(2)

	config.DB = config.DB
}

func getListenAddress() string {
	if config == nil {
		panic("Configuration not yet initialized")
	}

	return config.listenAddr
}

func parseFlags() {
	dbConnUrl := flag.String("db", "host=127.0.0.1 port=5432 dbname=gtest user=margus sslmode=disable", "database connection")
	listenAddr := flag.String("listen", "0.0.0.0:8080", "")

	flag.Parse()

	config = &Config{
		dbConnURL:  *dbConnUrl,
		listenAddr: *listenAddr,
	}
}
