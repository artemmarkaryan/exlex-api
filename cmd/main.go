package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/artemmarkaryan/exlex-backend/internal/migrations"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func main() {
	godotenv.Load(".env")
	godotenv.Load("/etc/exlex/.env")

	db, err := sql.Open("postgres", os.Getenv("PSQL_DSN"))
	if err != nil {
		log.Fatalln(err)
	}

	if err = goose.Up(db, "/"); err != nil {
		log.Fatalln(err)
	}

	log.Print("Serving...")
	http.HandleFunc("/live", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusTeapot)
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
