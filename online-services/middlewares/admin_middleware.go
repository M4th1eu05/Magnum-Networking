package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online-services/database"
	"online-services/models"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid, exist := c.Get("UUID")
		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{"error": "UUID not found"})
			c.Abort()
			return
		}

		var user models.User
		if err := database.DB.Where("uuid = ?", uuid).First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		if user.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to access this resource"})
			c.Abort()
			return
		}

		c.Next()
	}
}
