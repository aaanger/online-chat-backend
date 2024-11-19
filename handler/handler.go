package handler

import (
	"github.com/gin-gonic/gin"
	"onlineChat/internal/users"
	"onlineChat/internal/ws"
)

func PathHandler(userHandler *users.Handler, wsHandler *ws.Handler) *gin.Engine {
	r := gin.New()

	r.POST("/signup", userHandler.SignUp)
	r.POST("/signin", userHandler.SignIn)
	r.GET("/logout", userHandler.Logout)

	chat := r.Group("/chat", userHandler.UserIdentity)
	{
		chat.GET("/:chatID", wsHandler.ServeWS)
		chat.GET("/all", wsHandler.GetAllChats)
		chat.GET("/:chatID/clients", wsHandler.GetClientsByChatID)
	}

	return r
}
