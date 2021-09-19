package repository

import (
	"context"
	"time"

	"github.com/satorunooshie/Yatter/app/domain/object"
)

type Account interface {
	FindByID(ctx context.Context, id object.AccountID) (*object.Account, error)
	FindByIDs(ctx context.Context, ids []object.AccountID) ([]*object.Account, error)
	FindByUsername(ctx context.Context, username string) (*object.Account, error)
	Insert(ctx context.Context, username, passwordHash string, createAt time.Time) error
}
