//package controllers
//
//import (
//	"net/http"
//	"online-services/database"
//	"online-services/models"
//
//	"github.com/gin-gonic/gin"
//	"github.com/golang-jwt/jwt"
//)
//
//// GetStoreItems returns all available store items
//func GetStoreItems(c *gin.Context) {
//	var items []models.StoreItem
//	if err := database.DB.Find(&items).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching store items"})
//		return
//	}
//	c.JSON(http.StatusOK, items)
//}
//
//// GetUserItems returns all items owned by a user
//func GetUserItems(c *gin.Context) {
//	username := c.Param("username")
//	var user models.User
//	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//		return
//	}
//
//	var userItems []models.UserItem
//	if err := database.DB.Preload("StoreItem").Where("user_id = ?", user.ID).Find(&userItems).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user items"})
//		return
//	}
//
//	c.JSON(http.StatusOK, userItems)
//}
//
//// PurchaseItem handles the purchase of a store item
//func PurchaseItem(c *gin.Context) {
//	itemID := c.Param("id")
//
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
//	// Check if user already owns the item
//	var existingItem models.UserItem
//	result := database.DB.Where("user_id = ? AND item_id = ?", user.ID, itemID).First(&existingItem)
//	if result.RowsAffected > 0 {
//		c.JSON(http.StatusConflict, gin.H{"error": "Item already owned"})
//		return
//	}
//
//	// Create new user item
//	userItem := models.UserItem{
//		UserID: user.ID,
//		ItemID: uint(itemID),
//	}
//	if err := database.DB.Create(&userItem).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to purchase item"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Item purchased successfully"})
//}
//
//// EquipItem handles equipping/unequipping an item
//func EquipItem(c *gin.Context) {
//	itemID := c.Param("id")
//
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
//	// Find the user's item
//	var userItem models.UserItem
//	if err := database.DB.Where("user_id = ? AND item_id = ?", user.ID, itemID).First(&userItem).Error; err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found or not owned"})
//		return
//	}
//
//	// Toggle equipped status
//	userItem.Equipped = !userItem.Equipped
//	if err := database.DB.Save(&userItem).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item status"})
//		return
//	}
//
//	c.JSON(http.StatusOK, userItem)
//}
//
//// Admin routes for store management
//
//// CreateStoreItem creates a new store item (admin only)
//func CreateStoreItem(c *gin.Context) {
//	var item models.StoreItem
//	if err := c.ShouldBindJSON(&item); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	if err := database.DB.Create(&item).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create store item"})
//		return
//	}
//
//	c.JSON(http.StatusCreated, item)
//}
//
//// UpdateStoreItem updates an existing store item (admin only)
//func UpdateStoreItem(c *gin.Context) {
//	id := c.Param("id")
//	var item models.StoreItem
//
//	if err := database.DB.First(&item, id).Error; err != nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "Store item not found"})
//		return
//	}
//
//	if err := c.ShouldBindJSON(&item); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	if err := database.DB.Save(&item).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update store item"})
//		return
//	}
//
//	c.JSON(http.StatusOK, item)
//}
//
//// DeleteStoreItem deletes a store item (admin only)
//func DeleteStoreItem(c *gin.Context) {
//	id := c.Param("id")
//	if err := database.DB.Delete(&models.StoreItem{}, id).Error; err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete store item"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Store item deleted successfully"})
//}
