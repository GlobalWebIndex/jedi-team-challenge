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
	"time"
)

func TestChatRepository_CreateMessage(t *testing.T) {
	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		message *domain.Message
	}
	tests := []struct {
		name                          string
		args                          args
		mockSqlMessageQueryExpected   string
		mockInsertedMessageIdReturned uuid.UUID
		expectedMessageUid            uuid.UUID
	}{
		{
			name: "valid",
			args: args{
				message: &domain.Message{
					ChatSessionID: uuid.UUID{0x12, 0x34, 0x56, 0x78},
					Sender:        "USER",
					Content:       "ablaabla",
					CreatedAt:     time.Time{},
					Feedback:      nil,
				},
			},
			mockSqlMessageQueryExpected:   `INSERT INTO "messages" ("chat_session_id","sender","content","created_at","feedback") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`,
			mockInsertedMessageIdReturned: uuid.UUID{0x42, 0x34, 0x56, 0x78},
			expectedMessageUid:            uuid.UUID{0x42, 0x34, 0x56, 0x78},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MessageRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()

			mockDb.ExpectQuery(regexp.QuoteMeta(tt.mockSqlMessageQueryExpected)).
				WithArgs(
					tt.args.message.ChatSessionID, tt.args.message.Sender, tt.args.message.Content, sqlmock.AnyArg(), tt.args.message.Feedback,
				).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.mockInsertedMessageIdReturned))
			mockDb.ExpectCommit()

			actual, err := repo.CreateMessage(context.Background(), tt.args.message)
			if err != nil {
				t.Errorf("CreateMessage() error = %v", err)
			}

			assert.Equal(t, tt.expectedMessageUid, actual)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
