package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("BunnyPage listening on :%s", port)
	if err := http.ListenAndServe(":"+port, newHandler()); err != nil {
		log.Fatal(err)
	}
}

func newHandler() http.Handler {
	mux := http.NewServeMux()

	if apiURL := os.Getenv("API_URL"); apiURL != "" {
		target, err := url.Parse(apiURL)
		if err != nil {
			log.Fatalf("invalid API_URL %q: %v", apiURL, err)
		}
		proxy := httputil.NewSingleHostReverseProxy(target)
		mux.Handle("/api/", proxy)
		log.Printf("proxying /api/ to %s", target)
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))
	return mux
}
