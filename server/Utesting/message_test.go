package testing

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"server/internal/message"
)

func TestMessageRepository(t *testing.T) {

	repo := NewMockMessageService()

	t.Run("create and get message", func(t *testing.T) {
		msg := &message.Message{
			ID:        uuid.New().String(),
			RoomID:    uuid.New().String(),
			UserID:    uuid.New().String(),
			Username:  "testuser",
			Content:   "Test message",
			Type:      "text",
			Timestamp: time.Now(),
		}

		// Create message
		err := repo.SaveMessage(context.Background(), msg)
		assert.NoError(t, err)

		// Get message
		messages, err := repo.GetMessagesByRoom(context.Background(), msg.RoomID, 10, 0)
		assert.NoError(t, err)
		assert.NotEmpty(t, messages)
		assert.Equal(t, msg.Content, messages[0].Content)
	})

	t.Run("get messages by room with pagination", func(t *testing.T) {
		roomID := uuid.New().String()
		// Create multiple messages
		msgs := []*message.Message{
			{
				ID:        uuid.New().String(),
				RoomID:    roomID,
				UserID:    uuid.New().String(),
				Username:  "user1",
				Content:   "Message 1",
				Type:      "text",
				Timestamp: time.Now(),
			},
			{
				ID:        uuid.New().String(),
				RoomID:    roomID,
				UserID:    uuid.New().String(),
				Username:  "user2",
				Content:   "Message 2",
				Type:      "text",
				Timestamp: time.Now().Add(time.Second),
			},
		}

		for _, msg := range msgs {
			err := repo.SaveMessage(context.Background(), msg)
			assert.NoError(t, err)
		}

		// Get messages for room
		messages, err := repo.GetMessagesByRoom(context.Background(), roomID, 1, 0)
		assert.NoError(t, err)
		assert.Len(t, messages, 1)
	})
}
