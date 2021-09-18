package dao

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/satorunooshie/Yatter/app/domain/object"
	"github.com/satorunooshie/Yatter/app/domain/repository"
)

type (
	// Implementation for repository.Account
	account struct {
		db *sqlx.DB
	}
)

func NewAccount(db *sqlx.DB) repository.Account {
	return &account{db: db}
}

func (r *account) FindByID(ctx context.Context, id int64) (*object.Account, error) {
	entity := &object.Account{}
	if err := r.db.QueryRowxContext(ctx, "SELECT * FROM `account` WHERE `id` = ? AND `delete_at` IS NULL", id).StructScan(entity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return entity, nil
}

func (r *account) FindByUsername(ctx context.Context, username string) (*object.Account, error) {
	entity := &object.Account{}
	if err := r.db.QueryRowxContext(ctx, "SELECT * FROM `account` WHERE `username` = ? AND `delete_at` IS NULL", username).StructScan(entity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return entity, nil
}

func (r *account) Insert(ctx context.Context, username, passwordHash string) error {
	stmt, err := r.db.PreparexContext(ctx, "INSERT INTO `account` (`username`, `password_hash`) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("[WARN] dao::account::Insert::stmt.Close(): %v", err)
		}
	}()
	if _, err := stmt.ExecContext(ctx, username, passwordHash); err != nil {
		return err
	}
	return nil
}
