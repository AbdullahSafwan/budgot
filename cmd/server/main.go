package main

import (
	"context"
	"log"
	"net/http"

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

	if err := db.Schema.Create(context.Background()); err != nil {
		log.Fatal("failed creating schema resources: ", err)
	}
	log.Println("Database migrated successfully")

	r := router.NewRouter()
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
