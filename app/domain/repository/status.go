package repository

import (
	"context"

	"github.com/satorunooshie/Yatter/app/domain/object"
)

type Status interface {
	FindByID(ctx context.Context, id int64) (*object.Status, error)
	Select(ctx context.Context, minID, maxID, limit int64) ([]*object.Status, error)
	SelectOnlyMedia(ctx context.Context, minID, maxID, limit int64) ([]*object.Status, error)
	Insert(ctx context.Context, accountID int64, content string) (int64, error)
	Delete(ctx context.Context, id int64) error
}
