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
	// Implementation for repository.Status
	status struct {
		db *sqlx.DB
	}
)

func NewStatus(db *sqlx.DB) repository.Status {
	return &status{db: db}
}

func (r *status) FindByID(ctx context.Context, id int64) (*object.Status, error) {
	entity := &object.Status{}
	if err := r.db.QueryRowxContext(ctx, "SELECT * FROM `status` WHERE `id` = ? AND `delete_at` IS NULL", id).StructScan(entity); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return entity, nil
}

func (r *status) Insert(ctx context.Context, accountID int64, content string) (int64, error) {
	stmt, err := r.db.PreparexContext(ctx, "INSERT INTO `status` (`account_id`, `content`) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("[WARN] dao::status::Insert::stmt.Close(): %v", err)
		}
	}()

	res, err := stmt.ExecContext(ctx, accountID, content)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *status) Delete(ctx context.Context, id int64) error {
	stmt, err := r.db.PrepareContext(ctx, "UPDATE `status` SET `delete_at` = NOW() WHERE `id` = ?")
	if err != nil {
		return err
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("[WARN] dao::status::Delete::stmt.Close(): %v", err)
		}
	}()

	if _, err := stmt.ExecContext(ctx, id); err != nil {
		return err
	}
	return nil
}
