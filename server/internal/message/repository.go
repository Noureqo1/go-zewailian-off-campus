package message

import (
	"context"
	"database/sql"
	"time"
)

// Repository defines the interface for message data access
type Repository interface {
	// Message operations
	SaveMessage(ctx context.Context, message *Message) error
	GetMessagesByRoom(ctx context.Context, roomID string, limit, offset int) ([]*Message, error)
	GetMessageByID(ctx context.Context, id string) (*Message, error)
	
	// Room operations
	CreateRoom(ctx context.Context, room *Room) error
	GetRooms(ctx context.Context) ([]*Room, error)
	GetRoomByID(ctx context.Context, id string) (*Room, error)
	UpdateRoomActivity(ctx context.Context, roomID string) error
}

// PostgresRepository implements the Repository interface using PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{
		db: db,
	}
}

// SaveMessage stores a message in the database
func (r *PostgresRepository) SaveMessage(ctx context.Context, message *Message) error {
	query := `
		INSERT INTO messages (id, room_id, user_id, username, content, type, timestamp, recipient)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		message.ID,
		message.RoomID,
		message.UserID,
		message.Username,
		message.Content,
		message.Type,
		message.Timestamp,
		message.Recipient,
	)
	
	return err
}

// GetMessagesByRoom retrieves messages for a specific room with pagination
func (r *PostgresRepository) GetMessagesByRoom(ctx context.Context, roomID string, limit, offset int) ([]*Message, error) {
	query := `
		SELECT id, room_id, user_id, username, content, type, timestamp, recipient
		FROM messages
		WHERE room_id = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		err := rows.Scan(
			&msg.ID,
			&msg.RoomID,
			&msg.UserID,
			&msg.Username,
			&msg.Content,
			&msg.Type,
			&msg.Timestamp,
			&msg.Recipient,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	
	return messages, nil
}

// GetMessageByID retrieves a message by its ID
func (r *PostgresRepository) GetMessageByID(ctx context.Context, id string) (*Message, error) {
	query := `
		SELECT id, room_id, user_id, username, content, type, timestamp, recipient
		FROM messages
		WHERE id = $1
	`
	
	msg := &Message{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&msg.ID,
		&msg.RoomID,
		&msg.UserID,
		&msg.Username,
		&msg.Content,
		&msg.Type,
		&msg.Timestamp,
		&msg.Recipient,
	)
	
	if err != nil {
		return nil, err
	}
	
	return msg, nil
}

// CreateRoom creates a new chat room
func (r *PostgresRepository) CreateRoom(ctx context.Context, room *Room) error {
	query := `
		INSERT INTO rooms (id, name, owner_id, created, last_activity)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		room.ID,
		room.Name,
		room.OwnerID,
		room.Created,
		room.LastActivity,
	)
	
	return err
}

// GetRooms retrieves all available chat rooms
func (r *PostgresRepository) GetRooms(ctx context.Context) ([]*Room, error) {
	query := `
		SELECT id, name, owner_id, created, last_activity
		FROM rooms
		ORDER BY last_activity DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var rooms []*Room
	for rows.Next() {
		room := &Room{}
		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.OwnerID,
			&room.Created,
			&room.LastActivity,
		)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	
	return rooms, nil
}

// GetRoomByID retrieves a room by its ID
func (r *PostgresRepository) GetRoomByID(ctx context.Context, id string) (*Room, error) {
	query := `
		SELECT id, name, owner_id, created, last_activity
		FROM rooms
		WHERE id = $1
	`
	
	room := &Room{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.Name,
		&room.OwnerID,
		&room.Created,
		&room.LastActivity,
	)
	
	if err != nil {
		return nil, err
	}
	
	return room, nil
}

// UpdateRoomActivity updates the last_activity timestamp for a room
func (r *PostgresRepository) UpdateRoomActivity(ctx context.Context, roomID string) error {
	query := `
		UPDATE rooms
		SET last_activity = $1
		WHERE id = $2
	`
	
	_, err := r.db.ExecContext(ctx, query, time.Now(), roomID)
	return err
}
