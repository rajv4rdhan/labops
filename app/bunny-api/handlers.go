package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Bunny struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Color     string `json:"color"`
	CreatedAt string `json:"created_at"`
}

type createBunnyRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func newRouter(store *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/ready", handleReady(store))
	mux.HandleFunc("/api/bunnies", handleBunnies(store))
	mux.HandleFunc("/api/bunnies/", handleBunny(store))
	return withLogging(mux)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleReady(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := store.Ping(r.Context()); err != nil {
			writeError(w, http.StatusServiceUnavailable, "database not ready")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}

func handleBunnies(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			bunnies, err := store.ListBunnies(r.Context())
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to list bunnies")
				return
			}
			writeJSON(w, http.StatusOK, bunnies)
		case http.MethodPost:
			var req createBunnyRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeError(w, http.StatusBadRequest, "invalid json body")
				return
			}
			req.Name = strings.TrimSpace(req.Name)
			if req.Name == "" {
				writeError(w, http.StatusBadRequest, "name is required")
				return
			}
			bunny, err := store.CreateBunny(r.Context(), req)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to create bunny")
				return
			}
			writeJSON(w, http.StatusCreated, bunny)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func handleBunny(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/api/bunnies/"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}
		bunny, err := store.GetBunny(r.Context(), id)
		if err != nil {
			writeError(w, http.StatusNotFound, "bunny not found")
			return
		}
		writeJSON(w, http.StatusOK, bunny)
	}
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
