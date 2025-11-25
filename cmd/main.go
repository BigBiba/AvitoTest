package main

import (
	ihttp "AvitoTest/internal/http"
	"AvitoTest/internal/storage"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Hello World")
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables only")
	}

	db, err := storage.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	server := ihttp.NewServer(db)
	router := ihttp.NewRouter(server)
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}

}
