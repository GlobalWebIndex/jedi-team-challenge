package repositories

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func TestChatRepository_CreateChatSession(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		chat *domain.ChatSession
	}
	tests := []struct {
		name                       string
		args                       args
		mockSqlChatQueryExpected   string
		mockInsertedChatIdReturned uuid.UUID
		expectedChatUid            uuid.UUID
	}{
		{
			name: "valid",
			args: args{
				chat: &domain.ChatSession{
					UserID: uuid.UUID{
						0x12, 0x34, 0x56, 0x78,
					},
					Title:    "mockTitle",
					Messages: nil,
				},
			},
			mockSqlChatQueryExpected:   `INSERT INTO "chat_sessions" ("user_id","title","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`,
			mockInsertedChatIdReturned: uuid.UUID{0x12, 0x34, 0x56, 0x78},
			expectedChatUid:            uuid.UUID{0x12, 0x34, 0x56, 0x78},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ChatSessionRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlChatQueryExpected)).
				WithArgs(
					tt.args.chat.UserID, tt.args.chat.Title, sqlmock.AnyArg(), sqlmock.AnyArg(),
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.mockInsertedChatIdReturned))
			mockDb.ExpectCommit()

			actual, err := repo.CreateChatSession(context.Background(), tt.args.chat)
			if err != nil {
				t.Errorf("CreateChat() error = %v", err)
			}

			assert.Equal(t, tt.expectedChatUid, actual)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestChatRepository_CreateChatSessionWithError(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		chat *domain.ChatSession
	}
	tests := []struct {
		name                     string
		args                     args
		mockSqlChatQueryExpected string
		expectedErrorMessage     string
	}{
		{
			name: "random error",
			args: args{
				chat: &domain.ChatSession{
					Title:  "mockTitle",
					UserID: uuid.UUID{0x12, 0x34, 0x56, 0x78},
				},
			},
			mockSqlChatQueryExpected: `INSERT INTO "chat_sessions" ("user_id","title","created_at","updated_at") VALUES ($1,$2,$3,$4) RETURNING "id"`,
			expectedErrorMessage:     "random error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ChatSessionRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlChatQueryExpected)).
				WithArgs(
					tt.args.chat.UserID, tt.args.chat.Title, sqlmock.AnyArg(), sqlmock.AnyArg(),
				).
				WillReturnError(errors.New(tt.expectedErrorMessage))
			mockDb.ExpectRollback()

			_, err := repo.CreateChatSession(context.Background(), tt.args.chat)
			actualErrorMessage := err.Error()

			assert.Equal(t, tt.expectedErrorMessage, actualErrorMessage)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
