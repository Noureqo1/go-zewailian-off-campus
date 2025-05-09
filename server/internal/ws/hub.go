package ws

import "time"

type Room struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Clients      map[string]*Client `json:"clients"`
	OwnerID      string             `json:"owner_id,omitempty"`
	Created      time.Time          `json:"created,omitempty"`
	LastActivity time.Time          `json:"last_activity,omitempty"`
}

type Hub struct {
	Rooms             map[string]*Room
	Register          chan *Client
	Unregister        chan *Client
	Broadcast         chan *Message
	UpdateClientStatus chan *Client       // Channel for client status updates (typing, etc.)
	PrivateMessage     chan *Message      // Channel for private messages between users
}

func NewHub() *Hub {
	return &Hub{
		Rooms:              make(map[string]*Room),
		Register:           make(chan *Client),
		Unregister:         make(chan *Client),
		Broadcast:          make(chan *Message, 5),
		UpdateClientStatus: make(chan *Client, 5),
		PrivateMessage:     make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register: //join
			if _, ok := h.Rooms[cl.RoomID]; ok {
				r := h.Rooms[cl.RoomID]

				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
				}
			}
		case cl := <-h.Unregister:
			if _, ok := h.Rooms[cl.RoomID]; ok {
				if _, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					if len(h.Rooms[cl.RoomID].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  "user left the chat",
							RoomID:   cl.RoomID,
							Username: cl.Username,
						}
					}

					delete(h.Rooms[cl.RoomID].Clients, cl.ID)
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			if _, ok := h.Rooms[m.RoomID]; ok {
				// Update room's last activity timestamp
				h.Rooms[m.RoomID].LastActivity = time.Now()

				for _, cl := range h.Rooms[m.RoomID].Clients {
					cl.Message <- m
				}
			}
			
		case cl := <-h.UpdateClientStatus:
			// Handle client status updates (typing indicators, etc.)
			if _, ok := h.Rooms[cl.RoomID]; ok {
				// Update the client in the room
				if existingClient, ok := h.Rooms[cl.RoomID].Clients[cl.ID]; ok {
					existingClient.IsTyping = cl.IsTyping
					existingClient.LastActive = time.Now()
					
					// Notify all clients in the room about the status change
					statusMsg := &Message{
						Type:      MessageTypeTyping,
						Content:   map[bool]string{true: "true", false: "false"}[cl.IsTyping],
						RoomID:    cl.RoomID,
						Username:  cl.Username,
						Timestamp: time.Now(),
					}
					
					for _, client := range h.Rooms[cl.RoomID].Clients {
						// Don't send typing notification back to the sender
						if client.ID != cl.ID {
							client.Message <- statusMsg
						}
					}
				}
			}
			
		case m := <-h.PrivateMessage:
			// Handle private messages between users
			if _, ok := h.Rooms[m.RoomID]; ok {
				// Find the recipient client
				for _, cl := range h.Rooms[m.RoomID].Clients {
					if cl.Username == m.Recipient {
						// Send the private message to the recipient
						cl.Message <- m
						
						// Also send a copy back to the sender
						for _, sender := range h.Rooms[m.RoomID].Clients {
							if sender.Username == m.Username {
								sender.Message <- m
								break
							}
						}
						break
					}
				}
			}
		}
	}
}
