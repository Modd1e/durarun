package api

import (
	"log/slog"
	"net/http"

	"github.com/Modd1e/durarun/internal/postgres/dbgen"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type API struct {
	queries dbgen.Querier
	log     *slog.Logger
}

func New(queries dbgen.Querier, log *slog.Logger) *API {
	return &API{
		queries: queries,
		log:     log,
	}
}

func (a *API) Handler() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/jobs", func(r chi.Router) {
		r.Get("/count", a.jobsCount)
		r.Post("/", a.createJob)
	})

	return r
}
