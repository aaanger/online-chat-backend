package users

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	service *UserService
}

func NewUserHandler(service *UserService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) SignUp(c *gin.Context) {
	var input User
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	res, err := h.service.CreateUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) SignIn(c *gin.Context) {
	var input User
	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	user, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	c.SetCookie("jwt", user.accessToken, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"login successful": user.Username,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"status": "logged out",
	})
}
