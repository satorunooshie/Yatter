package dao

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

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

func (r *account) FindByID(ctx context.Context, id object.AccountID) (*object.Account, error) {
	entity := &object.Account{}
	if err := r.db.QueryRowxContext(ctx, "SELECT * FROM `account` WHERE `id` = ? AND `delete_at` IS NULL", id).StructScan(entity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return entity, nil
}

func (r *account) FindByIDs(ctx context.Context, ids []object.AccountID) ([]*object.Account, error) {
	query, params, err := sqlx.In("SELECT * FROM `account` WHERE `id` IN (?) AND `delete_at` IS NULL", ids)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryxContext(ctx, query, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("[WARN] dao::account::FindByIDs::rows.Close(): %v", err)
		}
	}()

	entities := make([]*object.Account, 0, len(ids))
	for rows.Next() {
		entity := &object.Account{}
		if err := rows.StructScan(&entity); err != nil {
			return nil, err
		}
		entities = append(entities, entity)
	}
	return entities, nil
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

func (r *account) Insert(ctx context.Context, username, passwordHash string, createAt time.Time) error {
	stmt, err := r.db.PreparexContext(ctx, "INSERT INTO `account` (`username`, `password_hash`, `create_at`) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("[WARN] dao::account::Insert::stmt.Close(): %v", err)
		}
	}()
	if _, err := stmt.ExecContext(ctx, username, passwordHash, createAt); err != nil {
		return err
	}
	return nil
}
