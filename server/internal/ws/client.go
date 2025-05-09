package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	MessageTypeChat    MessageType = "chat"    // Regular chat message
	MessageTypeJoin    MessageType = "join"    // User joined notification
	MessageTypeLeave   MessageType = "leave"   // User left notification
	MessageTypeSystem  MessageType = "system"  // System message
	MessageTypePrivate MessageType = "private" // Private message to specific user
	MessageTypeError   MessageType = "error"   // Error message
	MessageTypeTyping  MessageType = "typing"  // User is typing notification
)

type Client struct {
	Conn       *websocket.Conn
	Message    chan *Message
	ID         string    `json:"id"`
	RoomID     string    `json:"roomId"`
	Username   string    `json:"username"`
	IsActive   bool      `json:"isActive"`   // Whether client is currently active
	JoinedAt   time.Time `json:"joinedAt"`   // When client joined
	LastActive time.Time `json:"lastActive"` // Last activity timestamp
	IsTyping   bool      `json:"isTyping"`   // Whether client is currently typing
}

// Message represents a message sent between clients
type Message struct {
	ID        string      `json:"id,omitempty"`        // Unique message ID
	Type      MessageType `json:"type"`                // Message type
	Content   string      `json:"content"`             // Message content
	RoomID    string      `json:"roomId"`              // Room ID
	Username  string      `json:"username"`            // Sender username
	Timestamp time.Time   `json:"timestamp"`           // Message timestamp
	Recipient string      `json:"recipient,omitempty"` // For private messages
}

// writeMessage handles sending messages to the client
func (c *Client) writeMessage() {
	// Close connection when function returns
	defer func() {
		c.Conn.Close()
	}()

	// Set ping handler
	c.Conn.SetPingHandler(func(string) error {

		c.LastActive = time.Now()
		return c.Conn.WriteControl(websocket.PongMessage, []byte{}, time.Now().Add(10*time.Second))
	})

	pingTicker := time.NewTicker(30 * time.Second)
	defer pingTicker.Stop()

	for {
		select {
		case message, ok := <-c.Message:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.LastActive = time.Now()

			err := c.Conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error writing message to client %s: %v", c.ID, err)
				return
			}

		case <-pingTicker.C:
			if err := c.Conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second)); err != nil {
				log.Printf("Error sending ping to client %s: %v", c.ID, err)
				return
			}
		}
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
	c.Conn.SetCloseHandler(func(code int, text string) error {
		log.Printf("Client %s connection closed: %d %s", c.ID, code, text)
		return nil
	})

	for {
		c.LastActive = time.Now()

		_, rawMessage, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message from client %s: %v", c.ID, err)
			}
			break
		}

		var parsedMsg Message
		if err := json.Unmarshal(rawMessage, &parsedMsg); err == nil {
			if parsedMsg.Type == "" {
				parsedMsg.Type = MessageTypeChat
			}
			if parsedMsg.RoomID == "" {
				parsedMsg.RoomID = c.RoomID
			}
			if parsedMsg.Username == "" {
				parsedMsg.Username = c.Username
			}
			if parsedMsg.Timestamp.IsZero() {
				parsedMsg.Timestamp = time.Now()
			}

			if parsedMsg.Type == MessageTypeTyping {
				c.IsTyping = parsedMsg.Content == "true"
				hub.UpdateClientStatus <- c
				continue
			}

			if parsedMsg.Type == MessageTypePrivate && parsedMsg.Recipient != "" {
				hub.PrivateMessage <- &parsedMsg
				continue
			}

			hub.Broadcast <- &parsedMsg
		} else {
			msg := &Message{
				Type:      MessageTypeChat,
				Content:   string(rawMessage),
				RoomID:    c.RoomID,
				Username:  c.Username,
				Timestamp: time.Now(),
			}

			hub.Broadcast <- msg
		}
	}
}
