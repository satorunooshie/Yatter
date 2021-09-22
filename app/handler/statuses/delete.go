package statuses

import (
	"encoding/json"
	"net/http"

	"github.com/satorunooshie/Yatter/app/handler/auth"
	"github.com/satorunooshie/Yatter/app/handler/httperror"
	"github.com/satorunooshie/Yatter/app/handler/request"
)

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := request.IDOf(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	ctx := r.Context()
	statusRepo := h.app.Dao.Status() // domain/repository の取得

	account := auth.AccountOf(r)
	if account == nil {
		httperror.Error(w, http.StatusUnauthorized)
		return
	}

	status, err := statusRepo.FindByID(ctx, id)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}
	if status.Account.ID != account.ID {
		httperror.Error(w, http.StatusUnauthorized)
		return
	}

	if err := statusRepo.Delete(ctx, id, account.ID); err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&struct{}{}); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
