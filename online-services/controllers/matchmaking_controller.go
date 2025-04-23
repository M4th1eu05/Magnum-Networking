package controllers

import (
	"net/http"
	"online-services/database"
	"online-services/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GameServerInfo struct {
	IP	   string `form:"ip" json:"ip" binding:"required"`
	Port   int    `form:"port" json:"port" binding:"required"`
}

func RegisterServer(c *gin.Context) {
	var serverInfo GameServerInfo
	if err := c.ShouldBindJSON(&serverInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	server := models.GameServer{
		IP:   serverInfo.IP,
		Port: serverInfo.Port,
		Status: "available",
	}

	var existingServer models.GameServer
	if err := database.DB.Where("ip = ? AND port = ?", server.IP, server.Port).First(&existingServer).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Server already registered"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := database.DB.Create(&server).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register server"})
		return
	}

	c.JSON(http.StatusCreated, server)
}

// JoinQueue adds a player to the matchmaking queue
func JoinQueue(c *gin.Context) {
	
}

// LeaveQueue removes a player from the matchmaking queue
func LeaveQueue(c *gin.Context) {
	
}

// Internal function to match players
func matchPlayers() {
	
}
