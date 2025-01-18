package ws

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"onlineChat/pkg/redis"
	"sync"
	"time"
)

type Hub struct {
	Chats      map[int]*Chat
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	Redis      *redis.RedisClient
	mu         sync.Mutex
}

type Chat struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
	Clients   map[int]*Client
}

func NewHub(redisCfg redis.RedisConfig) *Hub {
	return &Hub{
		Chats:      make(map[int]*Chat),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message),
		Redis: redis.NewRedisClient(redis.RedisConfig{
			Address:  redisCfg.Address,
			Password: redisCfg.Password,
			DB:       redisCfg.DB,
		}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			chat, ok := h.Chats[client.ChatID]
			if ok {
				err := h.Redis.AddUser(chat.ID, client.ID)
				if err != nil {
					logrus.Warnf("Failed to add user %s to chat in Redis", client.Username)
					continue
				}
				chat.Clients[client.ID] = client

				h.Broadcast <- &Message{
					ChatID:   client.ChatID,
					Username: client.Username,
					Content:  fmt.Sprintf("User %s joined the chat", client.Username),
				}
			}
			h.mu.Unlock()
		case client := <-h.Unregister:
			h.mu.Lock()
			_, ok := h.Chats[client.ChatID].Clients[client.ID]
			if ok {
				err := h.Redis.RemoveUser(client.ChatID, client.ID)
				if err != nil {
					logrus.Warnf("Failed to remove user %s from chat in Redis", client.Username)
				}
				h.Broadcast <- &Message{
					ChatID:   client.ChatID,
					Username: client.Username,
					Content:  fmt.Sprintf("User %s left the chat", client.Username),
				}
				delete(h.Chats[client.ChatID].Clients, client.ID)
				close(client.Message)
			}
			h.mu.Unlock()
		case message := <-h.Broadcast:
			h.mu.Lock()
			chat, ok := h.Chats[message.ChatID]
			if ok {
				for _, client := range chat.Clients {
					client.Message <- message
				}
			}
			h.mu.Unlock()
		}
	}
}
