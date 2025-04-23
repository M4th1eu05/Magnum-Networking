package controllers

import (
	"errors"
	"log"
	"net/http"
	"online-services/database"
	"online-services/models"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GameServerInfo struct {
	IP   string `form:"ip" json:"ip" binding:"required"`
	Port int    `form:"port" json:"port" binding:"required"`
}

var queue = make([]models.User, 0)
var queueMutex sync.Mutex

func RegisterServer(c *gin.Context) {
	var serverInfo GameServerInfo
	if err := c.ShouldBindJSON(&serverInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	server := models.GameServer{
		IP:     serverInfo.IP,
		Port:   serverInfo.Port,
		Status: "available",
	}

	var existingServer models.GameServer
	if err := database.DB.Where("ip = ? AND port = ?", server.IP, server.Port).First(&existingServer).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Server already registered"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
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
	UUID, exists := c.Get("UUID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID not found"})
		return
	}

	var user models.User
	if err := database.DB.Where("uuid = ?", UUID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	queueMutex.Lock()
	queue = append(queue, user)
	queueMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Joined matchmaking queue"})
}

// Pulling request from players
func QueueStatus(c *gin.Context) {
	UUID, exists := c.Get("UUID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID not found"})
		return
	}

	var user models.User
	if err := database.DB.Where("uuid = ?", UUID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	queueMutex.Lock()
	inQueue := false
	for _, queuedUser := range queue {
		if queuedUser.ID == user.ID {
			inQueue = true
			break
		}
	}
	queueMutex.Unlock()

	if inQueue {
		c.JSON(http.StatusOK, gin.H{"status": "in queue"})
	} else {
		// get game from user
		var game models.Game
		if err := database.DB.Where("game_id = ?", user.GameID).First(&game).Error; err != nil {

		}
		c.JSON(http.StatusOK, gin.H{"status": "not in queue"})
	}
}

// LeaveQueue removes a player from the matchmaking queue
func LeaveQueue(c *gin.Context) {
	UUID, exists := c.Get("UUID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID not found"})
		return
	}

	var user models.User
	if err := database.DB.Where("uuid = ?", UUID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	queueMutex.Lock()
	for i, queuedUser := range queue {
		if queuedUser.ID == user.ID {
			queue = append(queue[:i], queue[i+1:]...)
			break
		}
	}
	queueMutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"message": "Left matchmaking queue"})
}

// Internal function to match players
func matchPlayers() {
	for {
		time.Sleep(1 * time.Second) // Intervalle pour vérifier la file d'attente

		queueMutex.Lock()
		if len(queue) >= 4 {
			// Extraire deux joueurs de la file d'attente
			players := queue[:4]
			queue = queue[4:]

			createMatch(players)
		}
		queueMutex.Unlock()
	}
}

func createMatch(players []models.User) {
	game := models.Game{
		Players: players,
	}
	if err := database.DB.Create(&game).Error; err != nil {
		log.Fatal(err)
		return
	}

}
