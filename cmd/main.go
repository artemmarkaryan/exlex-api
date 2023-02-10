package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/artemmarkaryan/exlex-backend/internal/migrations"
	"github.com/artemmarkaryan/exlex-backend/internal/server"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("PSQL_DSN"))
	if err != nil {
		log.Fatalln(err)
	}

	if err = goose.Up(db, "/"); err != nil {
		log.Fatalln(err)
	}

	server.Serve()
}
