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
	limit, sinceID, maxID, selectType, err := h.validateQuery(r)
	if err != nil {
		httperror.BadRequest(w, err)
		return
	}

	ctx := r.Context()

	statuses := make([]*object.Status, 0, limit)
	media := make([]*object.MediaAttachment, 0, limit)

	accountIDs := make([]int64, 0, limit)
	statusIDs := make([]int64, 0, limit)

	switch selectType {
	case OnlyMedia:
		mediaRepo := h.app.Dao.MediaAttachment()
		if media, err = mediaRepo.Select(ctx, sinceID, maxID, limit); err != nil {
			httperror.InternalServerError(w, err)
			return
		}

		for _, v := range media {
			statusIDs = append(statusIDs, v.StatusID)
		}

		if len(statusIDs) != 0 {
			statusRepo := h.app.Dao.Status()
			statuses, err = statusRepo.FindByIDs(ctx, statusIDs)
			if err != nil {
				httperror.InternalServerError(w, err)
				return
			}
		}
	case All, request.ParamNotFound:
		statusRepo := h.app.Dao.Status()
		statuses, err = statusRepo.Select(ctx, sinceID, maxID, limit)
		if err != nil {
			httperror.InternalServerError(w, err)
			return
		}

		for _, v := range statuses {
			statusIDs = append(statusIDs, v.ID)
		}

		if len(statusIDs) != 0 {
			mediaRepo := h.app.Dao.MediaAttachment()
			if media, err = mediaRepo.FindByStatusIDs(ctx, statusIDs); err != nil {
				httperror.InternalServerError(w, err)
				return
			}
		}
	}

	/* MediaAttachmentをレスポンスに詰める */
	if len(statusIDs) != 0 && len(media) != 0 {
		mediaAttachmentMap := make(map[object.StatusID][]*object.MediaAttachment, len(media))
		for _, v := range media {
			mediaAttachmentMap[v.StatusID] = append(mediaAttachmentMap[v.StatusID], v)
		}

		for _, v := range statuses {
			accountIDs = append(accountIDs, v.AccountID)

			if m, ok := mediaAttachmentMap[v.ID]; ok {
				v.MediaAttachment = m
			}
		}
	}

	/* AccountIDsからAccountを取得しレスポンスに詰める */
	if len(accountIDs) != 0 {
		accountRepo := h.app.Dao.Account()
		accounts, err := accountRepo.FindByIDs(ctx, accountIDs)
		if err != nil || accounts == nil {
			httperror.InternalServerError(w, err)
			return
		}
		accountMap := make(map[object.AccountID]*object.Account, len(accounts))
		for _, v := range accounts {
			accountMap[v.ID] = v
		}
		for i, v := range statuses {
			if m, ok := accountMap[v.AccountID]; ok {
				statuses[i].Account = m
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&statuses); err != nil {
		httperror.InternalServerError(w, err)
		return
	}
}

func (h *handler) validateQuery(r *http.Request) (limit, sinceID, maxID, selectType int64, err error) {
	limit, err = request.DecodeParam2Int64(r, "limit")
	if err != nil {
		return
	}
	if limit == request.ParamNotFound || limit > 80 {
		limit = 40
	}
	maxID, err = request.DecodeParam2Int64(r, "max_id")
	if err != nil {
		return
	}
	if maxID == request.ParamNotFound {
		maxID = math.MaxInt64
	}
	sinceID, err = request.DecodeParam2Int64(r, "since_id")
	if err != nil {
		return
	}
	if sinceID == request.ParamNotFound {
		sinceID = 0
	}
	selectType, err = request.DecodeParam2Int64(r, "only_media")
	if err != nil {
		return
	}
	if selectType == request.ParamNotFound {
		selectType = All
	}
	return limit, sinceID, maxID, selectType, nil
}
