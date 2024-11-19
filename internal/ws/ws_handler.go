package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"onlineChat/internal/users"
	"strconv"
)

type Handler struct {
	Hub     *Hub
	service *ChatService
}

type RoomReq struct {
	Name string `json:"name"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) ServeWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	chatID, err := strconv.Atoi(c.Param("chatID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	clientID, err := users.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	client := &Client{
		ID:         clientID,
		ChatID:     chatID,
		Connection: conn,
		Message:    make(chan *Message),
	}

	message := &Message{
		ChatID: chatID,
	}

	h.Hub.Register <- client
	h.Hub.Broadcast <- message

	go client.readMessage(h)
	go client.writeMessage(h)
}

func (h *Handler) CreateChat(c *gin.Context) {
	var input RoomReq
	userID, err := users.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err,
		})
	}

	err = c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	room, err := h.service.CreateChat(input.Name, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"room created": room,
	})
}

func (h *Handler) JoinChat(c *gin.Context) {
	userID, err := users.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err,
		})
	}

	chatID, err := strconv.Atoi(c.Param("chatID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid chat id",
		})
	}

	chat, err := h.service.JoinChat(userID, chatID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, chat)
}

func (h *Handler) GetAllChats(c *gin.Context) {
	var chats []Chat
	for _, ch := range h.Hub.Chats {
		chats = append(chats, Chat{
			ID:   ch.ID,
			Name: ch.Name,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"chats": chats,
	})
}

func (h *Handler) GetClientsByChatID(c *gin.Context) {
	chatID, err := strconv.Atoi(c.Param("chatID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}

	_, ok := h.Hub.Chats[chatID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "the chat doesnt exist",
		})
	}
	var clients []users.UserResponse
	for _, client := range h.Hub.Chats[chatID].Clients {
		clients = append(clients, users.UserResponse{
			ID:       client.ID,
			Username: client.Username,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"clients": clients,
	})
}
