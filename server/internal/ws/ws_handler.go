package ws

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type MessageService interface {
	SaveMessage(ctx context.Context, message interface{}) error
	CreateRoom(ctx context.Context, id, name, ownerID string) (interface{}, error)
	UpdateRoomActivity(ctx context.Context, roomID string) error
}

type Handler struct {
	hub            *Hub
	messageService MessageService
}

func NewHandler(h *Hub, messageService MessageService) *Handler {
	return &Handler{
		hub:            h,
		messageService: messageService,
	}
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(c *gin.Context) {
	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:           req.ID,
		Name:         req.Name,
		Clients:      make(map[string]*Client),
		OwnerID:      "system",
		Created:      time.Now(),
		LastActivity: time.Now(),
	}

	if h.messageService != nil {
		_, err := h.messageService.CreateRoom(c.Request.Context(), req.ID, req.Name, "system")
		if err != nil {
			log.Printf("Error saving room to database: %v", err)
		}
	}

	c.JSON(http.StatusOK, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) JoinRoom(c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Param("roomId")
	clientID := c.Query("userId")
	username := c.Query("username")

	cl := &Client{
		Conn:     conn,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	m := &Message{
		Content:   "A new user has joined the room",
		RoomID:    roomID,
		Username:  username,
		Type:      MessageTypeJoin,
		Timestamp: time.Now(),
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	if h.messageService != nil {
		if err := h.messageService.UpdateRoomActivity(c.Request.Context(), roomID); err != nil {
			log.Printf("Error updating room activity: %v", err)
		}
		
		if err := h.messageService.SaveMessage(c.Request.Context(), m); err != nil {
			log.Printf("Error saving join message: %v", err)
		}
	}

	cl.messageService = h.messageService

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]RoomRes, 0)

	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	c.JSON(http.StatusOK, rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *gin.Context) {
	var clients []ClientRes
	roomId := c.Param("roomId")

	if _, ok := h.hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, c := range h.hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}
