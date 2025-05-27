package repositories

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
	apierrors "github.com/loukaspe/jedi-team-challenge/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

func TestChatRepository_GetChatSession(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		uuid uuid.UUID
	}
	tests := []struct {
		name                         string
		args                         args
		mockSqlChatQueryExpected     string
		mockSqlMessagesQueryExpected string
		mockChatReturned             *ChatSession
		mockMessagesReturned         []*Message
		expected                     *domain.ChatSession
	}{
		{
			name: "valid",
			args: args{
				uuid: uuid.UUID{0x12, 0x34, 0x56, 0x78},
			},
			mockSqlChatQueryExpected:     `SELECT * FROM "chat_sessions" WHERE id = $1 LIMIT $2`,
			mockSqlMessagesQueryExpected: `SELECT * FROM "messages" WHERE "messages"."chat_session_id" = $1`,
			mockChatReturned: &ChatSession{
				ID:        uuid.UUID{0x12, 0x34, 0x56, 0x78},
				UserID:    uuid.UUID{0x22, 0x34, 0x56, 0x88},
				Title:     "mockTitle",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			mockMessagesReturned: []*Message{
				{
					ID:            uuid.UUID{0x02, 0x34, 0x56, 0x68},
					ChatSessionID: uuid.UUID{0x12, 0x34, 0x56, 0x78},
					Sender:        "SYSTEM",
					Content:       "MAY THE FORCE BE WITH YOU",
					CreatedAt:     time.Time{},
				},
				{
					ID:            uuid.UUID{0x052, 0x34, 0x56, 0x58},
					ChatSessionID: uuid.UUID{0x12, 0x34, 0x56, 0x78},
					Sender:        "USER",
					Content:       "2MAY THE FORCE BE WITH YOU2",
					CreatedAt:     time.Time{},
				},
			},
			expected: &domain.ChatSession{
				ID:        uuid.UUID{0x12, 0x34, 0x56, 0x78},
				UserID:    uuid.UUID{0x22, 0x34, 0x56, 0x88},
				Title:     "mockTitle",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				Messages: []*domain.Message{
					{
						ID:            uuid.UUID{0x02, 0x34, 0x56, 0x68},
						ChatSessionID: uuid.UUID{0x12, 0x34, 0x56, 0x78},
						Sender:        "SYSTEM",
						Content:       "MAY THE FORCE BE WITH YOU",
						CreatedAt:     time.Time{},
					},
					{
						ID:            uuid.UUID{0x052, 0x34, 0x56, 0x58},
						ChatSessionID: uuid.UUID{0x12, 0x34, 0x56, 0x78},
						Sender:        "USER",
						Content:       "2MAY THE FORCE BE WITH YOU2",
						CreatedAt:     time.Time{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ChatSessionRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlChatQueryExpected)).
				WithArgs(tt.args.uuid, 1).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "title", "user_id", "created_at", "updated_at"},
					).AddRow(
						tt.mockChatReturned.ID, tt.mockChatReturned.Title, tt.mockChatReturned.UserID, tt.mockChatReturned.CreatedAt, tt.mockChatReturned.UpdatedAt,
					),
				)

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlMessagesQueryExpected)).
				WithArgs(tt.args.uuid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "chat_session_id", "sender", "created_at", "content"},
					).AddRow(
						tt.mockMessagesReturned[0].ID, tt.mockMessagesReturned[0].ChatSessionID,
						tt.mockMessagesReturned[0].Sender, tt.mockMessagesReturned[0].CreatedAt, tt.mockMessagesReturned[0].Content,
					).AddRow(
						tt.mockMessagesReturned[1].ID, tt.mockMessagesReturned[1].ChatSessionID,
						tt.mockMessagesReturned[1].Sender, tt.mockMessagesReturned[1].CreatedAt, tt.mockMessagesReturned[1].Content,
					),
				)

			actual, err := repo.GetChatSession(context.Background(), tt.args.uuid)
			if err != nil {
				t.Errorf("GetChatSession() error = %v", err)
				return
			}

			assert.Equal(t, tt.expected, actual)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestChatRepository_GetChatSessionWithError(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		uuid uuid.UUID
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
				uuid: uuid.UUID{0x12, 0x34, 0x56, 0x78},
			},
			mockSqlChatQueryExpected: `SELECT * FROM "chat_sessions" WHERE id = $1 LIMIT $2`,
			mockSqlErrorReturned:     errors.New("random error"),
			expectedError:            errors.New("random error"),
		},
		{
			name: "session not found",
			args: args{
				uuid: uuid.UUID{0x12, 0x34, 0x56, 0x78},
			},
			mockSqlChatQueryExpected: `SELECT * FROM "chat_sessions" WHERE id = $1 LIMIT $2`,
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

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlChatQueryExpected)).
				WithArgs(tt.args.uuid, 1).
				WillReturnError(tt.mockSqlErrorReturned)

			_, actual := repo.GetChatSession(context.Background(), tt.args.uuid)

			assert.Equal(t, actual, tt.expectedError)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}

func TestChatRepository_GetUserChatSessions(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		uuid uuid.UUID
	}
	tests := []struct {
		name                         string
		args                         args
		mockSqlChatQueryExpected     string
		mockSqlMessagesQueryExpected string
		mockChatsReturned            []*ChatSession
		mockMessagesReturned         []*Message
		expected                     []*domain.ChatSession
	}{
		{
			name: "valid",
			args: args{
				uuid: uuid.UUID{0x22, 0x34, 0x56, 0x88},
			},
			mockSqlChatQueryExpected:     `SELECT * FROM "chat_sessions" WHERE user_id = $1`,
			mockSqlMessagesQueryExpected: `SELECT * FROM "messages" WHERE "messages"."chat_session_id" IN ($1,$2)`,
			mockChatsReturned: []*ChatSession{
				&ChatSession{
					ID:        uuid.UUID{0x32, 0x34, 0x56, 0x78},
					UserID:    uuid.UUID{0x22, 0x34, 0x56, 0x88},
					Title:     "mockTitle",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
				&ChatSession{
					ID:        uuid.UUID{0x42, 0x34, 0x56, 0x78},
					UserID:    uuid.UUID{0x22, 0x34, 0x56, 0x88},
					Title:     "mockTitle",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
			},
			mockMessagesReturned: []*Message{

				{
					ID:            uuid.UUID{0x02, 0x34, 0x56, 0x68},
					ChatSessionID: uuid.UUID{0x32, 0x34, 0x56, 0x78},
					Sender:        "SYSTEM",
					Content:       "MAY THE FORCE BE WITH YOU",
					CreatedAt:     time.Time{},
				},
				{
					ID:            uuid.UUID{0x052, 0x34, 0x56, 0x58},
					ChatSessionID: uuid.UUID{0x32, 0x34, 0x56, 0x78},
					Sender:        "USER",
					Content:       "2MAY THE FORCE BE WITH YOU2",
					CreatedAt:     time.Time{},
				},

				{
					ID:            uuid.UUID{0x02, 0x34, 0x56, 0x68},
					ChatSessionID: uuid.UUID{0x42, 0x34, 0x56, 0x78},
					Sender:        "SYSTEM",
					Content:       "3MAY THE FORCE BE WITH YOU3",
					CreatedAt:     time.Time{},
				},
				{
					ID:            uuid.UUID{0x052, 0x34, 0x56, 0x58},
					ChatSessionID: uuid.UUID{0x42, 0x34, 0x56, 0x78},
					Sender:        "USER",
					Content:       "4MAY THE FORCE BE WITH YOU4",
					CreatedAt:     time.Time{},
				},
			},
			expected: []*domain.ChatSession{
				&domain.ChatSession{
					ID:        uuid.UUID{0x32, 0x34, 0x56, 0x78},
					UserID:    uuid.UUID{0x22, 0x34, 0x56, 0x88},
					Title:     "mockTitle",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
					Messages: []*domain.Message{
						{
							ID:            uuid.UUID{0x02, 0x34, 0x56, 0x68},
							ChatSessionID: uuid.UUID{0x32, 0x34, 0x56, 0x78},
							Sender:        "SYSTEM",
							Content:       "MAY THE FORCE BE WITH YOU",
							CreatedAt:     time.Time{},
						},
						{

							ID:            uuid.UUID{0x052, 0x34, 0x56, 0x58},
							ChatSessionID: uuid.UUID{0x32, 0x34, 0x56, 0x78},
							Sender:        "USER",
							Content:       "2MAY THE FORCE BE WITH YOU2",
							CreatedAt:     time.Time{},
						},
					},
				},
				&domain.ChatSession{
					ID:        uuid.UUID{0x42, 0x34, 0x56, 0x78},
					UserID:    uuid.UUID{0x22, 0x34, 0x56, 0x88},
					Title:     "mockTitle",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
					Messages: []*domain.Message{
						{
							ID:            uuid.UUID{0x02, 0x34, 0x56, 0x68},
							ChatSessionID: uuid.UUID{0x42, 0x34, 0x56, 0x78},
							Sender:        "SYSTEM",
							Content:       "3MAY THE FORCE BE WITH YOU3",
							CreatedAt:     time.Time{},
						},
						{
							ID:            uuid.UUID{0x052, 0x34, 0x56, 0x58},
							ChatSessionID: uuid.UUID{0x42, 0x34, 0x56, 0x78},
							Sender:        "USER",
							Content:       "4MAY THE FORCE BE WITH YOU4",
							CreatedAt:     time.Time{},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &ChatSessionRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlChatQueryExpected)).
				WithArgs(tt.args.uuid).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "title", "user_id", "created_at", "updated_at"},
					).AddRow(
						tt.mockChatsReturned[0].ID, tt.mockChatsReturned[0].Title, tt.mockChatsReturned[0].UserID, tt.mockChatsReturned[0].CreatedAt, tt.mockChatsReturned[0].UpdatedAt,
					).AddRow(
						tt.mockChatsReturned[1].ID, tt.mockChatsReturned[1].Title, tt.mockChatsReturned[1].UserID, tt.mockChatsReturned[1].CreatedAt, tt.mockChatsReturned[1].UpdatedAt,
					),
				)

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlMessagesQueryExpected)).
				WithArgs(tt.mockChatsReturned[0].ID, tt.mockChatsReturned[1].ID).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "chat_session_id", "sender", "created_at", "content"},
					).AddRow(
						tt.mockMessagesReturned[0].ID, tt.mockMessagesReturned[0].ChatSessionID,
						tt.mockMessagesReturned[0].Sender, tt.mockMessagesReturned[0].CreatedAt, tt.mockMessagesReturned[0].Content,
					).AddRow(
						tt.mockMessagesReturned[1].ID, tt.mockMessagesReturned[1].ChatSessionID,
						tt.mockMessagesReturned[1].Sender, tt.mockMessagesReturned[1].CreatedAt, tt.mockMessagesReturned[1].Content,
					).AddRow(
						tt.mockMessagesReturned[2].ID, tt.mockMessagesReturned[2].ChatSessionID,
						tt.mockMessagesReturned[2].Sender, tt.mockMessagesReturned[2].CreatedAt, tt.mockMessagesReturned[2].Content,
					).AddRow(
						tt.mockMessagesReturned[3].ID, tt.mockMessagesReturned[3].ChatSessionID,
						tt.mockMessagesReturned[3].Sender, tt.mockMessagesReturned[3].CreatedAt, tt.mockMessagesReturned[3].Content,
					),
				)

			actual, err := repo.GetUserChatSessions(context.Background(), tt.args.uuid)
			if err != nil {
				t.Errorf("GetChatSession() error = %v", err)
				return
			}

			assert.Equal(t, tt.expected, actual)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
