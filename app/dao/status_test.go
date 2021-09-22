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

func Test_status_FindByID(t *testing.T) {
	createAt, _ := time.Parse("2012-01-02", "2020-01-01")

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	r := &status{
		db: db,
	}

	type args struct {
		ctx context.Context
		id  object.StatusID
	}
	tests := []struct {
		name    string
		query   func(s sqlxmock.Sqlmock)
		args    args
		want    *object.Status
		wantErr bool
	}{
		{
			name: "ok",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `status` WHERE `id` = ? AND `delete_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"account_id",
								"content",
								"create_at",
								"delete_at",
							},
						).
							AddRow(1, 1, "content", createAt, nil),
					)
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &object.Status{
				ID:        1,
				AccountID: 1,
				Content:   "content",
				CreateAt:  object.DateTime{Time: createAt},
				DeleteAt:  nil,
			},
		},
		{
			name: "no rows",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `status` WHERE `id` = ? AND `delete_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"account_id",
								"content",
								"create_at",
								"delete_at",
							},
						),
					)
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "error",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `status` WHERE `id` = ? AND `delete_at` IS NULL")).
					WithArgs(1).
					WillReturnError(errors.New("error"))
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.query(mock)
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("status.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("status.FindByID() returned diff (want -> got):\n%s", diff)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func Test_status_FindByIDs(t *testing.T) {
	createAt, _ := time.Parse("2006-01-02", "2020-01-01")

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	r := &status{
		db: db,
	}

	type args struct {
		ctx context.Context
		ids []object.StatusID
	}
	tests := []struct {
		name    string
		query   func(s sqlxmock.Sqlmock)
		args    args
		want    []*object.Status
		wantErr bool
	}{
		{
			name: "ok",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `status` WHERE `id` IN (?, ?, ?) AND `delete_at` IS NULL")).
					WithArgs(1, 2, 3).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"account_id",
								"content",
								"create_at",
								"delete_at",
							},
						).
							AddRow(1, 1, "content1", createAt, nil).
							AddRow(2, 2, "content2", createAt, nil).
							AddRow(3, 3, "content3", createAt, nil),
					)
			},
			args: args{
				ctx: context.Background(),
				ids: []object.StatusID{1, 2, 3},
			},
			want: []*object.Status{
				{
					ID:        1,
					AccountID: 1,
					Content:   "content1",
					CreateAt:  object.DateTime{Time: createAt},
					DeleteAt:  nil,
				},
				{
					ID:        2,
					AccountID: 2,
					Content:   "content2",
					CreateAt:  object.DateTime{Time: createAt},
					DeleteAt:  nil,
				},
				{
					ID:        3,
					AccountID: 3,
					Content:   "content3",
					CreateAt:  object.DateTime{Time: createAt},
					DeleteAt:  nil,
				},
			},
			wantErr: false,
		},
		{
			name: "no rows",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `status` WHERE `id` IN (?, ?, ?) AND `delete_at` IS NULL")).
					WithArgs(1, 2, 3).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"account_id",
								"content",
								"create_at",
								"delete_at",
							},
						),
					)
			},
			args: args{
				ctx: context.Background(),
				ids: []object.StatusID{1, 2, 3},
			},
			want:    []*object.Status{},
			wantErr: false,
		},
		{
			name: "error",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `status` WHERE `id` IN (?, ?, ?) AND `delete_at` IS NULL")).
					WithArgs(1, 2, 3).
					WillReturnError(errors.New("error"))
			},
			args: args{
				ctx: context.Background(),
				ids: []object.StatusID{1, 2, 3},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.query(mock)
			got, err := r.FindByIDs(tt.args.ctx, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("status.FindByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("status.FindByIDs() returned diff (want -> got):\n%s", diff)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func Test_status_Select(t *testing.T) {
	createAt, _ := time.Parse("2012-01-02", "2020-01-01")

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	r := &status{
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
		want    []*object.Status
		wantErr bool
	}{
		{
			name: "ok",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `status` WHERE `id` BETWEEN ? AND ? AND `delete_at` IS NULL ORDER BY `create_at` DESC LIMIT ?")).
					WithArgs(1, 100, 2).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"account_id",
								"content",
								"create_at",
								"delete_at",
							},
						).
							AddRow(1, 1, "content1", createAt, nil).
							AddRow(2, 2, "content2", createAt, nil),
					)
			},
			args: args{
				ctx:   context.Background(),
				minID: 1,
				maxID: 100,
				limit: 2,
			},
			want: []*object.Status{
				{
					ID:        1,
					AccountID: 1,
					Content:   "content1",
					CreateAt:  object.DateTime{Time: createAt},
					DeleteAt:  nil,
				},
				{
					ID:        2,
					AccountID: 2,
					Content:   "content2",
					CreateAt:  object.DateTime{Time: createAt},
					DeleteAt:  nil,
				},
			},
		},
		{
			name: "no rows",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `status` WHERE `id` BETWEEN ? AND ? AND `delete_at` IS NULL ORDER BY `create_at` DESC LIMIT ?")).
					WithArgs(1, 100, 2).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"account_id",
								"content",
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
			want:    []*object.Status{},
			wantErr: false,
		},
		{
			name: "error",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `status` WHERE `id` BETWEEN ? AND ? AND `delete_at` IS NULL ORDER BY `create_at` DESC LIMIT ?")).
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
				t.Errorf("status.Select() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("status.Select() returned diff (want -> got):\n%s", diff)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func Test_status_Insert(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	r := &status{
		db: db,
	}

	type args struct {
		ctx       context.Context
		accountID object.AccountID
		content   string
	}
	tests := []struct {
		name    string
		query   func(s sqlxmock.Sqlmock)
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "ok",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectPrepare(regexp.QuoteMeta("INSERT INTO `status` (`account_id`, `content`) VALUES (?, ?)")).
					ExpectExec().
					WithArgs(1, "content").
					WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			args: args{
				ctx:       context.Background(),
				accountID: 1,
				content:   "content",
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "error",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectPrepare(regexp.QuoteMeta("INSERT INTO `status` (`account_id`, `content`) VALUES (?, ?)")).
					ExpectExec().
					WithArgs(1, "content").
					WillReturnError(errors.New("error"))
			},
			args: args{
				ctx:       context.Background(),
				accountID: 1,
				content:   "content",
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.query(mock)
			got, err := r.Insert(tt.args.ctx, tt.args.accountID, tt.args.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("status.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("status.Insert() returned diff (want -> got):\n%s", diff)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func Test_status_Delete(t *testing.T) {
	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	r := &status{
		db: db,
	}

	type args struct {
		ctx       context.Context
		statusID  object.StatusID
		accountID object.AccountID
	}
	tests := []struct {
		name    string
		query   func(s sqlxmock.Sqlmock)
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectPrepare(regexp.QuoteMeta("UPDATE `status` SET `delete_at` = NOW() WHERE `id` = ?")).
					ExpectExec().
					WithArgs(1, 1).
					WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			args: args{
				ctx:       context.Background(),
				statusID:  1,
				accountID: 1,
			},
			wantErr: false,
		},
		{
			name: "error",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectPrepare(regexp.QuoteMeta("UPDATE `status` SET `delete_at` = NOW() WHERE `id` = ?")).
					ExpectExec().
					WithArgs(1, 1).
					WillReturnError(errors.New("error"))
			},
			args: args{
				ctx:       context.Background(),
				statusID:  1,
				accountID: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.query(mock)
			if err := r.Delete(tt.args.ctx, tt.args.statusID, tt.args.accountID); (err != nil) != tt.wantErr {
				t.Errorf("status.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}
