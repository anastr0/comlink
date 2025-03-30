package message

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *MessagesAPIHandler) CreateUserHandler(c *gin.Context) {
	// create user
	var user_json User

	if err := c.ShouldBindJSON(&user_json); err != nil {
		log.Printf("error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		result := h.db.Create(&user_json)

		if result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		} else {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "user": user_json})
		}
	}
}
func (h *MessagesAPIHandler) GetUsersHandler(c *gin.Context) {
	// get all users
	var users []User
	result := h.db.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "users": users})
	}
}
