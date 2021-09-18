package accounts

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/satorunooshie/Yatter/app/domain/object"
	"github.com/satorunooshie/Yatter/app/handler/httperror"
)

// Request body for `POST /v1/accounts`
type AccountCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Handle request for `POST /v1/accounts`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req AccountCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	account := new(object.Account)
	account.Username = req.Username
	if err := account.SetPassword(req.Password); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	ctx := r.Context()
	accountRepo := h.app.Dao.Account() // domain/repository の取得

	user, err := accountRepo.FindByUsername(ctx, account.Username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if user != nil {
		httperror.BadRequest(w, errors.New("account name is already in use"))
		return
	}

	if err := accountRepo.Insert(ctx, account.Username, account.PasswordHash); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
