package statuses

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"

	"github.com/satorunooshie/Yatter/app/handler/httperror"
	"github.com/satorunooshie/Yatter/app/handler/request"
)

// Handle request for `GET /v1/statuses/{id}`
func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := request.IDOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	ctx := r.Context()
	statusRepo := h.app.Dao.Status() // domain/repository の取得

	status, err := statusRepo.FindByID(ctx, id)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if status == nil {
		httperror.BadRequest(w, errors.New("status does not exist"))
		return
	}

	accountRepo := h.app.Dao.Account()
	account, err := accountRepo.FindByID(ctx, status.AccountID)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if account == nil {
		httperror.InternalServerError(w, errors.Errorf("account that has this status (%v) not found", status))
		return
	}

	status.Account = *account

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(status); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
