package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tsubasa283paris/HMDC/api"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

func main() {
	port := "8080"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	log.Printf("Starting up on http://localhost:%s", port)

	r := chi.NewRouter()

	r.Get("/users", api.GetUsers)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
