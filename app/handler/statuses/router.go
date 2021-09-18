package statuses

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/satorunooshie/Yatter/app/app"
	"github.com/satorunooshie/Yatter/app/handler/auth"
)

// Implementation of handler
type handler struct {
	app *app.App
}

// Create Handler for `/v1/statuses/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}
	r.With(auth.Middleware(app)).Post("/", h.Create)
	r.Get("/{id}", h.Get)
	r.With(auth.Middleware(app)).Delete("/{id}", h.Delete)

	return r
}
