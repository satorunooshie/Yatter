package handler

import (
	"github.com/satorunooshie/Yatter/app/handler/auth"
	"github.com/satorunooshie/Yatter/app/handler/statuses"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/satorunooshie/Yatter/app/app"
	"github.com/satorunooshie/Yatter/app/handler/accounts"
	"github.com/satorunooshie/Yatter/app/handler/health"
)

func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(newCORS().Handler)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/v1/health", health.NewRouter())

	r.Mount("/v1/accounts", accounts.NewRouter(app))

	/* Auth */
	r.Mount("/v1/statuses", auth.Middleware(app)(statuses.NewRouter(app)))

	return r
}

func newCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
	})
}
