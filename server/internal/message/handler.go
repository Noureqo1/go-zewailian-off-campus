package message

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateRoomRequest struct {
	ID      string `json:"id" binding:"required"`
	Name    string `json:"name" binding:"required"`
	OwnerID string `json:"ownerId" binding:"required"`
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetMessages retrieves messages for a specific room
func (h *Handler) GetMessages(c *gin.Context) {
	roomID := c.Param("roomId")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room ID is required"})
		return
	}

	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "50")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 50
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	// Get messages
	messages, err := h.service.GetMessagesByRoom(c.Request.Context(), roomID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve messages"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

// GetMessage retrieves a message by its ID
func (h *Handler) GetMessage(c *gin.Context) {
	messageID := c.Param("messageId")
	if messageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message ID is required"})
		return
	}

	// Get message
	message, err := h.service.GetMessageByID(c.Request.Context(), messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": message})
}

// CreateRoom creates a new chat room
func (h *Handler) CreateRoom(c *gin.Context) {
	var request CreateRoomRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create room
	room, err := h.service.CreateRoom(c.Request.Context(), request.ID, request.Name, request.OwnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create room"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"room": room})
}

// GetRooms retrieves all available chat rooms
func (h *Handler) GetRooms(c *gin.Context) {
	// Get rooms
	rooms, err := h.service.GetRooms(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve rooms"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

// GetRoom retrieves a room by its ID
func (h *Handler) GetRoom(c *gin.Context) {
	roomID := c.Param("roomId")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "room ID is required"})
		return
	}

	// Get room
	room, err := h.service.GetRoomByID(c.Request.Context(), roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve room"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"room": room})
}
