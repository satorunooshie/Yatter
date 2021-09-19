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
	mediaAttachment struct {
		db *sqlx.DB
	}
)

func NewMediaAttachment(db *sqlx.DB) repository.MediaAttachment {
	return &mediaAttachment{db: db}
}

func (r *mediaAttachment) FindByStatusIDs(ctx context.Context, statusIDs []int64) ([]*object.MediaAttachment, error) {
	query, params, err := sqlx.In("SELECT * FROM `media_attachment` WHERE `status_id` IN (?) AND `delete_at` IS NULL", statusIDs)
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
			log.Printf("[WARN] dao::media_attachment::FindByIDs::rows.Close(): %v", err)
		}
	}()

	entities := make([]*object.MediaAttachment, 0, len(statusIDs))
	for rows.Next() {
		entity := &object.MediaAttachment{}
		if err := rows.StructScan(&entity); err != nil {
			return nil, err
		}
		entity.SetMediaType()
		entities = append(entities, entity)
	}
	return entities, nil
}
