package accounts

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

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
	if err := account.SetPassword(req.Password); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	ctx := r.Context()
	accountRepo := h.app.Dao.Account() // domain/repository の取得

	accountInUse, err := accountRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if accountInUse != nil {
		httperror.BadRequest(w, errors.New("account name is already in use"))
		return
	}

	account.Username = req.Username
	account.CreateAt = object.DateTime{Time: time.Now()}
	if err := accountRepo.Insert(ctx, account.Username, account.PasswordHash, account.CreateAt.Time); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(account); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
