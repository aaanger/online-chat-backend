package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	maxMessageSize = 512
	pongWait       = 1 * time.Minute
	pingPeriod     = 50 * time.Second
)

type Client struct {
	Connection *websocket.Conn
	ID         int    `json:"id"`
	ChatID     int    `json:"chatID"`
	Username   string `json:"username"`
	Message    chan *Message
}

type Message struct {
	ChatID   int    `json:"id"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

func (c *Client) readMessage(h *Handler) {
	defer func() {
		h.Hub.Unregister <- c
		c.Connection.Close()
	}()

	c.Connection.SetReadLimit(maxMessageSize)
	c.Connection.SetPongHandler(func(string) error {
		return c.Connection.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, message, err := c.Connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		msg := &Message{
			ChatID:   c.ChatID,
			Username: c.Username,
			Content:  string(message),
		}

		err = h.service.SaveChat(c.ID, msg)
		if err != nil {
			log.Printf("error saving message: %v", err)
		}

		h.Hub.Broadcast <- msg

	}
}

func (c *Client) writeMessage(h *Handler) {
	defer func() {
		h.Hub.Unregister <- c
		c.Connection.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		c.Connection.WriteJSON(message)
	}
}
