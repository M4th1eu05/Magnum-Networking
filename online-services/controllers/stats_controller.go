package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"online-services/database"
	"online-services/models"
)

func GetStats(c *gin.Context) {

	UUID, exists := c.Get("UUID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UUID not found"})
		return
	}

	var user models.User
	if err := database.DB.Where("uuid = ?", UUID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
		return
	}

	var stats []models.Stat
	if err := database.DB.Where("user_id = ?", user.ID).Find(&stats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des statistiques"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"stats":    stats,
	})
}

func UpdateStats(uuid uuid.UUID, statsInfos []StatsInfo) error {
	var user models.User
	if err := database.DB.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		return err
	}
	for _, statsInfo := range statsInfos {
		var stat models.Stat
		if err := database.DB.Where("userID = ? AND name = ?", user.ID, statsInfo.StatName).First(&stat).Error; err != nil {
			stat = models.Stat{UserID: &user.ID, Name: statsInfo.StatName, Value: statsInfo.Value}
			if err := database.DB.Create(&stat).Error; err != nil {
				return err
			}
		} else {
			stat.Value += statsInfo.Value
			if err := database.DB.Save(&stat).Error; err != nil {
				return err
			}
		}
		if err := CheckAchievements(stat); err != nil {
			return err
		}
	}

	return nil
}
