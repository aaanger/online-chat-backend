package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"onlineChat/internal/users"
	"onlineChat/pkg/response"
	"strconv"
)

type Handler struct {
	Hub     *Hub
	service IChatService
}

func NewChatHandler(hub *Hub, service IChatService) *Handler {
	return &Handler{
		Hub:     hub,
		service: service,
	}
}

type ChatReq struct {
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
		response.Error(c, http.StatusInternalServerError, err.Error())
	}

	chatID, err := strconv.Atoi(c.Param("chatID"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid chat id")
	}

	clientID, err := users.GetUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user id not found")
		return
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
	var input ChatReq

	if input.Name == "" {
		response.Error(c, http.StatusBadRequest, "empty chat name")
	}

	userID, err := users.GetUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "user id not found")
	}

	err = c.BindJSON(&input)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid input parameters")
	}

	chat, err := h.service.CreateChat(input.Name, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "create chat error")
	}

	c.JSON(http.StatusOK, gin.H{
		"chat created": chat,
	})
}

func (h *Handler) JoinChat(c *gin.Context) {
	userID, err := users.GetUserID(c)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "invalid user id")
	}

	chatID, err := strconv.Atoi(c.Param("chatID"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid chat id")
	}

	_, err = h.service.JoinChat(userID, chatID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to join chat")
	}

	c.Redirect(http.StatusOK, "/:chatID")
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
		response.Error(c, http.StatusBadRequest, "invalid chat id")
	}

	_, ok := h.Hub.Chats[chatID]
	if !ok {
		response.Error(c, http.StatusNotFound, "chat does not exist")
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
