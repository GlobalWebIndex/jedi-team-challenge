package repositories

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	apierrors "github.com/loukaspe/jedi-team-challenge/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func TestChatRepository_UpdateChatSessionTitle(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		uuid  uuid.UUID
		title string
	}
	tests := []struct {
		name                      string
		args                      args
		mockSqlChatQueryExpected  string
		mockUpdatedChatIdReturned int
	}{
		{
			name: "valid",
			args: args{
				uuid:  uuid.UUID{0x12, 0x34, 0x56, 0x78},
				title: "mockTitle",
			},
			mockSqlChatQueryExpected:  `UPDATE "chat_sessions" SET "title"=$1,"updated_at"=$2 WHERE id = $3`,
			mockUpdatedChatIdReturned: 666,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ChatSessionRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlChatQueryExpected)).
				WithArgs(tt.args.title, sqlmock.AnyArg(), tt.args.uuid).
				WillReturnResult(sqlmock.NewResult(int64(tt.mockUpdatedChatIdReturned), 1))
			mockDb.ExpectCommit()

			err := repo.UpdateChatSessionTitle(
				context.Background(),
				tt.args.uuid,
				tt.args.title,
			)

			assert.NoError(t, err)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestChatRepository_UpdateChatSessionTitleWithError(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		uuid  uuid.UUID
		title string
	}
	tests := []struct {
		name                     string
		args                     args
		mockSqlChatQueryExpected string
		mockSqlErrorReturned     error
		expectedError            error
	}{
		{
			name: "random error",
			args: args{
				uuid:  uuid.UUID{0x12, 0x34, 0x56, 0x78},
				title: "mockTitle",
			},
			mockSqlChatQueryExpected: `UPDATE "chat_sessions" SET "title"=$1,"updated_at"=$2 WHERE id = $3`,
			mockSqlErrorReturned:     errors.New("random error"),
			expectedError:            errors.New("random error"),
		},
		{
			name: "title not found",
			args: args{
				uuid:  uuid.UUID{0x12, 0x34, 0x56, 0x78},
				title: "mockTitle",
			},
			mockSqlChatQueryExpected: `UPDATE "chat_sessions" SET "title"=$1,"updated_at"=$2 WHERE id = $3`,
			mockSqlErrorReturned:     gorm.ErrRecordNotFound,
			expectedError: apierrors.ResourceNotFoundErrorWrapper{
				OriginalError: errors.New("chatSessionID 12345678-0000-0000-0000-000000000000 not found"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ChatSessionRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlChatQueryExpected)).
				WithArgs(tt.args.title, sqlmock.AnyArg(), tt.args.uuid).
				WillReturnError(tt.mockSqlErrorReturned)
			mockDb.ExpectRollback()

			actual := repo.UpdateChatSessionTitle(
				context.Background(),
				tt.args.uuid,
				tt.args.title,
			)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
