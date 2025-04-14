package controllers

import (
	"net/http"
	"online-services/database"
	"online-services/models"
	"time"

	"github.com/gin-gonic/gin"
)

// GetAchievements returns all available achievements
func GetAchievements(c *gin.Context) {
	var achievements []models.Achievement
	if err := database.DB.Find(&achievements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching achievements"})
		return
	}
	c.JSON(http.StatusOK, achievements)
}

// GetUserAchievements returns achievements for a specific user
func GetUserAchievements(c *gin.Context) {
	username := c.Param("username")
	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var achievements []models.UserAchievement
	if err := database.DB.Preload("Achievement").Where("user_id = ?", user.ID).Find(&achievements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user achievements"})
		return
	}

	c.JSON(http.StatusOK, achievements)
}

// CheckAchievements checks and awards achievements based on user stats
func CheckAchievements(stats models.Stats) error {
	var achievements []models.Achievement
	if err := database.DB.Find(&achievements).Error; err != nil {
		return err
	}

	for _, achievement := range achievements {
		// Check if achievement conditions are met
		achieved := false
		switch achievement.Type {
		case "games_won":
			achieved = stats.NbrGamesWon >= achievement.Threshold
		case "games_played":
			achieved = stats.NbrGamesPlayed >= achievement.Threshold
		case "cubes_spawned":
			achieved = stats.NbrCubesSpawned >= achievement.Threshold
		case "spheres_spawned":
			achieved = stats.NbrSpheresSpawned >= achievement.Threshold
		}

		if achieved {
			// Check if user already has this achievement
			var userAchievement models.UserAchievement
			result := database.DB.Where("user_id = ? AND achievement_id = ?", stats.UserID, achievement.ID).First(&userAchievement)
			if result.RowsAffected == 0 {
				// Award new achievement
				newUserAchievement := models.UserAchievement{
					UserID:        uint(stats.UserID),
					AchievementID: achievement.ID,
					UnlockedAt:    time.Now().Unix(),
				}
				if err := database.DB.Create(&newUserAchievement).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Admin routes

// CreateAchievement creates a new achievement (admin only)
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

// UpdateAchievement updates an existing achievement (admin only)
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

// DeleteAchievement deletes an achievement (admin only)
func DeleteAchievement(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Achievement{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete achievement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Achievement deleted successfully"})
}
