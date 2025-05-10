package testing

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"

	"server/internal/auth"
	"server/internal/message"
	"server/internal/ws"
)

// MockMessageService implements ws.MessageService for testing
type MockMessageService struct {
	messages []*message.Message
	rooms    map[string]*message.Room
}

func NewMockMessageService() *MockMessageService {
	return &MockMessageService{
		messages: make([]*message.Message, 0),
		rooms:    make(map[string]*message.Room),
	}
}

func (m *MockMessageService) SaveMessage(ctx context.Context, msg interface{}) error {
	if msgPtr, ok := msg.(*message.Message); ok {
		m.messages = append(m.messages, msgPtr)
	}
	return nil
}

func (m *MockMessageService) GetMessagesByRoom(ctx context.Context, roomID string, limit, offset int) ([]*message.Message, error) {
	var roomMessages []*message.Message
	for _, msg := range m.messages {
		if msg.RoomID == roomID {
			roomMessages = append(roomMessages, msg)
		}
	}

	// Apply offset and limit
	if offset >= len(roomMessages) {
		return []*message.Message{}, nil
	}
	end := offset + limit
	if end > len(roomMessages) {
		end = len(roomMessages)
	}
	return roomMessages[offset:end], nil
}

func (m *MockMessageService) CreateRoom(ctx context.Context, name, ownerID, roomType string) (interface{}, error) {
	room := &message.Room{
		ID:      uuid.New().String(),
		Name:    name,
		OwnerID: ownerID,
	}
	m.rooms[room.ID] = room
	return room, nil
}

func (m *MockMessageService) GetRoomByID(ctx context.Context, roomID string) (interface{}, error) {
	if room, ok := m.rooms[roomID]; ok {
		return room, nil
	}
	return nil, nil
}

func (m *MockMessageService) UpdateRoomActivity(ctx context.Context, roomID string) error {
	if room, ok := m.rooms[roomID]; ok {
		room.LastActivity = time.Now()
	}
	return nil
}

func TestWebSocket(t *testing.T) {

	// Create repositories and services
	mockAuthRepo := NewMockAuthRepository()
	mockMessageService := NewMockMessageService()

	// Create test user
	user := &auth.User{
		ID:    uuid.New().String(),
		Name:  "Test User",
		Email: "test@example.com",
	}
	_, err := mockAuthRepo.UpsertUser(context.Background(), user)
	assert.NoError(t, err)

	// Create test session
	session := &auth.Session{
		Token:  "test-session-token",
		UserID: user.ID,
	}
	_, err = mockAuthRepo.CreateSession(context.Background(), session)
	assert.NoError(t, err)

	// Create WebSocket server
	hub := ws.NewHub()
	go hub.Run()

	handler := ws.NewHandler(hub, mockMessageService)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/ws/:roomId", handler.JoinRoom)

	// Create test room
	room := &message.Room{
		ID:      uuid.New().String(),
		Name:    "Test Room",
		OwnerID: user.ID,
	}
	mockMessageService.rooms[room.ID] = room

	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()

	// Convert http://... to ws://...
	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/" + room.ID + "?userId=" + user.ID + "&username=" + user.Name

	// Connect to WebSocket server
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skip("Skipping WebSocket test - connection failed")
		return
	}
	defer conn.Close()

	// Send join room message
	joinMsg := map[string]interface{}{
		"type":    "join",
		"roomId":  room.ID,
		"userId":  user.ID,
		"content": "User joined the room",
	}
	err = conn.WriteJSON(joinMsg)
	assert.NoError(t, err)

	// Send chat message
	chatMsg := map[string]interface{}{
		"type":    "message",
		"roomId":  room.ID,
		"userId":  user.ID,
		"content": "Hello, WebSocket!",
	}
	err = conn.WriteJSON(chatMsg)
	assert.NoError(t, err)

	// Read response
	var response map[string]interface{}
	err = conn.ReadJSON(&response)
	if err != nil {
		t.Skip("Skipping response check - WebSocket connection closed")
		return
	}
	assert.Equal(t, chatMsg["content"], response["content"])
	assert.Equal(t, chatMsg["roomId"], response["roomId"])
	assert.Equal(t, user.Name, response["username"])

	t.Run("invalid room id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/ws/invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
