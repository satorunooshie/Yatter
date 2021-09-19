package repository

import (
	"context"

	"github.com/satorunooshie/Yatter/app/domain/object"
)

type MediaAttachment interface {
	FindByStatusIDs(ctx context.Context, statusIDs []int64) ([]*object.MediaAttachment, error)
}
