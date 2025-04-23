package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online-services/database"
	"online-services/models"
)

func GetAchievements(c *gin.Context) {
	UUID, exists := c.Get("UUID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID not found"})
		return
	}

	var user models.User
	if err := database.DB.Where("uuid = ?", UUID).First(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	var achievements []models.Achievement
	if err := database.DB.Where("user_id = ?", user.ID).Find(&achievements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user achievements"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}


func CheckAchievements(stats models.Stats) error {
	var achievements []models.Achievement
	if err := database.DB.Where("StatsName = ?", stats.Name).Find(&achievements).Error; err != nil {
		return err
	}


	for _, achievement := range achievements {

		if stats.Value >= achievement.Threshold {
			var user models.User
			if err := database.DB.Where("id = ?", stats.ID).First(&user).Error; err != nil {
				return err
			}
			if err := database.DB.Model(&user).Association("Achievements").Append(&achievement); err != nil {
				return err
			}
		}
	}
	return nil
}


// Admin routes
func CreateAchievement(c *gin.Context) {
	var achievement models.Achievement
	if err := c.ShouldBindJSON(&achievement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&achievement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create achievement"})
		return
	}

	c.JSON(http.StatusCreated, achievement)
}

func UpdateAchievement(c *gin.Context) {
	id := c.Param("id")
	var achievement models.Achievement

	if err := database.DB.First(&achievement, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
		return
	}

	if err := c.ShouldBindJSON(&achievement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&achievement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update achievement"})
		return
	}

	c.JSON(http.StatusOK, achievement)
}

func DeleteAchievement(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Achievement{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete achievement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Achievement deleted successfully"})
}
