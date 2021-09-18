package dao

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/satorunooshie/Yatter/app/domain/object"
	"github.com/satorunooshie/Yatter/app/domain/repository"

	"github.com/jmoiron/sqlx"
)

type (
	// Implementation for repository.Account
	account struct {
		db *sqlx.DB
	}
)

// Create accout repository
func NewAccount(db *sqlx.DB) repository.Account {
	return &account{db: db}
}

// FindByUsername : ユーザ名からユーザを取得
func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := &object.Account{}
	if err := r.db.QueryRowxContext(ctx, "SELECT * FROM `account` WHERE `username` = ?", username).StructScan(entity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return entity, nil
}

func (r *account) Insert(ctx context.Context, username, passwordHash string) error {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO `account` (`username`, `password_hash`) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("[WARN] dao::account::Insert::stmt.Close(): %v", err)
		}
	}()
	if _, err := stmt.Exec(username, passwordHash); err != nil {
		return err
	}
	return nil
}
