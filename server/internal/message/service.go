package message

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	// Message operations
	SaveMessage(ctx context.Context, message *Message) error
	GetMessagesByRoom(ctx context.Context, roomID string, limit, offset int) ([]*Message, error)
	GetMessageByID(ctx context.Context, id string) (*Message, error)

	CreateRoom(ctx context.Context, id, name, ownerID string) (*Room, error)
	GetRooms(ctx context.Context) ([]*Room, error)
	GetRoomByID(ctx context.Context, id string) (*Room, error)
	UpdateRoomActivity(ctx context.Context, roomID string) error

	SetUserSession(ctx context.Context, userID, sessionData string, expiration time.Duration) error
	GetUserSession(ctx context.Context, userID string) (string, error)
	DeleteUserSession(ctx context.Context, userID string) error
}

type DefaultService struct {
	repo  Repository
	cache Cache
}

func NewService(repo Repository, cache Cache) Service {
	return &DefaultService{
		repo:  repo,
		cache: cache,
	}
}

func (s *DefaultService) SaveMessage(ctx context.Context, message *Message) error {
	if message.ID == "" {
		message.ID = uuid.New().String()
	}

	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}

	if err := s.repo.SaveMessage(ctx, message); err != nil {
		return err
	}

	if err := s.cache.CacheMessage(ctx, message); err != nil {
	}

	if err := s.UpdateRoomActivity(ctx, message.RoomID); err != nil {
	}

	return nil
}
func (s *DefaultService) GetMessagesByRoom(ctx context.Context, roomID string, limit, offset int) ([]*Message, error) {
	cachedMessages, err := s.cache.GetCachedRoomMessages(ctx, roomID)
	if err == nil && len(cachedMessages) > 0 {
		return cachedMessages, nil
	}

	messages, err := s.repo.GetMessagesByRoom(ctx, roomID, limit, offset)
	if err != nil {
		return nil, err
	}

	if len(messages) > 0 {
		if err := s.cache.CacheRoomMessages(ctx, roomID, messages); err != nil {
		}
	}

	return messages, nil
}

func (s *DefaultService) GetMessageByID(ctx context.Context, id string) (*Message, error) {
	cachedMessage, err := s.cache.GetCachedMessage(ctx, id)
	if err == nil {
		return cachedMessage, nil
	}

	message, err := s.repo.GetMessageByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.cache.CacheMessage(ctx, message); err != nil {
	}

	return message, nil
}

func (s *DefaultService) CreateRoom(ctx context.Context, id, name, ownerID string) (*Room, error) {
	room := &Room{
		ID:           id,
		Name:         name,
		OwnerID:      ownerID,
		Created:      time.Now(),
		LastActivity: time.Now(),
	}

	// Save to database
	if err := s.repo.CreateRoom(ctx, room); err != nil {
		return nil, err
	}

	// Update cache
	if err := s.cache.CacheRoom(ctx, room); err != nil {
		// Log error but don't fail the operation
		// log.Printf("Error caching room: %v", err)
	}

	// Invalidate room list cache by getting fresh data and updating
	rooms, err := s.repo.GetRooms(ctx)
	if err == nil {
		s.cache.CacheRoomList(ctx, rooms)
	}

	return room, nil
}

// GetRooms retrieves all available chat rooms with caching
func (s *DefaultService) GetRooms(ctx context.Context) ([]*Room, error) {
	// Try to get from cache first
	cachedRooms, err := s.cache.GetCachedRoomList(ctx)
	if err == nil && len(cachedRooms) > 0 {
		return cachedRooms, nil
	}

	// Get from database
	rooms, err := s.repo.GetRooms(ctx)
	if err != nil {
		return nil, err
	}

	// Update cache
	if len(rooms) > 0 {
		if err := s.cache.CacheRoomList(ctx, rooms); err != nil {
			// Log error but don't fail the operation
			// log.Printf("Error caching room list: %v", err)
		}
	}

	return rooms, nil
}

// GetRoomByID retrieves a room by its ID with caching
func (s *DefaultService) GetRoomByID(ctx context.Context, id string) (*Room, error) {
	// Try to get from cache first
	cachedRoom, err := s.cache.GetCachedRoom(ctx, id)
	if err == nil {
		return cachedRoom, nil
	}

	// Get from database
	room, err := s.repo.GetRoomByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update cache
	if err := s.cache.CacheRoom(ctx, room); err != nil {
		// Log error but don't fail the operation
		// log.Printf("Error caching room: %v", err)
	}

	return room, nil
}

// UpdateRoomActivity updates the last_activity timestamp for a room
func (s *DefaultService) UpdateRoomActivity(ctx context.Context, roomID string) error {
	// Update in database
	if err := s.repo.UpdateRoomActivity(ctx, roomID); err != nil {
		return err
	}

	// Get updated room to refresh cache
	room, err := s.repo.GetRoomByID(ctx, roomID)
	if err == nil {
		s.cache.CacheRoom(ctx, room)
	}

	return nil
}

// SetUserSession stores a user session in Redis
func (s *DefaultService) SetUserSession(ctx context.Context, userID, sessionData string, expiration time.Duration) error {
	return s.cache.SetSession(ctx, userID, sessionData, expiration)
}

// GetUserSession retrieves a user session from Redis
func (s *DefaultService) GetUserSession(ctx context.Context, userID string) (string, error) {
	return s.cache.GetSession(ctx, userID)
}

// DeleteUserSession removes a user session from Redis
func (s *DefaultService) DeleteUserSession(ctx context.Context, userID string) error {
	return s.cache.DeleteSession(ctx, userID)
}
