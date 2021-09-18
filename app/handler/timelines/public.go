package timelines

import (
	"encoding/json"
	"math"
	"net/http"

	"github.com/satorunooshie/Yatter/app/domain/object"
	"github.com/satorunooshie/Yatter/app/handler/httperror"
	"github.com/satorunooshie/Yatter/app/handler/request"
)

const (
	All = iota
	OnlyMedia
)

func (h *handler) GetPublic(w http.ResponseWriter, r *http.Request) {
	limit, err := request.DecodeParam2Int64(r, "limit")
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	if limit == request.ParamNotFound || limit > 80 {
		limit = 40
	}
	maxID, err := request.DecodeParam2Int64(r, "max_id")
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	if maxID == request.ParamNotFound {
		maxID = math.MaxInt64
	}
	sinceID, err := request.DecodeParam2Int64(r, "since_id")
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}
	if sinceID == request.ParamNotFound {
		sinceID = 0
	}
	selectType, err := request.DecodeParam2Int64(r, "only_media")
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	ctx := r.Context()
	statusRepo := h.app.Dao.Status()

	statuses := make([]*object.Status, 0, limit)
	switch selectType {
	case OnlyMedia:
		statuses, err = statusRepo.SelectOnlyMedia(ctx, sinceID, maxID, limit)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}
	case All, request.ParamNotFound:
		statuses, err = statusRepo.Select(ctx, sinceID, maxID, limit)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}
	}

	accountIDs := make([]int64, len(statuses))
	for i, v := range statuses {
		accountIDs[i] = v.AccountID
	}

	accountRepo := h.app.Dao.Account()
	accounts, err := accountRepo.FindByIDs(ctx, accountIDs)
	for i, v := range accounts {
		statuses[i].Account = *v
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&statuses); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}
