package tests

import (
	"net/http"
	"net/http/httptest"
	"online-services/controllers"
	"online-services/database"
	"online-services/middlewares"
	"online-services/models"
	"online-services/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGetStatsSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()
	utils.InitValidators()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := models.User{Username: "testuser", Password: string(hashedPassword)}
	database.DB.Create(&testUser)

	token, _ := middlewares.GenerateToken(testUser)

	router := gin.Default()
	router.GET("/stats", middlewares.AuthMiddleware(), controllers.GetStats)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stats", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")
	assert.Contains(t, w.Body.String(), "stats")
}

func TestGetStatsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()
	utils.InitValidators()

	router := gin.Default()
	router.GET("/stats", middlewares.AuthMiddleware(), controllers.GetStats)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stats", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetStatsUserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()
	utils.InitValidators()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	nonExistentUser := models.User{Username: "nonexistent", Password: string(hashedPassword)}
	token, _ := middlewares.GenerateToken(nonExistentUser)

	router := gin.Default()
	router.GET("/stats", middlewares.AuthMiddleware(), controllers.GetStats)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/stats", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Utilisateur non trouvé")
}

func TestUpdateStatsAndAchievements(t *testing.T) {
	gin.SetMode(gin.TestMode)
	database.ConnectDB()
	defer database.CloseDB()

	// Create test user
	user := models.User{Username: "testuser", Password: "password123"}
	database.DB.Create(&user)

	// Create an achievement
	achievement := models.Achievement{Name: "High Score", Threshold: 50, StatsName: "score"}
	database.DB.Create(&achievement)

	// Update stats
	stats := []controllers.StatsInfo{
		{UserUUID: user.UUID, StatName: "score", Value: 60},
	}
	err := controllers.UpdateStats(user.UUID, stats)
	assert.NoError(t, err)

	// Check if the achievement was unlocked
	var userAchievements []models.Achievement
	database.DB.Model(&user).Association("Achievements").Find(&userAchievements)
	assert.Equal(t, 1, len(userAchievements))
	assert.Equal(t, "High Score", userAchievements[0].Name)
}
