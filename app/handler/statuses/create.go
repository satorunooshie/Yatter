package statuses

import (
	"encoding/json"
	"net/http"

	"github.com/satorunooshie/Yatter/app/domain/object"
	"github.com/satorunooshie/Yatter/app/handler/auth"
	"github.com/satorunooshie/Yatter/app/handler/httperror"
)

// Request body for `POST /v1/statuses`
type StatusCreateRequest struct {
	Status   string  `json:"status"`
	MediaIDs []int64 `json:"media_ids"`
}

// Handle request for `POST /v1/statuses`
func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	account := auth.AccountOf(r)
	if account == nil {
		httperror.Error(w, http.StatusUnauthorized)
		return
	}

	var req StatusCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httperror.BadRequest(w, err)
		return
	}

	ctx := r.Context()
	statusRepo := h.app.Dao.Status() // domain/repository の取得

	id, err := statusRepo.Insert(ctx, account.ID, req.Status)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	status, err := statusRepo.FindByID(ctx, id)
	if err != nil {
		httperror.InternalServerError(w, err)
		return
	}

	res := &object.Status{
		ID:       id,
		Account:  *account,
		Content:  status.Content,
		CreateAt: status.CreateAt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
