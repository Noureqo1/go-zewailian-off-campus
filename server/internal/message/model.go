package message

import (
	"time"
)

// Message represents a chat message stored in the database
type Message struct {
	ID        string    `json:"id" db:"id"`
	RoomID    string    `json:"roomId" db:"room_id"`
	UserID    string    `json:"userId" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Content   string    `json:"content" db:"content"`
	Type      string    `json:"type" db:"type"`
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
	Recipient string    `json:"recipient,omitempty" db:"recipient"`
}

// Room represents a chat room stored in the database
type Room struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	OwnerID      string    `json:"ownerId" db:"owner_id"`
	Created      time.Time `json:"created" db:"created"`
	LastActivity time.Time `json:"lastActivity" db:"last_activity"`
}
