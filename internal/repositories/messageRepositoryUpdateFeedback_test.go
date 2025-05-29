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

func TestChatRepository_UpdateMessageFeedback(t *testing.T) {
	feedback := "ablabla"

	db, mockDb, err := sqlmock.New()
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))

	type args struct {
		uuid     uuid.UUID
		feedback string
	}
	tests := []struct {
		name                        string
		args                        args
		mockSqlMessageQueryExpected string
		expected                    *domain.Message
	}{
		{
			name: "valid",
			args: args{
				uuid:     uuid.UUID{0x12, 0x34, 0x56, 0x78},
				feedback: feedback,
			},
			mockSqlMessageQueryExpected: `UPDATE "messages" SET "feedback"=$1 WHERE id = $2`,
			expected: &domain.Message{
				ID:            uuid.UUID{0x12, 0x34, 0x56, 0x78},
				ChatSessionID: uuid.UUID{0x22, 0x34, 0x56, 0x78},
				Sender:        "user1",
				Content:       "Hello, world!",
				Feedback:      &feedback,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &MessageRepository{
				db: gormDb,
			}

			mockDb.ExpectBegin()
			mockDb.ExpectExec(regexp.QuoteMeta(tt.mockSqlMessageQueryExpected)).
				WithArgs(tt.args.feedback, tt.args.uuid).
				WillReturnResult(sqlmock.NewResult(0, 1))
			mockDb.ExpectCommit()

			err := repo.UpdateMessageFeedback(context.Background(), tt.args.uuid, tt.args.feedback)
			if err != nil {
				t.Errorf("UpdateMessageFeedback() error = %v", err)
				return
			}

			assert.Nil(t, err)

			if err = mockDb.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expections: %s", err)
			}
		})
	}
}
