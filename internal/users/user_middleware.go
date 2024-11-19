package users

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *Handler) UserIdentity(c *gin.Context) {
	header := c.GetHeader("Authorization")

	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "empty auth header",
		})
		return
	}

	headerParts := strings.Split(header, ", ")
	if len(headerParts) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid auth header",
		})
		return
	}

	userID, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	c.Set("userID", userID)
}

func GetUserID(c *gin.Context) (int, error) {
	id, ok := c.Get("userID")
	if !ok {
		return 0, fmt.Errorf("invalid user id")
	}

	userID, ok := id.(int)
	if !ok {
		return 0, fmt.Errorf("user id is not of type int")
	}

	return userID, nil
}
