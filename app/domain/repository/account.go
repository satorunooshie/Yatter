package repository

import (
	"context"

	"github.com/satorunooshie/Yatter/app/domain/object"
)

type Account interface {
	FindByID(ctx context.Context, id int64) (*object.Account, error)
	FindByIDs(ctx context.Context, ids []int64) ([]*object.Account, error)
	FindByUsername(ctx context.Context, username string) (*object.Account, error)
	Insert(ctx context.Context, username, passwordHash string) error
}
