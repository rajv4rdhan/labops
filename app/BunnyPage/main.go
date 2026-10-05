package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fs := http.FileServer(http.Dir("."))
	log.Printf("BunnyPage listening on :%s", port)
	if err := http.ListenAndServe(":"+port, fs); err != nil {
		log.Fatal(err)
	}
}
