package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"online-services/controllers"
	"online-services/database"
	"online-services/middlewares"
	"online-services/models"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupAchievementRouter() *gin.Engine {
	r := gin.Default()
	authorized := r.Group("/", middlewares.AuthMiddleware())
	{
		authorized.GET("/achievements", controllers.GetAchievements)
		authorized.GET("/user/:username/achievements", controllers.GetUserAchievements)
	}

	admin := r.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		admin.POST("/achievements", controllers.CreateAchievement)
		admin.PUT("/achievements/:id", controllers.UpdateAchievement)
		admin.DELETE("/achievements/:id", controllers.DeleteAchievement)
	}
	return r
}

func TestGetAchievements(t *testing.T) {
	// Setup
	database.ConnectDB()
	defer database.CloseDB()

	// Create test achievement
	testAchievement := models.Achievement{
		Name:        "Test Achievement",
		Description: "A test achievement",
		Condition:   "games_won",
		Threshold:   10,
		Type:        "games_won",
	}
	database.DB.Create(&testAchievement)

	// Create test user
	testUser := models.User{Username: "achiever", Password: "password123"}
	database.DB.Create(&testUser)
	token, _ := middlewares.GenerateToken(testUser)

	router := setupAchievementRouter()

	// Test GetAchievements
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/achievements", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var response []models.Achievement
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response), 1)
	assert.Equal(t, "Test Achievement", response[0].Name)
}

func TestGetUserAchievements(t *testing.T) {
	// Setup
	database.ConnectDB()
	defer database.CloseDB()

	// Create test achievement
	testAchievement := models.Achievement{
		Name:        "User Achievement",
		Description: "A user's achievement",
		Condition:   "games_won",
		Threshold:   5,
		Type:        "games_won",
	}
	database.DB.Create(&testAchievement)

	// Create test user with achievement
	testUser := models.User{Username: "achievement_owner", Password: "password123"}
	database.DB.Create(&testUser)
	
	userAchievement := models.UserAchievement{
		UserID:        testUser.ID,
		AchievementID: testAchievement.ID,
		UnlockedAt:    time.Now().Unix(),
	}
	database.DB.Create(&userAchievement)

	token, _ := middlewares.GenerateToken(testUser)
	router := setupAchievementRouter()

	// Test GetUserAchievements
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/user/achievement_owner/achievements", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var response []models.UserAchievement
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response))
}

func TestCreateAchievement(t *testing.T) {
	// Setup
	database.ConnectDB()
	defer database.CloseDB()

	// Create admin user
	adminUser := models.User{Username: "admin", Password: "admin123", IsAdmin: true}
	database.DB.Create(&adminUser)
	token, _ := middlewares.GenerateToken(adminUser)

	router := setupAchievementRouter()

	// Test create achievement
	newAchievement := `{
		"name": "New Achievement",
		"description": "A new achievement",
		"condition": "games_won",
		"threshold": 20,
		"type": "games_won"
	}`
	
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/admin/achievements", strings.NewReader(newAchievement))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	var response models.Achievement
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "New Achievement", response.Name)
}

func TestCheckAchievements(t *testing.T) {
	// Setup
	database.ConnectDB()
	defer database.CloseDB()

	// Create test user
	testUser := models.User{Username: "achievement_checker", Password: "password123"}
	database.DB.Create(&testUser)

	// Create test achievement
	testAchievement := models.Achievement{
		Name:        "Games Won Achievement",
		Description: "Win 5 games",
		Condition:   "games_won",
		Threshold:   5,
		Type:        "games_won",
	}
	database.DB.Create(&testAchievement)

	// Create user stats that should trigger the achievement
	stats := models.Stats{
		UserID:      uint64(testUser.ID),
		NbrGamesWon: 5,
	}

	// Check achievements
	err := controllers.CheckAchievements(stats)
	assert.NoError(t, err)

	// Verify achievement was awarded
	var userAchievement models.UserAchievement
	result := database.DB.Where("user_id = ? AND achievement_id = ?", testUser.ID, testAchievement.ID).First(&userAchievement)
	assert.Equal(t, int64(1), result.RowsAffected)
}
