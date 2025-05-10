package testing

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"server/internal/auth"
)

// MockAuthRepository implements auth.Repository interface for testing
type MockAuthRepository struct {
	users    map[string]*auth.User
	sessions map[string]*auth.Session
}

func NewMockAuthRepository() *MockAuthRepository {
	return &MockAuthRepository{
		users:    make(map[string]*auth.User),
		sessions: make(map[string]*auth.Session),
	}
}

func (m *MockAuthRepository) UpsertUser(ctx context.Context, user *auth.User) (*auth.User, error) {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	m.users[user.Email] = user
	return user, nil
}

func (m *MockAuthRepository) GetUserByID(ctx context.Context, id string) (*auth.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*auth.User, error) {
	if user, exists := m.users[email]; exists {
		return user, nil
	}
	return nil, nil
}

func (m *MockAuthRepository) CreateSession(ctx context.Context, session *auth.Session) (*auth.Session, error) {
	m.sessions[session.Token] = session
	return session, nil
}

func (m *MockAuthRepository) GetSessionByToken(ctx context.Context, token string) (*auth.Session, error) {
	if session, exists := m.sessions[token]; exists {
		return session, nil
	}
	return nil, nil
}

func (m *MockAuthRepository) DeleteSession(ctx context.Context, token string) error {
	delete(m.sessions, token)
	return nil
}

func (m *MockAuthRepository) GetUserByGoogleID(ctx context.Context, googleID string) (*auth.User, error) {
	for _, user := range m.users {
		if user.GoogleID == googleID {
			return user, nil
		}
	}
	return nil, nil
}

func TestSignup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockRepo := NewMockAuthRepository()
	authService := auth.NewService(mockRepo)
	authHandler := auth.NewHandler(authService)
	r.POST("/signup", authHandler.Signup)

	tests := []struct {
		name       string
		payload    auth.SignupRequest
		wantStatus int
	}{
		{
			name: "valid signup",
			payload: auth.SignupRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "testpass123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "duplicate email",
			payload: auth.SignupRequest{
				Name:     "Another User",
				Email:    "test@example.com",
				Password: "testpass123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
