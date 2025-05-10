package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Picture      string    `json:"picture,omitempty"`
	PasswordHash string    `json:"-"`
	GoogleID     string    `json:"google_id,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Service interface {
	Signup(ctx context.Context, req *SignupRequest) (*User, error)
	Login(ctx context.Context, req *LoginRequest) (*User, error)
	UpsertUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	CreateSession(ctx context.Context, userID string) (*Session, error)
	GetUserBySession(ctx context.Context, token string) (*User, error)
	DeleteSession(ctx context.Context, token string) error
}

type DefaultService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &DefaultService{repo: repo}
}

func (s *DefaultService) Signup(ctx context.Context, req *SignupRequest) (*User, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	// Create user
	user := &User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	return s.repo.UpsertUser(ctx, user)
}

func (s *DefaultService) Login(ctx context.Context, req *LoginRequest) (*User, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	return user, nil
}

func (s *DefaultService) UpsertUser(ctx context.Context, user *User) (*User, error) {
	if user.ID == "" {
		user.ID = uuid.New().String()
		user.CreatedAt = time.Now()
	}
	return s.repo.UpsertUser(ctx, user)
}

func (s *DefaultService) GetUserByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *DefaultService) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *DefaultService) CreateSession(ctx context.Context, userID string) (*Session, error) {
	session := &Session{
		ID:        uuid.New().String(),
		UserID:    userID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 1 week
	}
	return s.repo.CreateSession(ctx, session)
}

func (s *DefaultService) GetUserBySession(ctx context.Context, token string) (*User, error) {
	session, err := s.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if time.Now().After(session.ExpiresAt) {
		s.repo.DeleteSession(ctx, token)
		return nil, err
	}

	return s.repo.GetUserByID(ctx, session.UserID)
}

func (s *DefaultService) DeleteSession(ctx context.Context, token string) error {
	return s.repo.DeleteSession(ctx, token)
}
