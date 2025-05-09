package message

import (
	"context"
	"fmt"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/sony/gobreaker"
)

// CircuitBreakerConfig holds configuration for circuit breaker
type CircuitBreakerConfig struct {
	Name          string
	MaxRequests   uint32
	Interval      time.Duration
	Timeout       time.Duration
	ReadyToTrip   func(counts gobreaker.Counts) bool
	OnStateChange func(name string, from gobreaker.State, to gobreaker.State)
}

// RetryConfig holds configuration for retry mechanism
type RetryConfig struct {
	MaxElapsedTime time.Duration
	MaxInterval    time.Duration
	InitialInterval time.Duration
}

// ResilientService wraps a Service with circuit breaker and retry mechanisms
type ResilientService struct {
	service Service
	messageBreaker *gobreaker.CircuitBreaker
	roomBreaker    *gobreaker.CircuitBreaker
	retryConfig    RetryConfig
}

// NewResilientService creates a new resilient service wrapper
func NewResilientService(service Service, cbConfig CircuitBreakerConfig, retryConfig RetryConfig) *ResilientService {
	messageCBSettings := gobreaker.Settings{
		Name:          cbConfig.Name + "-messages",
		MaxRequests:   cbConfig.MaxRequests,
		Interval:      cbConfig.Interval,
		Timeout:       cbConfig.Timeout,
		ReadyToTrip:   cbConfig.ReadyToTrip,
		OnStateChange: cbConfig.OnStateChange,
	}

	roomCBSettings := gobreaker.Settings{
		Name:          cbConfig.Name + "-rooms",
		MaxRequests:   cbConfig.MaxRequests,
		Interval:      cbConfig.Interval,
		Timeout:       cbConfig.Timeout,
		ReadyToTrip:   cbConfig.ReadyToTrip,
		OnStateChange: cbConfig.OnStateChange,
	}

	return &ResilientService{
		service:        service,
		messageBreaker: gobreaker.NewCircuitBreaker(messageCBSettings),
		roomBreaker:    gobreaker.NewCircuitBreaker(roomCBSettings),
		retryConfig:    retryConfig,
	}
}

// executeWithResilience executes a function with circuit breaker and retry mechanisms
func (rs *ResilientService) executeWithResilience(ctx context.Context, breaker *gobreaker.CircuitBreaker, operation func() (interface{}, error)) (interface{}, error) {
	exponentialBackoff := backoff.NewExponentialBackOff()
	exponentialBackoff.MaxElapsedTime = rs.retryConfig.MaxElapsedTime
	exponentialBackoff.MaxInterval = rs.retryConfig.MaxInterval
	exponentialBackoff.InitialInterval = rs.retryConfig.InitialInterval

	var result interface{}
	var err error

	retryOperation := func() error {
		result, err = breaker.Execute(func() (interface{}, error) {
			return operation()
		})
		return err
	}

	err = backoff.Retry(retryOperation, backoff.WithContext(exponentialBackoff, ctx))
	if err != nil {
		return nil, fmt.Errorf("operation failed after retries: %w", err)
	}

	return result, nil
}

// SaveMessage implements Service with resilience
func (rs *ResilientService) SaveMessage(ctx context.Context, message *Message) error {
	_, err := rs.executeWithResilience(ctx, rs.messageBreaker, func() (interface{}, error) {
		return nil, rs.service.SaveMessage(ctx, message)
	})
	return err
}

// GetMessagesByRoom implements Service with resilience
func (rs *ResilientService) GetMessagesByRoom(ctx context.Context, roomID string, limit, offset int) ([]*Message, error) {
	result, err := rs.executeWithResilience(ctx, rs.messageBreaker, func() (interface{}, error) {
		return rs.service.GetMessagesByRoom(ctx, roomID, limit, offset)
	})
	if err != nil {
		return nil, err
	}
	return result.([]*Message), nil
}

// GetMessageByID implements Service with resilience
func (rs *ResilientService) GetMessageByID(ctx context.Context, id string) (*Message, error) {
	result, err := rs.executeWithResilience(ctx, rs.messageBreaker, func() (interface{}, error) {
		return rs.service.GetMessageByID(ctx, id)
	})
	if err != nil {
		return nil, err
	}
	return result.(*Message), nil
}

// CreateRoom implements Service with resilience
func (rs *ResilientService) CreateRoom(ctx context.Context, id, name, ownerID string) (*Room, error) {
	result, err := rs.executeWithResilience(ctx, rs.roomBreaker, func() (interface{}, error) {
		return rs.service.CreateRoom(ctx, id, name, ownerID)
	})
	if err != nil {
		return nil, err
	}
	return result.(*Room), nil
}

// GetRooms implements Service with resilience
func (rs *ResilientService) GetRooms(ctx context.Context) ([]*Room, error) {
	result, err := rs.executeWithResilience(ctx, rs.roomBreaker, func() (interface{}, error) {
		return rs.service.GetRooms(ctx)
	})
	if err != nil {
		return nil, err
	}
	return result.([]*Room), nil
}

// GetRoomByID implements Service with resilience
func (rs *ResilientService) GetRoomByID(ctx context.Context, id string) (*Room, error) {
	result, err := rs.executeWithResilience(ctx, rs.roomBreaker, func() (interface{}, error) {
		return rs.service.GetRoomByID(ctx, id)
	})
	if err != nil {
		return nil, err
	}
	return result.(*Room), nil
}

// UpdateRoomActivity implements Service with resilience
func (rs *ResilientService) UpdateRoomActivity(ctx context.Context, roomID string) error {
	_, err := rs.executeWithResilience(ctx, rs.roomBreaker, func() (interface{}, error) {
		return nil, rs.service.UpdateRoomActivity(ctx, roomID)
	})
	return err
}

// Session operations don't need circuit breakers as they're Redis-only operations
func (rs *ResilientService) SetUserSession(ctx context.Context, userID, sessionData string, expiration time.Duration) error {
	return rs.service.SetUserSession(ctx, userID, sessionData, expiration)
}

func (rs *ResilientService) GetUserSession(ctx context.Context, userID string) (string, error) {
	return rs.service.GetUserSession(ctx, userID)
}

func (rs *ResilientService) DeleteUserSession(ctx context.Context, userID string) error {
	return rs.service.DeleteUserSession(ctx, userID)
}
