//package controllers
//
//import (
//	"net/http"
//	"online-services/database"
//	"online-services/models"
//	"time"
//
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt"
//)
//
//// JoinQueue adds a player to the matchmaking queue
//func JoinQueue(c *gin.Context) {
//	// Get user from JWT token
//	tokenString := c.GetHeader("Authorization")
//	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		return []byte("your_secret_key"), nil
//	})
//	claims := token.Claims.(jwt.MapClaims)
//	username := claims["username"].(string)
//
//	var user models.User
//	if err := database.DB.Preload("Stats").Where("username = ?", username).First(&user).Error; err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//		return
//	}
//
//	// Check if player is already in queue
//	var existingQueue models.QueuedPlayer
//	result := database.DB.Where("user_id = ? AND status = ?", user.ID, "queued").First(&existingQueue)
//	if result.RowsAffected > 0 {
//		c.JSON(http.StatusConflict, gin.H{"error": "Already in queue"})
//		return
//	}
//
//	// Add player to queue
//	queuedPlayer := models.QueuedPlayer{
//		UserID:   user.ID,
//		JoinedAt: time.Now().Unix(),
//		Status:   "queued",
//	}
//	if err := database.DB.Create(&queuedPlayer).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to join queue"})
//		return
//	}
//
//	// Try to match players
//	go matchPlayers()
//
//	c.JSON(http.StatusOK, gin.H{"message": "Joined queue successfully"})
//}
//
//// LeaveQueue removes a player from the matchmaking queue
//func LeaveQueue(c *gin.Context) {
//	// Get user from JWT token
//	tokenString := c.GetHeader("Authorization")
//	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		return []byte("your_secret_key"), nil
//	})
//	claims := token.Claims.(jwt.MapClaims)
//	username := claims["username"].(string)
//
//	var user models.User
//	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//		return
//	}
//
//	// Remove player from queue
//	if err := database.DB.Model(&models.QueuedPlayer{}).Where("user_id = ? AND status = ?", user.ID, "queued").Update("status", "cancelled").Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to leave queue"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Left queue successfully"})
//}
//
//// RegisterServer registers a new game server
//func RegisterServer(c *gin.Context) {
//	var server models.GameServer
//	if err := c.ShouldBindJSON(&server); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	server.Status = "available"
//	if err := database.DB.Create(&server).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register server"})
//		return
//	}
//
//	c.JSON(http.StatusCreated, server)
//}
//
//// UpdateServerStatus updates a game server's status
//func UpdateServerStatus(c *gin.Context) {
//	var server models.GameServer
//	id := c.Param("id")
//
//	if err := database.DB.First(&server, id).Error; err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
//		return
//	}
//
//	var statusUpdate struct {
//		Status string `json:"status" binding:"required"`
//	}
//
//	if err := c.ShouldBindJSON(&statusUpdate); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	server.Status = statusUpdate.Status
//	if err := database.DB.Save(&server).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update server status"})
//		return
//	}
//
//	c.JSON(http.StatusOK, server)
//}
//
//// Internal function to match players
//func matchPlayers() {
//	var queuedPlayers []models.QueuedPlayer
//	database.DB.Preload("User.Stats").Where("status = ?", "queued").Order("joined_at asc").Find(&queuedPlayers)
//
//	if len(queuedPlayers) < 2 {
//		return
//	}
//
//	// Find available server
//	var server models.GameServer
//	if err := database.DB.Where("status = ?", "available").First(&server).Error; err != nil {
//		return
//	}
//
//	// Create new game
//	game := models.Game{
//		ServerID:  server.ID,
//		Status:    "waiting",
//		StartTime: time.Now().Unix(),
//	}
//
//	if err := database.DB.Create(&game).Error; err != nil {
//		return
//	}
//
//	// Match first two players in queue (simple matching for now)
//	player1 := queuedPlayers[0]
//	player2 := queuedPlayers[1]
//
//	// Update their queue status
//	database.DB.Model(&models.QueuedPlayer{}).Where("id IN ?", []uint{player1.ID, player2.ID}).Update("status", "matched")
//
//	// Add players to game
//	if err := database.DB.Model(&game).Association("Players").Append([]models.User{player1.User, player2.User}); err != nil {
//		return
//	}
//
//	// Update server status
//	server.Status = "in_game"
//	database.DB.Save(&server)
//}
