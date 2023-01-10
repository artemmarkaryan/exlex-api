package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	log.Print("Serving...")
	http.HandleFunc("/live", func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusTeapot)
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
