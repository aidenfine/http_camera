package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func StreamRouter() chi.Router {
	r := chi.NewRouter()

	r.Get("/", stream())

	return r
}

func stream() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode("Hello world")
	}
}
