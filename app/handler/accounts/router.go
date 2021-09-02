package accounts

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/satorunooshie/Yatter/app/app"
)

// Implementation of handler
type handler struct {
	app *app.App
}

// Create Handler for `/v1/accounts/`
func NewRouter(app *app.App) http.Handler {
	r := chi.NewRouter()

	h := &handler{app: app}
	r.Post("/", h.Create)

	return r
}
