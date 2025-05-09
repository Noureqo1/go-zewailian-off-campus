package ws

import (
	"context"
	"server/internal/message"
)

type MessageServiceAdapter struct {
	messageService message.Service
}

func NewMessageServiceAdapter(messageService message.Service) MessageService {
	return &MessageServiceAdapter{
		messageService: messageService,
	}
}

func (a *MessageServiceAdapter) SaveMessage(ctx context.Context, msg interface{}) error {
	wsMsg, ok := msg.(*Message)
	if !ok {
		return nil 
	}

	dbMsg := &message.Message{
		ID:        wsMsg.ID,
		RoomID:    wsMsg.RoomID,
		UserID:    "", 
		Username:  wsMsg.Username,
		Content:   wsMsg.Content,
		Type:      string(wsMsg.Type),
		Timestamp: wsMsg.Timestamp,
		Recipient: wsMsg.Recipient,
	}
	
	return a.messageService.SaveMessage(ctx, dbMsg)
}

func (a *MessageServiceAdapter) CreateRoom(ctx context.Context, id, name, ownerID string) (interface{}, error) {
	return a.messageService.CreateRoom(ctx, id, name, ownerID)
}

func (a *MessageServiceAdapter) UpdateRoomActivity(ctx context.Context, roomID string) error {
	return a.messageService.UpdateRoomActivity(ctx, roomID)
}
