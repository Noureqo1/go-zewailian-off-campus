package message

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const (
	// Cache keys and prefixes
	messageKeyPrefix = "message:"
	roomKeyPrefix    = "room:"
	roomListKey      = "rooms:list"
	sessionKeyPrefix = "session:"
	
	// Default expiration times
	defaultMessageExpiration = 24 * time.Hour
	defaultRoomExpiration    = 72 * time.Hour
	defaultSessionExpiration = 24 * time.Hour
)

// Cache defines the interface for caching message data
type Cache interface {
	// Message operations
	CacheMessage(ctx context.Context, message *Message) error
	GetCachedMessage(ctx context.Context, id string) (*Message, error)
	CacheRoomMessages(ctx context.Context, roomID string, messages []*Message) error
	GetCachedRoomMessages(ctx context.Context, roomID string) ([]*Message, error)
	
	// Room operations
	CacheRoom(ctx context.Context, room *Room) error
	GetCachedRoom(ctx context.Context, id string) (*Room, error)
	CacheRoomList(ctx context.Context, rooms []*Room) error
	GetCachedRoomList(ctx context.Context) ([]*Room, error)
	
	// Session operations
	SetSession(ctx context.Context, userID, sessionData string, expiration time.Duration) error
	GetSession(ctx context.Context, userID string) (string, error)
	DeleteSession(ctx context.Context, userID string) error
}

// RedisCache implements the Cache interface using Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache creates a new Redis cache
func NewRedisCache(addr string, password string, db int) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	
	return &RedisCache{
		client: client,
	}
}

// CacheMessage stores a message in Redis
func (c *RedisCache) CacheMessage(ctx context.Context, message *Message) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	
	key := fmt.Sprintf("%s%s", messageKeyPrefix, message.ID)
	return c.client.Set(ctx, key, data, defaultMessageExpiration).Err()
}

// GetCachedMessage retrieves a message from Redis
func (c *RedisCache) GetCachedMessage(ctx context.Context, id string) (*Message, error) {
	key := fmt.Sprintf("%s%s", messageKeyPrefix, id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	
	var message Message
	if err := json.Unmarshal(data, &message); err != nil {
		return nil, err
	}
	
	return &message, nil
}

// CacheRoomMessages stores a list of messages for a room in Redis
func (c *RedisCache) CacheRoomMessages(ctx context.Context, roomID string, messages []*Message) error {
	data, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	
	key := fmt.Sprintf("%s%s:messages", roomKeyPrefix, roomID)
	return c.client.Set(ctx, key, data, defaultMessageExpiration).Err()
}

// GetCachedRoomMessages retrieves messages for a room from Redis
func (c *RedisCache) GetCachedRoomMessages(ctx context.Context, roomID string) ([]*Message, error) {
	key := fmt.Sprintf("%s%s:messages", roomKeyPrefix, roomID)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	
	var messages []*Message
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil, err
	}
	
	return messages, nil
}

// CacheRoom stores a room in Redis
func (c *RedisCache) CacheRoom(ctx context.Context, room *Room) error {
	data, err := json.Marshal(room)
	if err != nil {
		return err
	}
	
	key := fmt.Sprintf("%s%s", roomKeyPrefix, room.ID)
	return c.client.Set(ctx, key, data, defaultRoomExpiration).Err()
}

// GetCachedRoom retrieves a room from Redis
func (c *RedisCache) GetCachedRoom(ctx context.Context, id string) (*Room, error) {
	key := fmt.Sprintf("%s%s", roomKeyPrefix, id)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	
	var room Room
	if err := json.Unmarshal(data, &room); err != nil {
		return nil, err
	}
	
	return &room, nil
}

// CacheRoomList stores the list of rooms in Redis
func (c *RedisCache) CacheRoomList(ctx context.Context, rooms []*Room) error {
	data, err := json.Marshal(rooms)
	if err != nil {
		return err
	}
	
	return c.client.Set(ctx, roomListKey, data, defaultRoomExpiration).Err()
}

// GetCachedRoomList retrieves the list of rooms from Redis
func (c *RedisCache) GetCachedRoomList(ctx context.Context) ([]*Room, error) {
	data, err := c.client.Get(ctx, roomListKey).Bytes()
	if err != nil {
		return nil, err
	}
	
	var rooms []*Room
	if err := json.Unmarshal(data, &rooms); err != nil {
		return nil, err
	}
	
	return rooms, nil
}

// SetSession stores a user session in Redis
func (c *RedisCache) SetSession(ctx context.Context, userID, sessionData string, expiration time.Duration) error {
	key := fmt.Sprintf("%s%s", sessionKeyPrefix, userID)
	
	if expiration == 0 {
		expiration = defaultSessionExpiration
	}
	
	return c.client.Set(ctx, key, sessionData, expiration).Err()
}

// GetSession retrieves a user session from Redis
func (c *RedisCache) GetSession(ctx context.Context, userID string) (string, error) {
	key := fmt.Sprintf("%s%s", sessionKeyPrefix, userID)
	return c.client.Get(ctx, key).Result()
}

// DeleteSession removes a user session from Redis
func (c *RedisCache) DeleteSession(ctx context.Context, userID string) error {
	key := fmt.Sprintf("%s%s", sessionKeyPrefix, userID)
	return c.client.Del(ctx, key).Err()
}
