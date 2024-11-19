package ws

import "fmt"

type Hub struct {
	Chats      map[int]*Chat
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

type Chat struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Message   string          `json:"message"`
	CreatedAt int             `json:"createdAt"`
	Clients   map[int]*Client `json:"clients"`
}

func NewHub() *Hub {
	return &Hub{
		Chats:      make(map[int]*Chat),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			_, ok := h.Chats[client.ChatID]
			if ok {
				chat := h.Chats[client.ChatID]

				_, ok := chat.Clients[client.ID]
				if ok {
					chat.Clients[client.ID] = client

					h.Broadcast <- &Message{
						ChatID:   client.ChatID,
						Username: client.Username,
						Content:  fmt.Sprintf("User %s joined the chat", client.Username),
					}
				}
			}
		case client := <-h.Unregister:
			_, ok := h.Chats[client.ChatID].Clients[client.ID]
			if ok {
				h.Broadcast <- &Message{
					ChatID:   client.ChatID,
					Username: client.Username,
					Content:  fmt.Sprintf("User %s left the chat", client.Username),
				}
				delete(h.Chats[client.ChatID].Clients, client.ID)
				close(client.Message)
			}
		case message := <-h.Broadcast:
			_, ok := h.Chats[message.ChatID]
			if ok {
				for _, client := range h.Chats[message.ChatID].Clients {
					client.Message <- message
				}
			}
		}
	}
}
