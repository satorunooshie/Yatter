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

func Test_account_FindByID(t *testing.T) {
	displayName := "ニックネーム"
	avatar := "http://example.com/avatar"
	header := "http://example.com/header"
	note := "一言"
	createAt, _ := time.Parse("2006-01-02", "2020-01-01")

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	r := &account{
		db: db,
	}

	type args struct {
		ctx context.Context
		id  object.AccountID
	}
	tests := []struct {
		name    string
		query   func(s sqlxmock.Sqlmock)
		args    args
		want    *object.Account
		wantErr bool
	}{
		{
			name: "ok",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `account` WHERE `id` = ? AND `delete_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"username",
								"password_hash",
								"display_name",
								"avatar",
								"header",
								"note",
								"create_at",
								"delete_at",
							},
						).
							AddRow(1, "名前", "hash", "ニックネーム", "http://example.com/avatar", "http://example.com/header", "一言", createAt, nil),
					)
			},
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			want: &object.Account{
				ID:           1,
				Username:     "名前",
				PasswordHash: "hash",
				DisplayName:  &displayName,
				Avatar:       &avatar,
				Header:       &header,
				Note:         &note,
				CreateAt:     object.DateTime{Time: createAt},
				DeleteAt:     nil,
			},
			wantErr: false,
		},
		{
			name: "no rows",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `account` WHERE `id` = ? AND `delete_at` IS NULL")).
					WithArgs(1).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"username",
								"password_hash",
								"display_name",
								"avatar",
								"header",
								"note",
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
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `account` WHERE `id` = ? AND `delete_at` IS NULL")).
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
				t.Errorf("account.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("account.FindByID() returned diff (want -> got):\n%s", diff)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func Test_account_FindByIDs(t *testing.T) {
	displayName := "ニックネーム"
	avatar := "http://example.com/avatar"
	header := "http://example.com/header"
	note := "一言"
	createAt, _ := time.Parse("2006-01-02", "2020-01-01")

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	r := &account{
		db: db,
	}

	type args struct {
		ctx context.Context
		ids []object.AccountID
	}
	tests := []struct {
		name    string
		query   func(s sqlxmock.Sqlmock)
		args    args
		want    []*object.Account
		wantErr bool
	}{
		{
			name: "ok",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `account` WHERE `id` IN (?, ?, ?) AND `delete_at` IS NULL")).
					WithArgs(1, 2, 3).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"username",
								"password_hash",
								"display_name",
								"avatar",
								"header",
								"note",
								"create_at",
								"delete_at",
							},
						).
							AddRow(1, "名前1", "hash1", "ニックネーム", "http://example.com/avatar", "http://example.com/header", "一言", createAt, nil).
							AddRow(2, "名前2", "hash2", "ニックネーム", "http://example.com/avatar", "http://example.com/header", "一言", createAt, nil).
							AddRow(3, "名前3", "hash3", "ニックネーム", "http://example.com/avatar", "http://example.com/header", "一言", createAt, nil),
					)
			},
			args: args{
				ctx: context.Background(),
				ids: []object.AccountID{1, 2, 3},
			},
			want: []*object.Account{
				{
					ID:           1,
					Username:     "名前1",
					PasswordHash: "hash1",
					DisplayName:  &displayName,
					Avatar:       &avatar,
					Header:       &header,
					Note:         &note,
					CreateAt:     object.DateTime{Time: createAt},
					DeleteAt:     nil,
				},
				{
					ID:           2,
					Username:     "名前2",
					PasswordHash: "hash2",
					DisplayName:  &displayName,
					Avatar:       &avatar,
					Header:       &header,
					Note:         &note,
					CreateAt:     object.DateTime{Time: createAt},
					DeleteAt:     nil,
				},
				{
					ID:           3,
					Username:     "名前3",
					PasswordHash: "hash3",
					DisplayName:  &displayName,
					Avatar:       &avatar,
					Header:       &header,
					Note:         &note,
					CreateAt:     object.DateTime{Time: createAt},
					DeleteAt:     nil,
				},
			},
			wantErr: false,
		},
		{
			name: "no rows",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `account` WHERE `id` IN (?, ?, ?) AND `delete_at` IS NULL")).
					WithArgs(1, 2, 3).
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"username",
								"password_hash",
								"display_name",
								"avatar",
								"header",
								"note",
								"create_at",
								"delete_at",
							},
						),
					)
			},
			args: args{
				ctx: context.Background(),
				ids: []object.AccountID{1, 2, 3},
			},
			want:    []*object.Account{},
			wantErr: false,
		},
		{
			name: "error",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `account` WHERE `id` IN (?, ?, ?) AND `delete_at` IS NULL")).
					WithArgs(1, 2, 3).
					WillReturnError(errors.New("error"))
			},
			args: args{
				ctx: context.Background(),
				ids: []object.AccountID{1, 2, 3},
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
				t.Errorf("account.FindByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("account.FindByIDs() returned diff (want -> got):\n%s", diff)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func Test_account_FindByUsername(t *testing.T) {
	displayName := "ニックネーム"
	avatar := "http://example.com/avatar"
	header := "http://example.com/header"
	note := "一言"
	createAt, _ := time.Parse("2006-01-02", "2020-01-01")

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	r := &account{
		db: db,
	}

	type args struct {
		ctx      context.Context
		username string
	}
	tests := []struct {
		name    string
		query   func(s sqlxmock.Sqlmock)
		args    args
		want    *object.Account
		wantErr bool
	}{
		{
			name: "ok",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `account` WHERE `username` = ? AND `delete_at` IS NULL")).
					WithArgs("名前").
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"username",
								"password_hash",
								"display_name",
								"avatar",
								"header",
								"note",
								"create_at",
								"delete_at",
							},
						).
							AddRow(1, "名前", "hash", "ニックネーム", "http://example.com/avatar", "http://example.com/header", "一言", createAt, nil),
					)
			},
			args: args{
				ctx:      context.Background(),
				username: "名前",
			},
			want: &object.Account{
				ID:           1,
				Username:     "名前",
				PasswordHash: "hash",
				DisplayName:  &displayName,
				Avatar:       &avatar,
				Header:       &header,
				Note:         &note,
				CreateAt:     object.DateTime{Time: createAt},
				DeleteAt:     nil,
			},
			wantErr: false,
		},
		{
			name: "no rows",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `account` WHERE `username` = ? AND `delete_at` IS NULL")).
					WithArgs("名前").
					WillReturnRows(
						sqlxmock.NewRows(
							[]string{
								"id",
								"username",
								"password_hash",
								"display_name",
								"avatar",
								"header",
								"note",
								"create_at",
								"delete_at",
							},
						),
					)
			},
			args: args{
				ctx:      context.Background(),
				username: "名前",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "error",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `account` WHERE `username` = ? AND `delete_at` IS NULL")).
					WithArgs("名前").
					WillReturnError(errors.New("error"))
			},
			args: args{
				ctx:      context.Background(),
				username: "名前",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.query(mock)
			got, err := r.FindByUsername(tt.args.ctx, tt.args.username)
			if (err != nil) != tt.wantErr {
				t.Errorf("account.FindByUsername() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("account.FindByUsername() returned diff (want -> got):\n%s", diff)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}

func Test_account_Insert(t *testing.T) {
	createAt, _ := time.Parse("2006-01-02", "2020-01-01")

	db, mock, err := sqlxmock.Newx()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = db.Close()
	}()

	r := &account{
		db: db,
	}

	type args struct {
		ctx          context.Context
		username     string
		passwordHash string
		createAt     time.Time
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
				s.ExpectPrepare(regexp.QuoteMeta("INSERT INTO `account` (`username`, `password_hash`, `create_at`) VALUES (?, ?, ?)")).
					ExpectExec().
					WithArgs("名前", "hash", createAt).
					WillReturnResult(sqlxmock.NewResult(1, 1))
			},
			args: args{
				ctx:          context.Background(),
				username:     "名前",
				passwordHash: "hash",
				createAt:     createAt,
			},
			wantErr: false,
		},
		{
			name: "error",
			query: func(s sqlxmock.Sqlmock) {
				s.ExpectPrepare(regexp.QuoteMeta("INSERT INTO `account` (`username`, `password_hash`, `create_at`) VALUES (?, ?, ?)")).
					ExpectExec().
					WithArgs("名前", "hash", createAt).
					WillReturnError(errors.New("error"))
			},
			args: args{
				ctx:          context.Background(),
				username:     "名前",
				passwordHash: "hash",
				createAt:     createAt,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.query(mock)
			if err := r.Insert(tt.args.ctx, tt.args.username, tt.args.passwordHash, tt.args.createAt); (err != nil) != tt.wantErr {
				t.Errorf("account.Insert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %v", err)
			}
		})
	}
}
