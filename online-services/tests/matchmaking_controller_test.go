package tests

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"online-services/controllers"
	"online-services/database"
	"online-services/middlewares"
	"online-services/models"
	"online-services/utils"
	"testing"
	"time"
)

func TestRegisterServer(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	router := gin.Default()
	router.POST("/register-server", controllers.RegisterServer)

	serverInfo := models.GameServer{
		IP:   "127.0.0.1",
		Port: 8080,
	}
	body, _ := json.Marshal(serverInfo)

	req, _ := http.NewRequest("POST", "/register-server", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestJoinQueue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	utils.StartMatchmaking()

	router := gin.Default()
	router.POST("/join-queue", middlewares.AuthMiddleware(), controllers.JoinQueue)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := models.User{Username: "testuser", Password: string(hashedPassword)}
	database.DB.Create(&testUser)

	token, _ := middlewares.GenerateToken(testUser)

	req, _ := http.NewRequest("POST", "/join-queue", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Joined matchmaking queue")
}

func TestQueueStatusInQueue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	utils.StartMatchmaking()

	router := gin.Default()
	router.POST("/join-queue", middlewares.AuthMiddleware(), controllers.JoinQueue)
	router.GET("/queue-status", middlewares.AuthMiddleware(), controllers.QueueStatus)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := models.User{Username: "testuser", Password: string(hashedPassword)}
	database.DB.Create(&testUser)

	token, _ := middlewares.GenerateToken(testUser)

	joinReq, _ := http.NewRequest("POST", "/join-queue", nil)
	joinReq.Header.Set("Authorization", token)
	joinW := httptest.NewRecorder()
	router.ServeHTTP(joinW, joinReq)

	assert.Equal(t, http.StatusOK, joinW.Code)

	req, _ := http.NewRequest("GET", "/queue-status", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooEarly, w.Code)
	assert.Contains(t, w.Body.String(), "in queue")
}

func TestQueueStatusNotInQueue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	utils.StartMatchmaking()

	router := gin.Default()
	router.GET("/queue-status", middlewares.AuthMiddleware(), controllers.QueueStatus)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := models.User{Username: "testuser", Password: string(hashedPassword)}
	database.DB.Create(&testUser)

	token, _ := middlewares.GenerateToken(testUser)

	req, _ := http.NewRequest("GET", "/queue-status", nil)
	req.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "not in queue")
}

func TestLeaveQueue(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	utils.StartMatchmaking()

	router := gin.Default()
	router.POST("/join-queue", middlewares.AuthMiddleware(), controllers.JoinQueue)
	router.POST("/leave-queue", middlewares.AuthMiddleware(), controllers.LeaveQueue)
	router.GET("/queue-status", middlewares.AuthMiddleware(), controllers.QueueStatus)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := models.User{Username: "testuser", Password: string(hashedPassword)}
	database.DB.Create(&testUser)
	token, _ := middlewares.GenerateToken(testUser)

	joinReq, _ := http.NewRequest("POST", "/join-queue", nil)
	joinReq.Header.Set("Authorization", token)
	joinW := httptest.NewRecorder()
	router.ServeHTTP(joinW, joinReq)

	assert.Equal(t, http.StatusOK, joinW.Code)

	leaveReq, _ := http.NewRequest("POST", "/leave-queue", nil)
	leaveReq.Header.Set("Authorization", token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, leaveReq)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Left matchmaking queue")

	statusReq, _ := http.NewRequest("GET", "/queue-status", nil)
	statusReq.Header.Set("Authorization", token)
	statusW := httptest.NewRecorder()
	router.ServeHTTP(statusW, statusReq)

	assert.Equal(t, http.StatusBadRequest, statusW.Code)
	assert.Contains(t, statusW.Body.String(), "not in queue")
}

func TestStartGame(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	router := gin.Default()
	router.POST("/start-game", controllers.StartGame)

	server := models.GameServer{
		IP:   "127.0.0.1",
		Port: 8080,
	}
	database.DB.Create(&server)

	body, _ := json.Marshal(server)
	req, _ := http.NewRequest("POST", "/start-game", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Match started")
}

func TestEndGame(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	router := gin.Default()
	router.POST("/end-game", controllers.EndGame)

	// Create test data
	game := models.Game{}
	database.DB.Create(&game)

	player := models.User{GameID: &game.ID}
	database.DB.Create(&player)

	stats := []controllers.StatsInfo{
		{UserUUID: player.UUID, StatName: "kills", Value: 10},
	}
	gameInfo := controllers.GameInfo{
		GameID: game.ID,
		Stats:  stats,
	}
	body, _ := json.Marshal(gameInfo)

	req, _ := http.NewRequest("POST", "/end-game", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Match ended")
}

func TestMatchPlayers(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	utils.StartMatchmaking()

	router := gin.Default()
	router.POST("/join-queue", middlewares.AuthMiddleware(), controllers.JoinQueue)
	router.POST("/queue-status", middlewares.AuthMiddleware(), controllers.QueueStatus)

	users := []models.User{
		{Username: "user1", Password: "password1"},
		{Username: "user2", Password: "password2"},
		{Username: "user3", Password: "password3"},
		{Username: "user4", Password: "password4"},
	}
	server := models.GameServer{IP: "127.0.0.1", Port: 8080}

	tokens := make([]string, len(users))
	database.DB.Create(&users)
	database.DB.Create(&server)

	for i, user := range users {
		token, _ := middlewares.GenerateToken(user)
		tokens[i] = token

		req := httptest.NewRequest("POST", "/join-queue", nil)
		req.Header.Set("Authorization", token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Joined matchmaking queue")
	}

	time.Sleep(2 * time.Second)

	for _, token := range tokens {
		req := httptest.NewRequest("POST", "/queue-status", nil)
		req.Header.Set("Authorization", token)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "server")
		assert.Contains(t, w.Body.String(), server.IP)
	}
}
