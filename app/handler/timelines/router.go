package timelines

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/satorunooshie/Yatter/app/app"
)

// Implementation of handler
type handler struct {
	app *app.App
}

// Create Handler for `/v1/timelines/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}
	r.Get("/public", h.GetPublic)

	return r
}