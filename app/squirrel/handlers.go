package main

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//go:embed web
var webFS embed.FS

type noteRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func newRouter(store *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/ready", handleReady(store))
	mux.HandleFunc("/api/notes", handleNotes(store))
	mux.HandleFunc("/api/notes/", handleNote(store))

	web, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatalf("web assets: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(web)))

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

func handleNotes(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			notes, err := store.ListNotes(r.Context())
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to list notes")
				return
			}
			writeJSON(w, http.StatusOK, notes)
		case http.MethodPost:
			req, ok := decodeNote(w, r)
			if !ok {
				return
			}
			note, err := store.CreateNote(r.Context(), req.Title, req.Body)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "failed to create note")
				return
			}
			writeJSON(w, http.StatusCreated, note)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func handleNote(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/api/notes/"), 10, 64)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid id")
			return
		}

		switch r.Method {
		case http.MethodGet:
			note, err := store.GetNote(r.Context(), id)
			if err != nil {
				writeNoteError(w, err, "failed to get note")
				return
			}
			writeJSON(w, http.StatusOK, note)
		case http.MethodPut:
			req, ok := decodeNote(w, r)
			if !ok {
				return
			}
			note, err := store.UpdateNote(r.Context(), id, req.Title, req.Body)
			if err != nil {
				writeNoteError(w, err, "failed to update note")
				return
			}
			writeJSON(w, http.StatusOK, note)
		case http.MethodDelete:
			if err := store.DeleteNote(r.Context(), id); err != nil {
				writeNoteError(w, err, "failed to delete note")
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func decodeNote(w http.ResponseWriter, r *http.Request) (noteRequest, bool) {
	var req noteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return noteRequest{}, false
	}
	req.Title = strings.TrimSpace(req.Title)
	req.Body = strings.TrimSpace(req.Body)
	if req.Title == "" && req.Body == "" {
		writeError(w, http.StatusBadRequest, "note is empty")
		return noteRequest{}, false
	}
	return req, true
}

func writeNoteError(w http.ResponseWriter, err error, fallback string) {
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "note not found")
		return
	}
	writeError(w, http.StatusInternalServerError, fallback)
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
