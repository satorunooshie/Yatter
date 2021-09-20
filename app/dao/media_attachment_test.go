package dao

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	sqlxmock "github.com/zhashkevych/go-sqlxmock"

	"github.com/satorunooshie/Yatter/app/domain/object"
)

func Test_mediaAttachment_FindByStatusIDs(t *testing.T) {
	createAt, _ := time.Parse("2006-01-02", "2020-01-01")
	mediaType := "image"

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	r := &mediaAttachment{
		db: db,
	}

	type args struct {
		ctx       context.Context
		statusIDs []object.StatusID
	}
	tests := []struct {
		name    string
		query   func(s sqlxmock.Sqlmock)
		args    args
		want    []*object.MediaAttachment
		wantErr bool
	}{
		{
			name: "ok",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `media_attachment` WHERE `status_id` IN (?, ?, ?) AND `delete_at` IS NULL")).
					WithArgs(1, 2, 3).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"status_id",
								"type",
								"url",
								"description",
								"create_at",
								"delete_at",
							},
						).
							AddRow(1, 1, object.TypeImage, "http://example.com/1", "description1", createAt, nil).
							AddRow(2, 2, object.TypeImage, "http://example.com/2", "description2", createAt, nil).
							AddRow(3, 3, object.TypeImage, "http://example.com/3", "description3", createAt, nil),
					)
			},
			args: args{
				ctx:       context.Background(),
				statusIDs: []object.StatusID{1, 2, 3},
			},
			want: []*object.MediaAttachment{
				{
					ID:          1,
					StatusID:    1,
					Type:        object.TypeImage,
					URL:         "http://example.com/1",
					Description: "description1",
					CreateAt:    object.DateTime{Time: createAt},
					DeleteAt:    nil,

					MediaType: &mediaType,
				},
				{
					ID:          2,
					StatusID:    2,
					Type:        object.TypeImage,
					URL:         "http://example.com/2",
					Description: "description2",
					CreateAt:    object.DateTime{Time: createAt},
					DeleteAt:    nil,

					MediaType: &mediaType,
				},
				{
					ID:          3,
					StatusID:    3,
					Type:        object.TypeImage,
					URL:         "http://example.com/3",
					Description: "description3",
					CreateAt:    object.DateTime{Time: createAt},
					DeleteAt:    nil,

					MediaType: &mediaType,
				},
			},
			wantErr: false,
		},
		{
			name: "no rows",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `media_attachment` WHERE `status_id` IN (?, ?, ?) AND `delete_at` IS NULL")).
					WithArgs(1, 2, 3).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"status_id",
								"type",
								"url",
								"description",
								"create_at",
								"delete_at",
							},
						),
					)
			},
			args: args{
				ctx:       context.Background(),
				statusIDs: []object.StatusID{1, 2, 3},
			},
			want:    []*object.MediaAttachment{},
			wantErr: false,
		},
		{
			name: "error",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `media_attachment` WHERE `status_id` IN (?, ?, ?) AND `delete_at` IS NULL")).
					WithArgs(1, 2, 3).
					WillReturnError(errors.New("error"))
			},
			args: args{
				ctx:       context.Background(),
				statusIDs: []object.StatusID{1, 2, 3},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.query(mock)
			got, err := r.FindByStatusIDs(tt.args.ctx, tt.args.statusIDs)
			if (err != nil) != tt.wantErr {
				t.Errorf("mediaAttachment.FindByStatusIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("mediaAttachment.FindByStatusIDs() returned diff (want -> got):\n%s", diff)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func Test_mediaAttachment_Select(t *testing.T) {
	mediaType := "image"
	createAt, _ := time.Parse("2012-01-02", "2020-01-01")

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	r := &mediaAttachment{
		db: db,
	}

	type args struct {
		ctx   context.Context
		minID object.StatusID
		maxID object.StatusID
		limit int64
	}
	tests := []struct {
		name    string
		query   func(s sqlxmock.Sqlmock)
		args    args
		want    []*object.MediaAttachment
		wantErr bool
	}{
		{
			name: "ok",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `media_attachment` WHERE `status_id` BETWEEN ? AND ? AND `delete_at` IS NULL ORDER BY `create_at` DESC LIMIT ?")).
					WithArgs(1, 100, 2).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"status_id",
								"type",
								"url",
								"description",
								"create_at",
								"delete_at",
							},
						).
							AddRow(1, 1, object.TypeImage, "http://example.com/1", "description1", createAt, nil).
							AddRow(2, 2, object.TypeImage, "http://example.com/2", "description2", createAt, nil),
					)
			},
			args: args{
				ctx:   context.Background(),
				minID: 1,
				maxID: 100,
				limit: 2,
			},
			want: []*object.MediaAttachment{
				{
					ID:          1,
					StatusID:    1,
					Type:        object.TypeImage,
					URL:         "http://example.com/1",
					Description: "description1",
					CreateAt:    object.DateTime{Time: createAt},
					DeleteAt:    nil,

					MediaType: &mediaType,
				},
				{
					ID:          2,
					StatusID:    2,
					Type:        object.TypeImage,
					URL:         "http://example.com/2",
					Description: "description2",
					CreateAt:    object.DateTime{Time: createAt},
					DeleteAt:    nil,

					MediaType: &mediaType,
				},
			},
		},
		{
			name: "no rows",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `media_attachment` WHERE `status_id` BETWEEN ? AND ? AND `delete_at` IS NULL ORDER BY `create_at` DESC LIMIT ?")).
					WithArgs(1, 100, 2).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"status_id",
								"type",
								"url",
								"description",
								"create_at",
								"delete_at",
							},
						),
					)
			},
			args: args{
				ctx:   context.Background(),
				minID: 1,
				maxID: 100,
				limit: 2,
			},
			want:    []*object.MediaAttachment{},
			wantErr: false,
		},
		{
			name: "error",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `media_attachment` WHERE `status_id` BETWEEN ? AND ? AND `delete_at` IS NULL ORDER BY `create_at` DESC LIMIT ?")).
					WithArgs(1, 100, 2).
					WillReturnError(errors.New("error"))
			},
			args: args{
				ctx:   context.Background(),
				minID: 1,
				maxID: 100,
				limit: 2,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.query(mock)
			got, err := r.Select(tt.args.ctx, tt.args.minID, tt.args.maxID, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("mediaAttachment.Select() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("mediaAttachment.Select() returned diff (want -> got):\n%s", diff)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}
