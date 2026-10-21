package app

import (
	"finance-tracker/internal/handler"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/health", handler.Health)

	return r
}
