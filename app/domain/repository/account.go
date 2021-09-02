package repository

import (
	"context"

	"github.com/satorunooshie/Yatter/app/domain/object"
)

type Account interface {
	// Fetch account which has specified username
	FindByUsername(ctx context.Context, username string) (*object.Account, error)
	// TODO: Add Other APIs
}
