package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tsubasa283paris/HMDC/sqlc/db"

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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		userId := r.FormValue("userId")
		log.Println("userId: ", userId)

		ctx := context.Background()

		pgUser := os.Getenv("DB_USER")
		pgPassword := os.Getenv("DB_PASSWORD")
		if pgUser == "" || pgPassword == "" {
			w.WriteHeader(http.StatusBadRequest)
			log.Println("Environment variable DB_USER or DB_PASSWORD not set")
			return
		}

		pgTarget := fmt.Sprintf("user=%s password=%s dbname=hmdc sslmode=require", pgUser, pgPassword)
		dbCnx, err := sql.Open("postgres", pgTarget)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		queries := db.New(dbCnx)

		user, err := queries.GetUser(ctx, userId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		bodyStr, err := json.Marshal(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(bodyStr))
	})

	log.Fatal(http.ListenAndServe(":"+port, r))
}
