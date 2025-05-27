package repositories

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/loukaspe/jedi-team-challenge/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

func TestChatRepository_GetMessage(t *testing.T) {
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
		name                        string
		args                        args
		mockSqlMessageQueryExpected string
		mockMessageReturned         *Message
		expected                    *domain.Message
	}{
		{
			name: "valid",
			args: args{
				uuid: uuid.UUID{0x12, 0x34, 0x56, 0x78},
			},
			mockSqlMessageQueryExpected: `SELECT * FROM "messages" WHERE id = $1 LIMIT $2`,
			mockMessageReturned: &Message{
				ID:            uuid.UUID{0x12, 0x34, 0x56, 0x78},
				ChatSessionID: uuid.UUID{0x22, 0x34, 0x56, 0x78},
				Sender:        "user1",
				Content:       "Hello, world!",
			},
			expected: &domain.Message{
				ID:            uuid.UUID{0x12, 0x34, 0x56, 0x78},
				ChatSessionID: uuid.UUID{0x22, 0x34, 0x56, 0x78},
				Sender:        "user1",
				Content:       "Hello, world!",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MessageRepository{
				db: gormDb,
			}

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlMessageQueryExpected)).
				WithArgs(tt.args.uuid, 1).
				WillReturnRows(
					sqlmock.NewRows(
						[]string{"id", "chat_session_id", "sender", "created_at", "content"},
					).AddRow(
						tt.mockMessageReturned.ID, tt.mockMessageReturned.ChatSessionID,
						tt.mockMessageReturned.Sender, tt.mockMessageReturned.CreatedAt, tt.mockMessageReturned.Content,
					))

			actual, err := repo.GetMessage(context.Background(), tt.args.uuid)
			if err != nil {
				t.Errorf("GetMessage() error = %v", err)
				return
			}

			assert.Equal(t, tt.expected, actual)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
