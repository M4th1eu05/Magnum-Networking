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

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupStoreRouter() *gin.Engine {
	r := gin.Default()
	authorized := r.Group("/", middlewares.AuthMiddleware())
	{
		authorized.GET("/store/items", controllers.GetStoreItems)
		authorized.GET("/user/:username/items", controllers.GetUserItems)
		authorized.POST("/store/items/:id/purchase", controllers.PurchaseItem)
		authorized.POST("/store/items/:id/equip", controllers.EquipItem)
	}

	admin := r.Group("/admin", middlewares.AuthMiddleware(), middlewares.AdminMiddleware())
	{
		admin.POST("/store/items", controllers.CreateStoreItem)
		admin.PUT("/store/items/:id", controllers.UpdateStoreItem)
		admin.DELETE("/store/items/:id", controllers.DeleteStoreItem)
	}
	return r
}

func TestGetStoreItems(t *testing.T) {
	// Setup
	database.ConnectDB()
	defer database.CloseDB()

	// Create test items
	testItem := models.StoreItem{
		Name:        "Test Item",
		Description: "A test item",
		Price:       100,
		Type:        "skin",
	}
	database.DB.Create(&testItem)

	router := setupStoreRouter()

	// Create test user and generate token
	testUser := models.User{Username: "testuser", Password: "password123", IsAdmin: false}
	database.DB.Create(&testUser)
	token, _ := middlewares.GenerateToken(testUser)

	// Test GetStoreItems
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/store/items", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var response []models.StoreItem
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response), 1)
	assert.Equal(t, "Test Item", response[0].Name)
}

func TestPurchaseItem(t *testing.T) {
	// Setup
	database.ConnectDB()
	defer database.CloseDB()

	// Create test item
	testItem := models.StoreItem{
		Name:        "Purchase Test Item",
		Description: "An item to purchase",
		Price:       100,
		Type:        "skin",
	}
	database.DB.Create(&testItem)

	// Create test user
	testUser := models.User{Username: "purchaser", Password: "password123", IsAdmin: false}
	database.DB.Create(&testUser)
	token, _ := middlewares.GenerateToken(testUser)

	router := setupStoreRouter()

	// Test purchase
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/store/items/"+string(testItem.ID), nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify item is in user's inventory
	var userItems []models.UserItem
	database.DB.Where("user_id = ? AND item_id = ?", testUser.ID, testItem.ID).Find(&userItems)
	assert.Equal(t, 1, len(userItems))
}

func TestEquipItem(t *testing.T) {
	// Setup
	database.ConnectDB()
	defer database.CloseDB()

	// Create test item
	testItem := models.StoreItem{
		Name:        "Equip Test Item",
		Description: "An item to equip",
		Price:       100,
		Type:        "skin",
	}
	database.DB.Create(&testItem)

	// Create test user and give them the item
	testUser := models.User{Username: "equipper", Password: "password123", IsAdmin: false}
	database.DB.Create(&testUser)
	userItem := models.UserItem{
		UserID:    testUser.ID,
		ItemID:    testItem.ID,
		Equipped:  false,
	}
	database.DB.Create(&userItem)

	token, _ := middlewares.GenerateToken(testUser)
	router := setupStoreRouter()

	// Test equip
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/store/items/"+string(testItem.ID)+"/equip", nil)
	req.Header.Set("Authorization", token)
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify item is equipped
	var updatedUserItem models.UserItem
	database.DB.Where("user_id = ? AND item_id = ?", testUser.ID, testItem.ID).First(&updatedUserItem)
	assert.True(t, updatedUserItem.Equipped)
}

func TestCreateStoreItem(t *testing.T) {
	// Setup
	database.ConnectDB()
	defer database.CloseDB()

	// Create admin user
	adminUser := models.User{Username: "admin", Password: "admin123", IsAdmin: true}
	database.DB.Create(&adminUser)
	token, _ := middlewares.GenerateToken(adminUser)

	router := setupStoreRouter()

	// Test create item
	newItem := `{"name":"New Item","description":"A new item","price":200,"type":"effect"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/admin/store/items", strings.NewReader(newItem))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, w.Code)
	var response models.StoreItem
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "New Item", response.Name)
}
