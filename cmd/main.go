package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
