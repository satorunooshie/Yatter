package accounts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/satorunooshie/Yatter/app/handler/httperror"
)

// Handle request for `GET /v1/accounts/{username}`
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		httperror.BadRequest(w, errors.New("invalid params"))
		return
	}

	ctx := r.Context()
	accountRepo := h.app.Dao.Account() // domain/repository の取得

	user, err := accountRepo.FindByUsername(ctx, username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if user == nil {
		httperror.BadRequest(w, errors.New("account does not exist"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&user); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
