package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/artemmarkaryan/exlex-backend/internal/migrations"
	_ "github.com/joho/godotenv/autoload"
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

	log.Print("Serving...")
	http.HandleFunc("/live", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusTeapot)
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
