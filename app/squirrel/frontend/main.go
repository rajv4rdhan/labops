package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	port := getenv("PORT", "8080")

	log.Printf("squirrel-frontend listening on :%s", port)
	if err := http.ListenAndServe(":"+port, newHandler()); err != nil {
		log.Fatal(err)
	}
}

func newHandler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	if apiURL := os.Getenv("API_URL"); apiURL != "" {
		target, err := url.Parse(apiURL)
		if err != nil {
			log.Fatalf("invalid API_URL %q: %v", apiURL, err)
		}
		mux.Handle("/api/", httputil.NewSingleHostReverseProxy(target))
		log.Printf("proxying /api/ to %s", target)
	}

	mux.Handle("/", http.FileServer(http.Dir(".")))
	return mux
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
