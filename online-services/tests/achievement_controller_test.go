package tests

import (
	"net/http"
	"net/http/httptest"
	"online-services/controllers"
	"online-services/database"
	"online-services/middlewares"
	"online-services/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetAchievements(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	router := gin.Default()
	router.GET("/achievements", middlewares.AuthMiddleware(), controllers.GetAchievements)

	// Create test user and achievements
	user := models.User{
		Username: "testuser",
		Password: "password123",
		Stats: []models.Stat{
			{Name: "Wins", Value: 1},
		},
		Achievements: []models.Achievement{
			{Name: "First Win", Threshold: 1},
		},
	}
	database.DB.Create(&user)

	token, _ := middlewares.GenerateToken(user)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/achievements", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "First Win")
}
