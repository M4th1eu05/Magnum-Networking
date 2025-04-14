package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"online-services/controllers"
	"online-services/database"
	"online-services/middlewares"
	"online-services/models"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupMatchmakingRouter() *gin.Engine {
	r := gin.Default()
	authorized := r.Group("/", middlewares.AuthMiddleware())
	{
		authorized.POST("/queue/join", controllers.JoinQueue)
		authorized.POST("/queue/leave", controllers.LeaveQueue)
	}

	admin := r.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		admin.POST("/servers", controllers.RegisterServer)
		admin.PUT("/servers/:id/status", controllers.UpdateServerStatus)
	}
	return r
}

func TestJoinQueue(t *testing.T) {
	// Setup
	database.ConnectDB()
	defer database.CloseDB()

	// Create test user
	testUser := models.User{
		Username: "queuer",
		Password: "password123",
		Stats: models.Stats{
			Rank: models.Silver,
		},
	}
	database.DB.Create(&testUser)
	token, _ := middlewares.GenerateToken(testUser)

	router := setupMatchmakingRouter()

	// Test joining queue
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/queue/join", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify user is in queue
	var queuedPlayer models.QueuedPlayer
	result := database.DB.Where("user_id = ? AND status = ?", testUser.ID, "queued").First(&queuedPlayer)
	assert.Equal(t, int64(1), result.RowsAffected)
}

func TestLeaveQueue(t *testing.T) {
	// Setup
	database.ConnectDB()
	defer database.CloseDB()

	// Create test user
	testUser := models.User{Username: "queue_leaver", Password: "password123"}
	database.DB.Create(&testUser)

	// Add user to queue
	queuedPlayer := models.QueuedPlayer{
		UserID:   testUser.ID,
		JoinedAt: time.Now().Unix(),
		Status:   "queued",
	}
	database.DB.Create(&queuedPlayer)

	token, _ := middlewares.GenerateToken(testUser)
	router := setupMatchmakingRouter()

	// Test leaving queue
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/queue/leave", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify user's queue status is updated
	var updatedPlayer models.QueuedPlayer
	database.DB.First(&updatedPlayer, queuedPlayer.ID)
	assert.Equal(t, "cancelled", updatedPlayer.Status)
}

func TestRegisterServer(t *testing.T) {
	// Setup
	database.ConnectDB()
	defer database.CloseDB()

	// Create admin user
	adminUser := models.User{Username: "admin", Password: "admin123", IsAdmin: true}
	database.DB.Create(&adminUser)
	token, _ := middlewares.GenerateToken(adminUser)

	router := setupMatchmakingRouter()

	// Test server registration
	w := httptest.NewRecorder()
	serverData := `{
		"address": "localhost",
		"port": 8081,
		"max_players": 10,
		"status": "available"
	}`
	req, _ := http.NewRequest("POST", "/admin/servers", strings.NewReader(serverData))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	var response models.GameServer
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "localhost", response.Address)
	assert.Equal(t, "available", response.Status)
}

func TestMatchmaking(t *testing.T) {
	// Setup
	database.ConnectDB()
	defer database.CloseDB()

	// Create game server
	server := models.GameServer{
		Address:    "localhost",
		Port:      8081,
		MaxPlayers: 2,
		Status:    "available",
	}
	database.DB.Create(&server)

	// Create two users with similar ranks
	user1 := models.User{
		Username: "player1",
		Password: "password123",
		Stats: models.Stats{
			Rank: models.Silver,
		},
	}
	user2 := models.User{
		Username: "player2",
		Password: "password123",
		Stats: models.Stats{
			Rank: models.Silver,
		},
	}
	database.DB.Create(&user1)
	database.DB.Create(&user2)

	// Add both users to queue
	queue1 := models.QueuedPlayer{
		UserID:   user1.ID,
		JoinedAt: time.Now().Unix(),
		Status:   "queued",
	}
	queue2 := models.QueuedPlayer{
		UserID:   user2.ID,
		JoinedAt: time.Now().Unix() + 1,
		Status:   "queued",
	}
	database.DB.Create(&queue1)
	database.DB.Create(&queue2)

	// Trigger matchmaking
	controllers.matchPlayers()

	// Verify match was created
	var game models.Game
	err := database.DB.Preload("Players").Where("server_id = ?", server.ID).First(&game).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(game.Players))

	// Verify players were matched
	var matchedPlayers []models.QueuedPlayer
	database.DB.Where("status = ?", "matched").Find(&matchedPlayers)
	assert.Equal(t, 2, len(matchedPlayers))

	// Verify server status was updated
	var updatedServer models.GameServer
	database.DB.First(&updatedServer, server.ID)
	assert.Equal(t, "in_game", updatedServer.Status)
}
