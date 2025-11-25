package main

import (
	ihttp "AvitoTest/internal/http"
	"AvitoTest/internal/storage"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Hello World")
	log.Println("Hello World2")
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables only")
	}
	log.Println("Attempting to connect to Postgres...")
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
	if err := http.ListenAndServe("0.0.0.0:"+port, router); err != nil {
		log.Fatal(err)
	}

}
