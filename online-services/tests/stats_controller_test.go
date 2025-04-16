package tests

import (
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
)

func TestGetUserStatsSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()
	utils.InitValidators()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := models.User{Username: "testuser", Password: string(hashedPassword)}
	database.DB.Create(&testUser)

	token, _ := middlewares.GenerateToken(testUser)

	router := gin.Default()
	utils.InitValidators()
	router.GET("/user/:username/stats", middlewares.AuthMiddleware(), controllers.GetUserStats)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/testuser/stats", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")
	assert.Contains(t, w.Body.String(), "score")
}

func TestGetUserStatsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()
	utils.InitValidators()

	router := gin.Default()
	utils.InitValidators()
	router.GET("/user/:username/stats", middlewares.AuthMiddleware(), controllers.GetUserStats)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/testuser/stats", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetUserStatsUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()
	utils.InitValidators()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := models.User{Username: "testuser", Password: string(hashedPassword)}
	database.DB.Create(&testUser)

	token, _ := middlewares.GenerateToken(testUser)

	router := gin.Default()
	utils.InitValidators()
	router.GET("/user/:username/stats", middlewares.AuthMiddleware(), controllers.GetUserStats)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/unknownuser/stats", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
