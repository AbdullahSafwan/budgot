package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"budgot/internal/configs"
	"budgot/internal/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := configs.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := router.NewRouter(db)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
